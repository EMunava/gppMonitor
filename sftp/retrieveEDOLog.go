package sftp

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strings"
)

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
		sendError("EDO Posing request file successfully sent", nil, false)
		log.Println("EDO Posing request file successfully sent")
	} else {
		sendError("EDO Posing request file send failed!!", nil, false)
		log.Println("EDO Posing request file send failed!!")
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
			//r.Close()
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
			count -= 1
			r.Close()
			return count, nil

		case err != nil:
			count -= 1
			r.Close()
			return count, err
		}
	}
}
