package sftp

import (
	"log"
	"fmt"
)

func RetreiveEDOLog() {

	fileList, err := GetFilesInPath("/cdwasha/connectdirect/outgoing/EDO_DirectDebitRequest")
	if err != nil {
		log.Print(err)
	}

	fmt.Print(fileList)
}
