package waitSchduleBatch



import (
	"fmt"
	"github.com/tebeka/selenium"
	"strings"
	"time"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/gppSelenium"
	"github.com/zamedic/go2hal/alert"
	"github.com/pkg/errors"
)

type Service interface {
	ConfirmWaitSchedSubBatch()
}

type service struct {
	selenium gppSelenium.Service
	alert alert.Service
}

func NewService(alert alert.Service) Service {
	return &service{alert:alert}
}

func (s *service)ConfirmWaitSchedSubBatch() {

	s.selenium = gppSelenium.NewService(s.alert)

	defer s.selenium.Driver().Close()

	wd := s.selenium.Driver()

	defer func() {
		if err := recover(); err != nil {
			s.selenium.HandleSeleniumError(true,errors.New(fmt.Sprint(err)))
			s.selenium.LogOut()
		}
	}()

	s.selenium.LogIn()

	s.navigateToSubBatchDates(wd)

	s.selenium.WaitFor(selenium.ByClassName, "ui-grid-cell-contents")

	subBatchAmount := s.extractSubBatchDates(wd)

	s.selenium.HandleSeleniumError(false,fmt.Errorf("Scheduled transactions: ", subBatchAmount))

	s.selenium.LogOut()
}

func (s *service)navigateToSubBatchDates(wd selenium.WebDriver) {

	s.selenium.ClickByClassName("dh-navigation-tabs-current-tab-button")

	s.selenium.ClickByCSSSelector( "#main-content > div.dh-main-container.ng-scope > div > div > div:nth-child(2) > div.dh-navigation-tabs > div.dialer-container > ul > li:nth-child(1) > button")

	s.selenium.WaitFor( selenium.ByXPATH, "//*[contains(text(), 'Individual Messages (')]")

	s.selenium.ClickByXPath( "//*[contains(text(), 'Individual Messages (')]")

	s.selenium.ClickByXPath( "//*[contains(text(), 'Waiting')]")

	s.selenium.ClickByXPath( "//*[contains(text(), 'Wait Sched Sub Batch')]")

}

func (s *service) extractSubBatchDates(wd selenium.WebDriver) int {

	Success := 0

	dates, err := wd.FindElements(selenium.ByClassName, "ui-grid-cell-contents")
	if err != nil {
		panic(err)
	}

	for _, date := range dates {

		Success += s.extractionLoopSubBatch(date)
	}
	return Success
}

func (s *service)extractionLoopSubBatch(date selenium.WebElement) int {
	sp, dValue := s.extract(date)

	if len(sp) != 1 {
		success := dateConfirmSubBatch(dValue)
		return success
	}
	return 0
}

func (s *service) extract(date selenium.WebElement) ([]string, string) {
	dValue, err := date.GetAttribute("innerText")
	if err != nil {
		s.selenium.HandleSeleniumError(true, err)
	}
	sp := strings.Split(dValue, "/")
	return sp, dValue
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
