package sftp

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strings"
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
	ll, _, err := lastLine(numLines)
	if err != nil {
		log.Print(err)
	}

	if strings.Contains(ll, "successful") {
		sendAlert("EDO Posting request file successfully sent")
		log.Println("EDO Posting request file successfully sent")
	} else if strings.Contains(ll, "failed") {
		sendAlert("EDO Posting request file send failed")
		log.Println("EDO Posting request file send failed!!")
	} else {
		log.Println("The last line did not contain success/failiure information")
	}

	os.Remove("/tmp/EDO.log")
}

func lastLine(lineNum int) (line string, lastLine int, err error) {
	r, err := os.Open("/tmp/EDO.log")
	if err != nil {
		log.Print(err)
	}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			r.Close()
			return sc.Text(), lastLine, sc.Err()
		}
	}
	r.Close()
	return line, lastLine, io.EOF
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
