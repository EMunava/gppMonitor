package extractFooterTransactions

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/sftp"
	"github.com/jasonlvhit/gocron"
	"github.com/matryer/try"
	"github.com/zamedic/go2hal/alert"
	"log"
	"os"
	"path/filepath"
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
}

type service struct {
	sftpService  sftp.Service
	alertService alert.Service
}

type fileInfo struct {
	Name    string
	ModTime string
	Size    int64
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
			r = fmt.Errorf("%v file has not arrived", contains)
		}
	}()

	var (
		fPath, fName string
		e            error
	)

	if len(exclude) != 0 {
		fPath, fName, e = pathToMostRecentFile("/cdwasha/connectdirect/incoming/EDO_DirectDebitRequest/", contains, exclude[0])

	} else {
		fPath, fName, e = pathToMostRecentFile("/cdwasha/connectdirect/incoming/EDO_DirectDebitRequest/", contains)
	}

	if e != nil {
		panic(e)
	}

	s.sftpService.RetrieveFile(fPath, fName)

	SAPTransAmount := extractTransactionAmount(lastLines(fPath + fName))

	s.alertService.SendHeartbeatGroupAlert(string(SAPTransAmount))

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

func lastLines(logFile string) string {

	f := openFile(logFile)

	buf := make([]string, 32*1024)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		buf = append(buf, line)
	}

	if buf[len(buf)-1] == "" {
		result := buf[len(buf)-2]
		return result
	}
	result := buf[len(buf)-1]

	return result
}

func extractTransactionAmount(trans string) int {
	trans = trans[32 : len(trans)-1]
	re := regexp.MustCompile("[0-9]+")
	ar := re.FindAllString(trans, -1)
	transInt, err := strconv.Atoi(ar[0])
	if err != nil {
		panic("String conversion failed")
	}
	return transInt
}

func fileStat(logFile string) fileInfo {
	info, err := os.Stat(logFile)
	if err != nil {
		panic("File not found")
	}
	fileI := fileInfo{
		Name:    info.Name(),
		ModTime: info.ModTime().Format("02/01/2006"),
		Size:    info.Size(),
	}
	return fileI
}

func pathToMostRecentFile(dirPath, fileContains string, exclude ...string) (string, string, error) {

	fileList := []string{}
	currentDate := time.Now().Format("02/01/2006")
	filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	for _, file := range fileList {
		cont := strings.Contains(file, fileContains)
		ex := false
		fileI := fileStat(file)
		if len(exclude) != 0 {
			ex = strings.Contains(file, exclude[0])
		}
		if fileI.ModTime == currentDate && cont == true && ex == false {
			return dirPath, fileI.Name, nil
		}
	}
	return "", "", fmt.Errorf("%v file has not arrived yet", fileContains)
}

func (s *service) RetrieveSAPTransactions() {
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		err = s.RetrieveSAPTransactionsMethod()
		if err != nil {
			log.Println("SAP file not detected. Next attempt in 2 minutes")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return attempt < 120, err //120 attempts. next 4 hours
	})
	if err != nil {
		log.Println(err)
	}
}
func (s *service) RetrieveLEGTransactions() {
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		err = s.RetrieveLEGTransactionsMethod()
		if err != nil {
			log.Println("LEG file not detected. Next attempt in 2 minutes")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return attempt < 120, err //120 attempts. next 4 hours
	})
	if err != nil {
		log.Println(err)
	}
}
func (s *service) RetrieveLEGSAPTransactions() {
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		err = s.RetrieveLEGSAPTransactionsMethod()
		if err != nil {
			log.Println("LEG.SAP file not detected. Next attempt in 2 minutes")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return attempt < 120, err //120 attempts. next 4 hours
	})
	if err != nil {
		log.Println(err)
	}
}
