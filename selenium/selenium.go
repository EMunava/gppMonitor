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
)

type alertMessage struct {
	Message, Image string
	InternalError  bool
}

func init() {
	go func() {
		schedule()
	}()

}

//CallSeleniumDateCheck confirms the GPP transaction date rollover buy comparing the current/tommorow's date and the date logged.
func CallSeleniumDateCheck() {
	seleniumDateRolloverCheck()
}

func CallWaitSchedBatchCheck() {
	seleniumWaitSchedBatchCheck()
}

func schedule() {
	sel := gocron.NewScheduler()
	sel.Every(1).Day().At("23:30").Do(seleniumDateRolloverCheck)
	sel.Every(1).Day().At("00:30").Do(seleniumDateRolloverCheck)
	sel.Every(1).Day().At("01:30").Do(seleniumDateRolloverCheck)
	sel.Every(1).Day().At("00:35").Do(seleniumWaitSchedBatchCheck)
	_, schedule := gocron.NextRun()
	log.Println(schedule)

	<-sel.Start()
}

func seleniumDateRolloverCheck() {
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
		return
	}

	confirmDateRollOver(webDriver)

}

func seleniumWaitSchedBatchCheck() {
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
		return
	}

	confirmWaitSchedSubBatch(webDriver)

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

func waitFor(webDriver selenium.WebDriver, findBy, selector string) error {

	e := webDriver.Wait(func(wb selenium.WebDriver) (bool, error) {

		elem, err := wb.FindElement(findBy, selector)
		if err != nil {
			return false, nil
		}
		r, err := elem.IsDisplayed()
		return r, nil
	})
	return e
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

func byClassName(wd selenium.WebDriver, cn string) {
	item, err := wd.FindElement(selenium.ByClassName, cn)
	if err != nil {
		panic(err)
	}

	if err = item.Click(); err != nil {
		panic(err)
	}
}

func byXPath(wd selenium.WebDriver, xp string) {
	item, err := wd.FindElement(selenium.ByXPATH, xp)
	if err != nil {
		panic(err)
	}

	if err = item.Click(); err != nil {
		panic(err)
	}
}

func byCSSSelector(wd selenium.WebDriver, cs string) {
	item, err := wd.FindElement(selenium.ByCSSSelector, cs)
	if err != nil {
		panic(err)
	}

	if err = item.Click(); err != nil {
		panic(err)
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
