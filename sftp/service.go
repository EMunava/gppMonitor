package sftp

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"

	"os"
	"strings"
	"time"
)

type Service interface {
	/*
		GetFilesInPath will return all the files within a path
	*/
	GetFilesInPath(path string) ([]File, error)

	RetrieveFile(path, file string)
}

type service struct {
}

func NewService() Service {
	return &service{}
}

type File struct {
	Name         string
	Path         string
	Size         int64
	LastModified time.Time
}

func (s *service) connect() (*sftp.Client, error) {

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

func (s *service) GetFilesInPath(path string) ([]File, error) {
	client, err := s.connect()
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

func (s *service) RetrieveFile(path, file string) {

	client, err := s.connect()
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

func sshUser() string {
	return os.Getenv("SSH_USER")
}

func privateKey() string {
	return strings.Replace(os.Getenv("SSH_KEY"), "*", "\n", -1)
}

func sshEndpoint() string {
	return os.Getenv("SSH_ENDPOINT")
}
