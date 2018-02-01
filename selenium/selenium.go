package selenium

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/jasonlvhit/gocron"
	"github.com/tebeka/selenium"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"database/sql/driver"
)

type Service interface{

}

type chromeService struct{
	driver selenium.WebDriver
}




func handleSeleniumError(err error) {
	debug.PrintStack()
	if driver == nil {
		sendError(err.Error(), nil, true)
		return
	}
	bytes, error := driver.Screenshot()
	if error != nil {
		// Couldnt get a screenshot - lets end the original error
		sendError(err.Error(), nil, true)
		return
	}
	sendError(err.Error(), bytes, true)
}

func sendError(message string, image []byte, internalError bool) {
	a := alertMessage{Message: message, InternalError: internalError}
	if image != nil {
		a.Image = base64.StdEncoding.EncodeToString(image)
	}

	request, _ := json.Marshal(a)

	response, err := http.Post(errorEndpoint(), "application/json", bytes.NewReader(request))
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Println(ioutil.ReadAll(response.Body))
	}

}

func endpoint() string {
	return os.Getenv("GPP_ENDPOINT")
}

func seleniumServer() string {
	return os.Getenv("SELENIUM_SERVER")
}

func errorEndpoint() string {
	return os.Getenv("HAL_ENDPOINT")
}
