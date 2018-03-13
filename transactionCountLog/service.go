package transactionCountLog

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/kyokomi/emoji"
	"github.com/matryer/try"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/gppMonitor/sftp"
	"golang.org/x/net/context"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//Service interface exposes the transaction retrieval methods
type Service interface {
	RetrieveSAPTransactions()
	RetrieveLEGTransactions()
	RetrieveLEGSAPTransactions()
	retreiveTransactions(contains string, exclude ...string) error
}

type service struct {
	sftpService  sftp.Service
	alertService alert.Service
}

type transactionStatus struct {
	rejected  int
	tracking  int
	processed int
}

//NewService function creates instances of required external service structs for local use
func NewService(sftpService sftp.Service, alertService alert.Service) Service {
	s := &service{sftpService: sftpService, alertService: alertService}
	go func() {
		s.schedule()
	}()
	return s
}

func (s *service) schedule() {
	retreiveSAP := gocron.NewScheduler()
	retreiveLEG := gocron.NewScheduler()
	retreiveLEGSAP := gocron.NewScheduler()

	go func() {
		retreiveSAP.Every(1).Day().At("00:05").Do(s.RetrieveSAPTransactions)
		<-retreiveSAP.Start()
	}()
	go func() {
		retreiveLEG.Every(1).Day().At("01:32").Do(s.RetrieveLEGTransactions)
		<-retreiveLEG.Start()
	}()
	go func() {
		retreiveLEGSAP.Every(1).Day().At("00:20").Do(s.RetrieveLEGSAPTransactions)
		<-retreiveLEGSAP.Start()
	}()
}

func (s *service) retreiveTransactions(contains string, exclude ...string) (r error) {

	defer func() {
		if err := recover(); err != nil {

			if e, ok := err.(error); ok {
				r = errors.New(e.Error())
			}
			r = fmt.Errorf("%v file has not yet arrived from EDO at: %v", contains, time.Now().Format("3:04PM"))
		}
	}()

	var (
		fPath, fName string
		e            error
	)

	if len(exclude) != 0 {
		fPath, fName, e = s.pathToMostRecentFile(transactionFileLocation(), contains, exclude[0])

	} else {
		fPath, fName, e = s.pathToMostRecentFile(transactionFileLocation(), contains)
	}

	if e != nil {
		panic(e)
	}

	s.sftpService.RetrieveFile(fPath, fName)

	lineCount, transactionBreakDown := lastLines("/tmp/", fName)

	SAPTransAmount := extractTransactionAmount(lineCount)

	s.alertService.SendAlert(context.TODO(), emoji.Sprintf(":white_check_mark: %v transaction count for %v: %v\nSuccessful: %v\nRejected: %v\nTracking: %v", contains, time.Now().Format("02/01/2006"), strconv.Itoa(SAPTransAmount), transactionBreakDown.processed, transactionBreakDown.rejected, transactionBreakDown.tracking))

	log.Printf("%v transaction count for %v: %v\nSuccessful: %v\nRejected: %v\nTracking: %v", contains, time.Now().Format("02/01/2006"), strconv.Itoa(SAPTransAmount), transactionBreakDown.processed, transactionBreakDown.rejected, transactionBreakDown.tracking)

	os.Remove("/tmp/" + fName)

	return nil
}

func (s *service) RetrieveSAPTransactionsMethod() error {
	err := s.retreiveTransactions("RESPONSE.SAP")
	return err
}

func (s *service) RetrieveLEGTransactionsMethod() error {
	err := s.retreiveTransactions("RESPONSE.LEG", "SAP")
	return err

}

func (s *service) RetrieveLEGSAPTransactionsMethod() error {
	err := s.retreiveTransactions("RESPONSE.LEG.SAP")
	return err
}

func openFile(targetFile string) *os.File {
	f, err := os.Open(targetFile)
	if err != nil {
		panic(err)
	}
	return f
}

func lastLines(logLocalLocation, logFile string) (string, transactionStatus) {

	f := openFile(logLocalLocation + logFile)

	buf := make([]string, 32*1024)
	scanner := bufio.NewScanner(f)
	ts := transactionStatus{}
	for scanner.Scan() {
		line := scanner.Text()
		responseCode := line[67:69]

		if responseCode == "00" {
			ts.processed += 1
		}

		if responseCode == "02" {
			ts.tracking += 1
		}
		if responseCode == "99" {
			ts.rejected += 1
		}
		buf = append(buf, line)
	}

	if buf[len(buf)-1] == "" {
		result := buf[len(buf)-2]
		return result, ts
	}
	result := buf[len(buf)-1]

	return result, ts
}

func extractTransactionAmount(trans string) int {
	trans = trans[32:]
	re := regexp.MustCompile("[0-9]+")
	ar := re.FindAllString(trans, -1)
	transInt, err := strconv.Atoi(ar[0])
	if err != nil {
		panic("String conversion failed")
	}
	return transInt
}

func (s *service) pathToMostRecentFile(dirPath, fileContains string, exclude ...string) (string, string, error) {

	fileList, err := s.sftpService.GetFilesInPath(dirPath)
	if err != nil {

	}

	currentDate := time.Now().Format("02/01/2006")

	for _, file := range fileList {
		cont := strings.Contains(file.Name, fileContains)
		ex := false
		if len(exclude) != 0 {
			ex = strings.Contains(file.Name, exclude[0])
		}
		daDate := file.LastModified.Format("02/01/2006")
		if daDate == currentDate && cont == true && ex == false {
			return dirPath, file.Name, nil
		}
	}
	return "", "", fmt.Errorf("%v file has not arrived yet", fileContains)
}

func (s *service) RetrieveSAPTransactions() {
	err := try.Do(func(attempt int) (bool, error) {
		try.MaxRetries = 120
		var err error
		err = s.RetrieveSAPTransactionsMethod()
		if err != nil {
			log.Println("SAP file not detected. Next attempt in 2 minutes")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}

		return true, err
	})
	if err != nil {
		log.Println(err)
		s.alertService.SendAlert(context.TODO(), err.Error())
	}
}
func (s *service) RetrieveLEGTransactions() {
	err := try.Do(func(attempt int) (bool, error) {
		try.MaxRetries = 120
		var err error
		err = s.RetrieveLEGTransactionsMethod()
		if err != nil {
			log.Println("LEG file not detected. Next attempt in 2 minutes")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return true, err
	})
	if err != nil {
		log.Println(err)
		s.alertService.SendAlert(context.TODO(), err.Error())
	}
}
func (s *service) RetrieveLEGSAPTransactions() {
	err := try.Do(func(attempt int) (bool, error) {
		try.MaxRetries = 120
		var err error
		err = s.RetrieveLEGSAPTransactionsMethod()
		if err != nil {
			log.Println("LEG.SAP file not detected. Next attempt in 2 minutes")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return true, err
	})
	if err != nil {
		log.Println(err)
		s.alertService.SendAlert(context.TODO(), err.Error())
	}
}

func transactionFileLocation() string {
	return os.Getenv("TRANSACTION_LOCATION")
}