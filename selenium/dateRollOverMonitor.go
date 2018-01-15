package selenium

import (
	"github.com/tebeka/selenium"
	"log"
	"strings"
	"time"
)

func confirmDateRollOver(wd selenium.WebDriver) {

	defer func() {
		if err := recover(); err != nil {
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
	user, err := wd.FindElement(selenium.ByName, "txtUserId")
	if err != nil {
		panic(err)
	}

	if err := user.Clear(); err != nil {
		panic(err)
	}

	pass, err := wd.FindElement(selenium.ByName, "txtPassword")
	if err != nil {
		panic(err)
	}

	if err := pass.Clear(); err != nil {
		panic(err)
	}

	loginButton, err := wd.FindElement(selenium.ByXPATH, "//*[contains(text(), 'Sign In')]")

	user.SendKeys("")
	pass.SendKeys("")
	loginButton.Submit()

	//Wait for successful login
	waitFor(wd, "dh-customer-logo")
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

	waitFor(wd, "ft-grid-click")
}

func logOut(wd selenium.WebDriver) {

	signOutButton, err := wd.FindElement(selenium.ByXPATH, "//*[contains(text(), 'Sign Out')]")
	if err != nil {
		log.Println(err.Error())
	}
	err = signOutButton.Click()
}

func waitFor(webDriver selenium.WebDriver, selector string) {

	webDriver.Wait(func(wb selenium.WebDriver) (bool, error) {

		elem, err := wb.FindElement(selenium.ByClassName, selector)
		if err != nil {
			return false, nil
		}
		r, err := elem.IsDisplayed()
		return r, nil
	})
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
	cd := currentDate.Format("02/01/2006")
	c := strings.Compare(d1, cd)
	if c == 0 {
		return 1
	}
	return 0

}
