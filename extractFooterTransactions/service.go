package extractFooterTransactions

import (
	"bufio"
	"errors"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/sftp"
	"github.com/zamedic/go2hal/alert"
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
	return &service{sftpService: sftpService, alertService: alertService}
}

func (s *service) retreiveTransactions(contains string, exclude ...string) {

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
}

func (s *service) RetrieveSAPTransactions() {
	s.retreiveTransactions("RESPONSE.SAP")

}

func (s *service) RetrieveLEGTransactions() {
	s.retreiveTransactions("RESPONSE.LEG", "SAP")
}

func (s *service) RetrieveLEGSAPTransactions() {
	s.retreiveTransactions("RESPONSE.LEG.SAP")
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
	return "", "", errors.New("File has not arrived yet")
}
