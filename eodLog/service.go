package eodLog

import (
	"bufio"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/kyokomi/emoji"
	"github.com/matryer/try"
	"github.com/pkg/errors"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/callout"
	"github.com/weAutomateEverything/gppMonitor/sftp"
	"github.com/weAutomateEverything/gppMonitor/transactionCountLog"
	"golang.org/x/net/context"
	"log"
	"os"
	"strings"
	"time"
	"github.com/weAutomateEverything/gppMonitor/hal"
)

type Service interface {
	/*
		RetrieveEDOLog copies contents of EDO.log to a local file of the same name which is then analysed for the success/failure of Edo Posing request file send
	*/
	RetrieveEDOLog()
	response(message, filename, dateStamp, timeStamp string) string
}

type service struct {
	sftpService        sftp.Service
	alertService       alert.Service
	calloutService     callout.Service
	transactionService transactionCountLog.Service
}

func NewService(callout callout.Service, sftpService sftp.Service, alertService alert.Service) Service {
	s := &service{calloutService: callout, sftpService: sftpService, alertService: alertService}
	go func() {
		s.schedule()
	}()
	return s
}

func (s *service) schedule() {
	retreiveEDOLog := gocron.NewScheduler()

	go func() {
		retreiveEDOLog.Every(1).Day().At("01:10").Do(s.RetrieveEDOLog)
		<-retreiveEDOLog.Start()
	}()
}

func (s *service) RetrieveEDOLogMethod() (r error) {

	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(error); ok {
				r = errors.New(e.Error())
			}
			r = errors.New("EDO file confirmation failed")
		}
	}()

	s.sftpService.RetrieveFile(getEDOLogLocation(), "EDO.log")

	dateLine, lastLine := lastLines()
	if len(dateLine) < 30 {
		panic(errors.Errorf(emoji.Sprintf(":rotating_light: Timestamp line within EDO.log is smaller than 30 characters which will result in and index out of bounds error.")))
	}
	year := dateLine[18:22]
	month := dateLine[23:25]
	day := dateLine[26:28]

	time := dateLine[29:]

	date := day + month + year
	dateStamp, timeStamp := dateTimeConvert(date, time)

	if dateStamp == "01/01/0001" {
		s.alertService.SendAlert(context.TODO(), hal.Chatid(),emoji.Sprintf(":rotating_light: EDO.log timestamp format has changed. Unable to parse date/time."))
		log.Println("EDO.log timestamp format has changed. Unable to parse date/time.")
	}
	fileName := fileNameExtract(lastLine)

	s.alertService.SendAlert(context.TODO(),hal.Chatid(), s.response(lastLine, fileName, dateStamp, timeStamp))

	os.Remove("/tmp/EDO.log")

	return nil
}

func lastLines() (string, string) {

	f := openFile("/tmp/EDO.log")

	buf := make([]string, 32*1024)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		buf = append(buf, line)
	}

	if buf[len(buf)-1] == "" {
		date := buf[len(buf)-3]
		result := buf[len(buf)-2]
		return date, result
	}
	date := buf[len(buf)-2]
	result := buf[len(buf)-1]

	return date, result
}

func openFile(targetFile string) *os.File {
	f, err := os.Open(targetFile)
	if err != nil {
		panic(err)
	}
	return f
}

func dateTimeConvert(date, logTime string) (string, string) {

	d, e := time.Parse("02012006", date)
	if e != nil {
		panic(fmt.Errorf("EDO.log timestamp(date portion specifically) failed to parse with the following error: %v", e))
	}
	t, e := time.Parse("15:04:05", logTime)
	if e != nil {
		panic(fmt.Errorf("EDO.log timestamp(time portion specifically) failed  to parse with the following error: %v", e))
	}
	dFormat := d.Format("02/01/2006")
	tFormat := t.Format("3:04PM")
	return dFormat, tFormat
}

func fileNameExtract(logEntry string) string {
	fileName := logEntry[13:50]
	fileName = strings.Replace(fileName, "_", "-", -1)
	return fileName
}

func (s *service) response(message, filename, dateStamp, timeStamp string) string {

	currentDate := time.Now()
	cd := currentDate.Format("02/01/2006")

	if strings.Contains(message, "successful") && cd == dateStamp {
		transactionsToBeProcessed := s.transactionService.RetrieveNightFileTransactions(filename)
		return emoji.Sprintf(":white_check_mark: EDO Posting request file '%s' successfully sent on the: %s at %s\n----------\nTransactions: ", filename, dateStamp, timeStamp, transactionsToBeProcessed)
	} else if strings.Contains(message, "failed") && cd == dateStamp {
		s.calloutService.InvokeCallout(context.TODO(), hal.Chatid(),"EDO Posting request file send failed", fmt.Sprintf("EDO Posting request file '%s' send failed on the: %s at %s", filename, dateStamp, timeStamp),hal.AlexaVars())
		return emoji.Sprintf(":rotating_light: EDO Posting request file '%s' send failed on the: %s at %s", filename, dateStamp, timeStamp)
	} else if cd != dateStamp {
		panic("Retrying EDO.log retrieval")
	}
	return emoji.Sprintf(":red_circle: Error extracting log timestamp or success/failure result on the: %s at %s. Please consult EDO log file directly", dateStamp, timeStamp)
}

func (s *service) RetrieveEDOLog() {
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		err = s.RetrieveEDOLogMethod()
		if err != nil {
			log.Println("next attempt in 2 minutes")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return attempt < 5, err //5 attempts
	})
	if err != nil {
		log.Println(err)
	}
}

func getEDOLogLocation() string {
	return os.Getenv("EDO_LOCATION")
}
