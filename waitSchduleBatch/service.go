package waitSchduleBatch

import (
	"fmt"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/gppSelenium"
	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
	"github.com/zamedic/go2hal/alert"
	"strings"
	"time"
)

type Service interface {
	ConfirmWaitSchedSubBatch()
}

type service struct {
	selenium gppSelenium.Service
	alert    alert.Service
}

func NewService(alert alert.Service, selenium gppSelenium.Service) Service {
	return &service{alert: alert, selenium: selenium}
}

func (s *service) ConfirmWaitSchedSubBatch() {
	s.selenium.NewClient()
	defer s.selenium.Driver().Close()

	defer func() {
		if err := recover(); err != nil {
			s.selenium.HandleSeleniumError(true, errors.New(fmt.Sprint(err)))
			s.selenium.LogOut()
		}
	}()

	s.selenium.LogIn()

	s.navigateToSubBatchDates()

	s.selenium.WaitFor(selenium.ByClassName, "ui-grid-cell-contents")

	subBatchAmount := s.extractDates()

	s.selenium.HandleSeleniumError(false, fmt.Errorf("'Wait Scheduled Sub Batch' transaction total: %v", subBatchAmount))

	s.selenium.LogOut()
}

func (s *service) ConfirmPostingExecution() {
	s.selenium.NewClient()
	defer s.selenium.Driver().Close()

	defer func() {
		if err := recover(); err != nil {
			s.selenium.HandleSeleniumError(true, errors.New(fmt.Sprint(err)))
			s.selenium.LogOut()
		}
	}()

	s.selenium.LogIn()

	s.navigateToPostingExce()

	s.selenium.WaitFor(selenium.ByClassName, "ui-grid-cell-contents")

	postexAmount := s.extractDates()

	s.selenium.HandleSeleniumError(false, fmt.Errorf("'Posting Exception' amount: %v", postexAmount))

	s.selenium.LogOut()
}

func (s *service) navigateToSubBatchDates() {

	s.selenium.ClickByClassName("dh-navigation-tabs-current-tab-button")

	s.selenium.ClickByCSSSelector("#main-content > div.dh-main-container.ng-scope > div > div > div:nth-child(2) > div.dh-navigation-tabs > div.dialer-container > ul > li:nth-child(1) > button")

	s.selenium.WaitFor(selenium.ByXPATH, "//*[contains(text(), 'Individual Messages (')]")

	s.selenium.ClickByXPath("//*[contains(text(), 'Individual Messages (')]")

	s.selenium.ClickByXPath("//*[contains(text(), 'Waiting')]")

	s.selenium.ClickByXPath("//*[contains(text(), 'Wait Sched Sub Batch')]")
}

func (s *service) navigateToPostingExce() {

	s.selenium.ClickByClassName("dh-navigation-tabs-current-tab-button")

	s.selenium.ClickByCSSSelector("#main-content > div.dh-main-container.ng-scope > div > div > div:nth-child(2) > div.dh-navigation-tabs > div.dialer-container > ul > li:nth-child(1) > button")

	s.selenium.WaitFor(selenium.ByXPATH, "//*[contains(text(), 'Individual Messages (')]")

	s.selenium.ClickByXPath("//*[contains(text(), 'Individual Messages (')]")

	s.selenium.ClickByXPath("//*[contains(text(), 'Exception')]")

	s.selenium.ClickByXPath("//*[contains(text(), 'Posting Exception')]")
}

func (s *service) extractDates() int {

	Success := 0

	dates, err := s.selenium.Driver().FindElements(selenium.ByClassName, "ui-grid-cell-contents")
	if err != nil {
		panic(err)
	}

	for _, date := range dates {

		Success += s.extractionLoopSubBatch(date)
	}
	return Success
}

func (s *service) extractionLoopSubBatch(date selenium.WebElement) int {
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
