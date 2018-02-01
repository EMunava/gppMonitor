package eodLog

import (
	"os"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/sftp"
	"bufio"
	"time"
	"strings"
	"github.com/zamedic/go2hal/alert"
)

type Service interface {
	/*
	RetrieveEDOLog copies contents of EDO.log to a local file of the same name which is then analysed for the success/failure of Edo Posing request file send
	*/
	RetrieveEDOLog()
}

type service struct {
	sftpService sftp.Service
	alertService alert.Service
}

func NewService(sftpService sftp.Service,alertService alert.Service) Service {
	return &service{sftpService:sftpService,alertService:alertService}
}

func (s *service) RetrieveEDOLog() {

	s.sftpService.RetrieveFile("/cdwasha/connectdirect/outgoing/EDO_DirectDebitRequest/", "EDO.log")

	dateLine, lastLine := lastLines()

	dateStamp := dateConvert(dateLine)

	s.alertService.SendAlert(response(lastLine, dateStamp))

	os.Remove("/tmp/EDO.log")
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
	dt, _ := time.Parse("Mon Jan _2 15:04:05 MST 2006", dtstr1)
	dtstr2 := dt.Format("02/01/2006")
	return dtstr2
}

func response(message, dateStamp string) string {

	currentDate := time.Now()
	cd := currentDate.Format("02/01/2006")

	if strings.Contains(message, "successful") && cd == dateStamp {
		return "EDO Posting request file successfully sent"
	} else if strings.Contains(message, "failed") && cd == dateStamp {
		return "EDO Posting request file send failed"
	} else if cd != dateStamp {
		return "Last log entry and current date do not correlate"
	}
	return "Error extracting log timestamp or success/failure result. Please consult log EDO file directly"
}
