package selenium

import (
	"fmt"
	"github.com/tebeka/selenium"
	"log"
	"os"
	"strings"
	"time"
)

func confirmDateRollOver(wd selenium.WebDriver) {

	defer func() {
		if err := recover(); err != nil {
			img, _ := wd.Screenshot()
			sendError(fmt.Sprint(err), img, true)
			logOut(wd)
		}
	}()

	logIn(wd)

	navigateToDates(wd)

	Success := extractDates(wd)

	if Success == 2 {
		sendError("Dates successfully rolled over to next business day", nil, false)
	} else {
		img, _ := wd.Screenshot()
		sendError("One or more dates have not rolled over to next business day", img, false)
	}

	logOut(wd)
}

func logIn(wd selenium.WebDriver) {

	if err := waitFor(wd, "dh-input-field"); err != nil {
		panic(err)
	}

	user, err := wd.FindElement(selenium.ByName, "txtUserId")
	if err != nil {
		panic(err)
	}

	pass, err := wd.FindElement(selenium.ByName, "txtPassword")
	if err != nil {
		panic(err)
	}

	loginButton, err := wd.FindElement(selenium.ByXPATH, "//*[contains(text(), 'Sign In')]")

	user.SendKeys(gppUser())
	pass.SendKeys(gppPass())
	loginButton.Submit()

	//Wait for successful login
	if err := waitFor(wd, "dh-customer-logo"); err != nil {
		panic(err)
	}
}

func navigateToDates(wd selenium.WebDriver) {

	waitForWaitFor(wd)

	bs, err := wd.FindElement(selenium.ByXPATH, "//*[contains(text(), 'Business Setup')]")
	if err != nil {
		panic(err)
	}

	if err = bs.Click(); err != nil {
		panic(err)
	}

	if err := waitFor(wd, "ft-grid-click"); err != nil {
		panic(err)
	}
}

func logOut(wd selenium.WebDriver) {

	signOutButton, err := wd.FindElement(selenium.ByXPATH, "//*[contains(text(), 'Sign Out')]")
	if err != nil {
		img, _ := wd.Screenshot()
		sendError(fmt.Sprint(err), img, true)
		log.Println(err.Error())
		return
	}
	err = signOutButton.Click()
	if err != nil {
		img, _ := wd.Screenshot()
		sendError(fmt.Sprint(err), img, true)
		log.Println(err.Error())
		return
	}

	if err := waitFor(wd, "dh-input-field"); err != nil {
		log.Println(err.Error())
		return
	}
}

func waitFor(webDriver selenium.WebDriver, selector string) error {

	e := webDriver.Wait(func(wb selenium.WebDriver) (bool, error) {

		elem, err := wb.FindElement(selenium.ByClassName, selector)
		if err != nil {
			return false, nil
		}
		r, err := elem.IsDisplayed()
		return r, nil
	})
	return e
}

func extractDates(wd selenium.WebDriver) int {

	Success := 0

	dates, err := wd.FindElements(selenium.ByClassName, "ui-grid-cell-contents")
	if err != nil {
		panic(err)
	}

	for _, date := range dates {

		Success += extractionLoop(date)
	}
	return Success
}

func extractionLoop(date selenium.WebElement) int {
	dValue, err := date.GetAttribute("innerText")
	if err != nil {
		sendError(err.Error(), nil, true)
	}
	sp := strings.Split(dValue, "/")

	if len(sp) != 1 {
		success := dateConfirm(dValue)
		return success
	}
	return 0
}

func dateConfirm(d1 string) int {

	currentDate := time.Now()
	tomorrowDate := currentDate.AddDate(0, 0, 1)

	cd := currentDate.Format("02/01/2006")
	td := tomorrowDate.Format("02/01/2006")

	if currentDate.Before(time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 24, 0, 0, 0, currentDate.Location())) && currentDate.After(time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 23, 0, 0, 0, currentDate.Location())) {
		t := strings.Compare(d1, td)
		if t == 0 {
			return 1
		}
		return 0
	}
	c := strings.Compare(d1, cd)
	if c == 0 {
		return 1
	}
	return 0
}

func gppUser() string {
	return os.Getenv("GPP_USER")
}

func gppPass() string {
	return os.Getenv("GPP_PASSWORD")
}
