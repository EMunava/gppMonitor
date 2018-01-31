package selenium

import (
	"fmt"
	"github.com/tebeka/selenium"
	"strings"
	"time"
)

func confirmWaitSchedSubBatch(wd selenium.WebDriver) {

	defer func() {
		if err := recover(); err != nil {
			img, _ := wd.Screenshot()
			sendError(fmt.Sprint(err), img, true)
			logOut(wd)
		}
	}()

	logIn(wd)

	navigateToSubBatchDates(wd)

	waitFor(wd, selenium.ByClassName, "ui-grid-cell-contents")

	subBatchAmount := extractSubBatchDates(wd)

	sendError(fmt.Sprint("Scheduled transactions: ", subBatchAmount), nil, false)

	logOut(wd)
}

func navigateToSubBatchDates(wd selenium.WebDriver) {

	byClassName(wd, "dh-navigation-tabs-current-tab-button")

	byCSSSelector(wd, "#main-content > div.dh-main-container.ng-scope > div > div > div:nth-child(2) > div.dh-navigation-tabs > div.dialer-container > ul > li:nth-child(1) > button")

	waitFor(wd, selenium.ByXPATH, "//*[contains(text(), 'Individual Messages (')]")

	byXPath(wd, "//*[contains(text(), 'Individual Messages (')]")

	byXPath(wd, "//*[contains(text(), 'Waiting')]")

	byXPath(wd, "//*[contains(text(), 'Wait Sched Sub Batch')]")

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
	sp, dValue := extract(date)

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
