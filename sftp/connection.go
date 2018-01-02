package sftp

import (
	"os"
	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
	"time"
)

type File struct {
	Name         string
	Path         string
	Size         int64
	LastModified time.Time
}

func connect() (*sftp.Client, error) {
	signer, _ := ssh.ParsePrivateKey([]byte(privateKey()))
	clientConfig := &ssh.ClientConfig{
		User: sshUser(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
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

func sshUser() string {
	return os.Getenv("SSH_USER")
}
func privateKey() string {
	return os.Getenv("SSH_KEY")
}
func sshEndpoint() string {
	return os.Getenv("SSH_ENDPOINT")
}
