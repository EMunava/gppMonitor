package sftp

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

/*
RetrieveEDOLog copies contents of EDO.log to a local file of the same name which is then analysed for the success/failure of Edo Posing request file send
*/
func RetrieveEDOLog() {

	retrieveFile("/cdwasha/connectdirect/outgoing/EDO_DirectDebitRequest/", "EDO.log")

	dateLine, lastLine := lastLines()

	dateStamp := dateConvert(dateLine)

	sendAlert(response(lastLine, dateStamp))

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

func dateConvert(date string) string {
	dtstr1 := date
	dt, _ := time.Parse("Mon Jan _2 15:04:05 MST 2006", dtstr1)
	dtstr2 := dt.Format("02/01/2006")
	return dtstr2
}

func openFile(targetFile string) *os.File {
	f, err := os.Open("/tmp/EDO.log")
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func response(message, dateStamp string) string {

	currentDate := time.Now()
	cd := currentDate.Format("02/01/2006")

	if strings.Contains(message, "successful") && cd == dateStamp {
		log.Println("EDO Posting request file successfully sent")
		return "EDO Posting request file successfully sent"
	} else if strings.Contains(message, "failed") && cd == dateStamp {
		log.Println("EDO Posting request file send failed!!")
		return "EDO Posting request file send failed"
	} else if cd != dateStamp {
		log.Println("Last log entry timestamp and current date do not correlate")
		return "Last log entry and current date do not correlate"
	}
	return "Error extracting log timestamp or success/failure result. Please consult log EDO file directly"
}
