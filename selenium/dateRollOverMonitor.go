package selenium

import (
	"github.com/tebeka/selenium"
	"time"
	"log"
	"strings"
)

func confirmDateRollOver(wd selenium.WebDriver) {

	defer func() {
		if err := recover(); err != nil {
			logOut(wd)
		}
	}()

	//============================Login========+======================================
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

	user.SendKeys("A229343")
	pass.SendKeys("Gr33nfus")
	loginButton.Submit()

	//Wait for successful login
	waitFor(wd, "dh-customer-logo")

	//============================Confirm Dates=======================================

	waitForWaitFor(wd)

	bs, err := wd.FindElement(selenium.ByXPATH, "//*[contains(text(), 'Business Setup')]")
	if err != nil {
		panic(err)
	}

	if err = bs.Click(); err != nil {
		panic(err)
	}

	waitFor(wd, "ft-grid-click")

	//Values to check
	dates, err := wd.FindElements(selenium.ByClassName, "ui-grid-cell-contents")
	if err != nil {
		panic(err)
	}

	//Successful date rollover amount
	Success := 0

	for  _, date := range dates{

		dValue, err := date.GetAttribute("innerText")
		if err != nil {
			panic(err)
		}
		sp := strings.Split(dValue, "/")

		if len(sp) != 1 {
			success := dateConfirm(dValue)
			Success += success
		}
	}

	if Success == 2 {
		sendError("Dates successfully rolled over to next business day", nil, false )
	} else {
		img,_ := wd.Screenshot()
		sendError("One or more dates have not rolled over to next business day", img,false)
	}


	//============================Logout==============================================
	logOut(wd)
}

func dateConfirm(d1 string)(int) {

	currentDate := time.Now()
	cd := currentDate.Format("02/01/2006")
	c := strings.Compare(d1, cd)
	if c == 0{
		return 1
	} else {
		return 0
	}
}

func logOut(wd selenium.WebDriver){

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