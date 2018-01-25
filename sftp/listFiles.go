package sftp

import (
	"fmt"
	"log"
)

/*
ListFiles lists all files in target directory and stores file info within struct for future use
*/
func ListFiles() {

	fileList, err := GetFilesInPath("/cdwasha/connectdirect/outgoing/EDO_DirectDebitRequest")
	if err != nil {
		log.Print(err)
	}

	fmt.Print(fileList)
}
