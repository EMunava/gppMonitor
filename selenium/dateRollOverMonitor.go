package selenium

import (
	"github.com/tebeka/selenium"
	"time"
	"fmt"
	"log"
	"strings"
	"image"
	"bytes"
	"os"
	"image/png"
)

func confirmDateRollOver(wd selenium.WebDriver) {


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

	loginButton, err := wd.FindElement(selenium.ByCSSSelector, "#main-content > div > div > div.ft-form-container > div:nth-child(2) > div > form > button")

	user.SendKeys("A229343")
	pass.SendKeys("Gr33nfus")
	loginButton.Submit()

	//Wait for successful login
	wd.Wait(func(wb selenium.WebDriver) (bool, error) {
		elem, err := wb.FindElement(selenium.ByClassName, "dh-customer-logo")
		if err != nil {
			return false, nil
		}
		r, err := elem.IsDisplayed()
		return r, nil
	})

	//============================Confirm Dates=======================================

	wd.SetImplicitWaitTimeout(5000)

	bs, err := wd.FindElement(selenium.ByCSSSelector, "#main-content > div.ft-pagination-container.ng-scope > div > div > div.dh-navigation-panel > ul > li:nth-child(4) > button")
	if err != nil {
		log.Println(err.Error())
		logOut(wd)
		panic(err)
	}
	err = bs.Click()

	wd.SetImplicitWaitTimeout(5 * time.Second)
	img,_ := wd.Screenshot()
	saveImage(img)
	sendError("Here be your pic",img,false)

	//test, err := wd.FindElement(selenium.ByCSSSelector, "#main-content > div.dh-main-container.ng-scope > div > div > div.dh-main-right-container.ng-scope > div > div.ft-top-grid-action > div.pull-left > div > span:nth-child(2)")
	//if err != nil {
	//	log.Println("Element NOT found")
	//	logOut(wd)
	//	panic(err)
	//}
	//log.Println("Element found")
	//err = test.Click()

	//Values to check
	date1, err := wd.FindElement(selenium.ByCSSSelector, "#\\31 515146185525-0-uiGrid-000D-cell > div")
	if err != nil {
		log.Println(err.Error())
		logOut(wd)
		panic(err)
	}
	d1, err := date1.GetAttribute("innerText")

	date2, err := wd.FindElement(selenium.ByXPATH, "//*[@id='1515064846508-1-uiGrid-000H-cell']")
	if err != nil {
		log.Println(err.Error())
		logOut(wd)
		panic(err)
	}
	d2, err := date2.GetAttribute("innerText")

	fmt.Println(d1)
	fmt.Println(d2)

	dateConfirm(d1, d2)

	//=======================Log Out==================================================

	time.Sleep(1 * time.Minute)
	logOut(wd)
}

func dateConfirm(d1, d2 string) {

	//d1 := "06/01/2018"
	// d2 := "07/01/2018"

	currentDate := time.Now()
	cd := currentDate.Format("02/01/2006")

	c := strings.Compare(d1, cd)
	if c == 0 {
		fmt.Println("Dates correlate")
	} else {
		fmt.Println("Dates do not correlate")
	}
}

func logOut(wd selenium.WebDriver){

	signOutButton, err := wd.FindElement(selenium.ByCSSSelector, "#main-content > div.ft-pagination-container.ng-scope > div > div > div.dh-top-panel > div.dh-top-panel-menu > div.dh-menu-buttons-wrapper > button:nth-child(2)")
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
	err = signOutButton.Click()
}

func saveImage(img []byte){

	image, _, _ := image.Decode(bytes.NewReader(img))
	out, err := os.Create("/home/jurgen/Desktop/test/QRImg.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = png.Encode(out, image)

}
