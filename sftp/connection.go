package sftp

import (
	"os"
	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
	"time"
	"log"
	"strings"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

type File struct {
	Name         string
	Path         string
	Size         int64
	LastModified time.Time
}

func connect() (*sftp.Client, error) {

	der := decrypt([]byte(privateKey()), []byte(privatePass()))
	key, err := x509.ParsePKCS1PrivateKey(der)
	signer, err := ssh.NewSignerFromKey(key)
	if err != nil{
		log.Print(err)
	}

	//signer, err := ssh.ParsePrivateKey([]byte(privateKey()))
	//if err != nil{
	//	log.Print(err)
	//}

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

func decrypt(key []byte, password []byte) []byte {
	block, rest := pem.Decode(key)
	if len(rest) > 0 {
		log.Fatalf("Extra data included in key")
	}
	der, err := x509.DecryptPEMBlock(block, password)
	if err != nil {
		log.Fatalf("Decrypt failed: %v", err)
	}
	return der
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

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func retrieveFile(path, file string){

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

	// Create the destination file
	dstFile, err := os.Create(localPath + filename)
	if err != nil {
		log.Println(err)
	}
	defer dstFile.Close()
	// Copy the file
	srcFile.WriteTo(dstFile)
}





func sshUser() string {
	return os.Getenv("SSH_USER")
}
func privateKey() string {
	return strings.Replace(os.Getenv("SSH_KEY"), "*", "\n", -1)
}
func privatePass() string {
	return os.Getenv("DEC_PASS")
}
func sshEndpoint() string {
	return os.Getenv("SSH_ENDPOINT")
}
