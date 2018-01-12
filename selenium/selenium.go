package selenium

import (
	"github.com/tebeka/selenium"
	"encoding/base64"
	"net/http"
	"encoding/json"
	"bytes"
	"log"
	"io/ioutil"
	"runtime/debug"
	"os"
	"github.com/jasonlvhit/gocron"
	"fmt"
)

type alertMessage struct {
	Message, Image string
	internalError bool
}

func init(){
	go func() {
		schedule()
	}()

}

func schedule(){
	gocron.Every(1).Day().At("11:30").Do(doSelenium)
	gocron.Every(1).Day().At("00:30").Do(doSelenium)
	gocron.Every(1).Day().At("01:30").Do(doSelenium)
	_, schedule := gocron.NextRun()
	fmt.Println(schedule)

	<- gocron.Start()
}

func doSelenium() {
	var webDriver selenium.WebDriver
	var err error
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "chrome"})
	caps["chrome.switches"] = []string{"--ignore-certificate-errors"}

	if webDriver, err = selenium.NewRemote(caps, seleniumServer()); err != nil {
		handleSeleniumError(err, nil)
		return
	}

	defer webDriver.Quit()

	err = webDriver.Get(endpoint())
	if err != nil {
		handleSeleniumError(err, webDriver)
	}

	if err != nil {
		handleSeleniumError(err, webDriver)
		return
	}

	confirmDateRollOver(webDriver)

}

func waitForWaitFor(webDriver selenium.WebDriver) error {
	return webDriver.Wait(func(wb selenium.WebDriver) (bool, error) {
		elem, err := wb.FindElement(selenium.ByCSSSelector, "body > div.dh-notification.ng-scope.success > div")
		if err != nil {
			return true, nil
		}
		r, err := elem.IsDisplayed()
		return !r, nil
	})
}

func handleSeleniumError(err error, driver selenium.WebDriver) {
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
	a := alertMessage{Message: message, internalError:internalError}
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


