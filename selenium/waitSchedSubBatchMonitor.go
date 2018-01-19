package selenium

import (
	"github.com/tebeka/selenium"
	"strings"
	"time"
)

func confirmWaitSchedSubBatch(wd selenium.WebDriver) {

	defer func() {
		if err := recover(); err != nil {
			logOut(wd)
		}
	}()

	logIn(wd)

	navigateToSubBatchDates()

	subBatchAmount := extractSubBatchDates(wd)



	logOut(wd)
}

func navigateToSubBatchDates() {


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