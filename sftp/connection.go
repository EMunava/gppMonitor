package sftp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	go func() {
		schedule()
	}()

}

func schedule() {
	gocron.Every(1).Day().At("01:10").Do(RetrieveEDOLog)

	_, schedule := gocron.NextRun()
	fmt.Println(schedule)

	<-gocron.Start()
}

type alert struct {
	Message string
}

type File struct {
	Name         string
	Path         string
	Size         int64
	LastModified time.Time
}

func connect() (*sftp.Client, error) {

	signer, err := ssh.ParsePrivateKey([]byte(privateKey()))
	if err != nil {
		log.Print(err)
	}
	clientConfig := &ssh.ClientConfig{
		User: sshUser(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshClient, err := ssh.Dial("tcp", sshEndpoint(), clientConfig)
	if err != nil {
		return nil, err
	}
	client, err := sftp.NewClient(sshClient)
	return client, err
}

/*
GetFilesInPath will return all the files within a path
*/
func GetFilesInPath(path string) ([]File, error) {
	client, err := connect()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	files, err := client.ReadDir(path)
	if err != nil {
		return nil, err
	}

	result := make([]File, len(files))

	for x, file := range files {
		result[x] = File{Name: file.Name(), LastModified: file.ModTime(), Path: path, Size: file.Size()}
	}
	return result, nil
}

func retrieveFile(path, file string) {

	client, err := connect()
	if err != nil {
		log.Println(err)
	}
	defer client.Close()

	srcPath := path
	localPath := "/tmp/"
	filename := file

	srcFile, err := client.Open(srcPath + filename)
	if err != nil {
		log.Println(err)
	}
	defer srcFile.Close()
	dstFile, err := os.Create(localPath + filename)
	if err != nil {
		log.Println(err)
	}
	defer dstFile.Close()
	srcFile.WriteTo(dstFile)
}

func sendAlert(message string) {
	a := alert{Message: message}

	request, _ := json.Marshal(a)

	response, err := http.Post(alertEndpoint(), "application/json", bytes.NewReader(request))
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Println(ioutil.ReadAll(response.Body))
	}

}

func sshUser() string {
	return os.Getenv("SSH_USER")
}

func privateKey() string {
	return strings.Replace(os.Getenv("SSH_KEY"), "*", "\n", -1)
}

func sshEndpoint() string {
	return os.Getenv("SSH_ENDPOINT")
}

func alertEndpoint() string {
	return os.Getenv("HAL_ENDPOINT_ALERT")
}
