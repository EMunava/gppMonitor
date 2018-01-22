package sftp

import (
	"fmt"
	"log"
)

func ListFiles() {

	fileList, err := GetFilesInPath("/cdwasha/connectdirect/outgoing/EDO_DirectDebitRequest")
	if err != nil {
		log.Print(err)
	}

	fmt.Print(fileList)
}
