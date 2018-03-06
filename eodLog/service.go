package eodLog

import (
	"bufio"
	"github.com/jasonlvhit/gocron"
	"github.com/kyokomi/emoji"
	"github.com/matryer/try"
	"github.com/pkg/errors"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/gppMonitor/sftp"
	"golang.org/x/net/context"
	"log"
	"os"
	"strings"
	"time"
)

type Service interface {
	/*
		RetrieveEDOLog copies contents of EDO.log to a local file of the same name which is then analysed for the success/failure of Edo Posing request file send
	*/
	RetrieveEDOLog()
}

type service struct {
	sftpService  sftp.Service
	alertService alert.Service
}

func NewService(sftpService sftp.Service, alertService alert.Service) Service {
	s := &service{sftpService: sftpService, alertService: alertService}
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
	dateLine = dateLine[18:]
	dateStamp := dateConvert(dateLine)
	if dateStamp == "01/01/0001" {
		s.alertService.SendHeartbeatGroupAlert(context.TODO(), "EDO.log timestamp format has changed. Unable to parse date/time.")
		log.Println("EDO.log timestamp format has changed. Unable to parse date/time.")
	}
	fileName := fileNameExtract(lastLine)

	s.alertService.SendAlert(context.TODO(), response(lastLine, fileName, dateStamp))

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

func dateConvert(date string) string {
	dtstr1 := date
	dt, _ := time.Parse("060102 030405", dtstr1)
	dtstr2 := dt.Format("02/01/2006")
	return dtstr2
}

func fileNameExtract(logEntry string) string {
	fileName := logEntry[13:50]
	fileName = strings.Replace(fileName,"_","-", -1)
	return fileName
}

func response(message, filename, dateStamp string) string {

	currentDate := time.Now()
	cd := currentDate.Format("02/01/2006")

	if strings.Contains(message, "successful") && cd == dateStamp {
		return emoji.Sprintf(":white_check_mark: EDO Posting request file '%s' successfully sent at: %s", filename, dateStamp)
	} else if strings.Contains(message, "failed") && cd == dateStamp {
		return emoji.Sprintf(":rotating_light: EDO Posting request file '%s' send failed at: %s", filename, dateStamp)
	} else if cd != dateStamp {
		panic("Retrying EDO.log retrieval")
	}
	return emoji.Sprintf(":red_circle: Error extracting log timestamp or success/failure result. Please consult EDO log file directly")
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
