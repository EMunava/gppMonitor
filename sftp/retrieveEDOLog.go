package sftp

import (
	"bufio"
	"bytes"
	"io"
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

	numLines, err := lineCounter()
	if err != nil {
		log.Print(err)
	}
	ll, dateLine, _, err := lastLine(numLines)
	if err != nil {
		log.Print(err)
	}

	dateStamp := dateConvert(dateLine)

	currentDate := time.Now()
	cd := currentDate.Format("02/01/2006")

	if strings.Contains(ll, "successful") && cd == dateStamp {
		sendAlert("EDO Posting request file successfully sent")
		log.Println("EDO Posting request file successfully sent")
	} else if strings.Contains(ll, "failed") && cd != dateStamp {
		sendAlert("EDO Posting request file send failed")
		log.Println("EDO Posting request file send failed!!")
	} else if cd != dateStamp {
		sendAlert("Last log entry and current date do not correlate")
		log.Println("Last log entry timestamp and current date do not correlate")
	} else {
		log.Println("Error extracting log timestamp or success/failure result. Please consult log EDO file")
	}

	os.Remove("/tmp/EDO.log")
}

func lastLine(lineNum int) (line, dateLine string, lastLine int, err error) {
	r, err := os.Open("/tmp/EDO.log")
	if err != nil {
		log.Print(err)
	}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == (lineNum - 1) {
			dateLine = sc.Text()
		}
		if lastLine == lineNum {
			r.Close()
			return sc.Text(), dateLine, lastLine, sc.Err()
		}
	}
	r.Close()
	return line, dateLine, lastLine, io.EOF
}

func lineCounter() (int, error) {
	r, err := os.Open("/tmp/EDO.log")
	if err != nil {
		log.Print(err)
	}
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			r.Close()
			return count, nil

		case err != nil:
			r.Close()
			return count, err
		}
	}
}

func dateConvert(date string) string {
	dtstr1 := date
	dt, _ := time.Parse("Mon Jan _2 15:04:05 MST 2006", dtstr1)
	dtstr2 := dt.Format("02/01/2006")
	return dtstr2
}
