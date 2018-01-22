package selenium

import (
	"github.com/tebeka/selenium"
	"strings"
	"time"
	"fmt"
)

func confirmWaitSchedSubBatch(wd selenium.WebDriver) {

	defer func() {
		if err := recover(); err != nil {
			logOut(wd)
		}
	}()

	logIn(wd)

	navigateToSubBatchDates(wd)

	subBatchAmount := extractSubBatchDates(wd)

	sendError(fmt.Sprint("Scheduled transactions: %d", subBatchAmount), nil, false)

	logOut(wd)
}

func navigateToSubBatchDates(wd selenium.WebDriver) {

	grid, err := wd.FindElement(selenium.ByCSSSelector, "dh-navigation-tabs-current-tab-button")
	if err != nil {
		panic(err)
	}

	if err = grid.Click(); err != nil {
		panic(err)
	}

	sq, err := wd.FindElement(selenium.ByCSSSelector, "#main-content > div.dh-main-container.ng-scope > div > div > div:nth-child(2) > div.dh-navigation-tabs > div.dialer-container > ul > li:nth-child(1) > button")
	if err != nil {
		panic(err)
	}

	if err = sq.Click(); err != nil {
		panic(err)
	}

	im, err := wd.FindElement(selenium.ByXPATH, "//*[contains(text(), 'Individual Messages')]")
	if err != nil {
		panic(err)
	}

	if err = im.Click(); err != nil {
		panic(err)
	}

	waiting, err := wd.FindElement(selenium.ByXPATH, "//*[contains(text(), 'Waiting')]")
	if err != nil {
		panic(err)
	}

	if err = waiting.Click(); err != nil {
		panic(err)
	}

	waitSchedSubBatch, err := wd.FindElement(selenium.ByXPATH, "//*[contains(text(), 'Wait Sched Sub Batch')]")
	if err != nil {
		panic(err)
	}

	waitSchedSubBatch.Click()

}

func extractSubBatchDates(wd selenium.WebDriver) int {

	Success := 0

	dates, err := wd.FindElements(selenium.ByClassName, "ui-grid-cell-contents")
	if err != nil {
		panic(err)
	}

	for _, date := range dates {

		Success += extractionLoopSubBatch(date)
	}
	return Success
}

func extractionLoopSubBatch(date selenium.WebElement) int {
	dValue, err := date.GetAttribute("innerText")
	if err != nil {
		sendError(err.Error(), nil, true)
	}
	sp := strings.Split(dValue, "/")

	if len(sp) != 1 {
		success := dateConfirmSubBatch(dValue)
		return success
	}
	return 0
}

func dateConfirmSubBatch(d1 string) int {

	currentDate := time.Now()
	tomorrowDate := currentDate.AddDate(0, 0, 1)

	td := tomorrowDate.Format("02/01/2006")

		t := strings.Compare(d1, td)
		if t == 0 {
			return 1
		}
		return 0
}