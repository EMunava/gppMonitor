package waitSchduleBatch

import (
	"fmt"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/gppSelenium"
	"github.com/matryer/try"
	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
	"github.com/zamedic/go2hal/alert"
	"log"
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

func (s *service) ConfirmWaitSchedSubBatchMethod() (r error) {
	s.selenium.NewClient()
	defer s.selenium.Driver().Close()

	defer func() {
		if err := recover(); err != nil {
			s.selenium.HandleSeleniumError(true, errors.New(fmt.Sprint(err)))
			s.selenium.LogOut()
			r = errors.New("Wait Scheduled Sub Batch transaction amount retrieval failed")
		}
	}()

	s.selenium.LogIn()

	s.navigateToBatchDates()

	s.selenium.ClickByXPath("//*[contains(text(), 'Wait Sched Sub Batch')]")

	s.selenium.WaitFor(selenium.ByClassName, "ui-grid-cell-contents")

	subBatchAmount := s.extractBatchDates(true)

	s.selenium.ClickByXPath("//*[contains(text(), 'Wait Posting')]")

	s.selenium.WaitFor(selenium.ByClassName, "ui-grid-cell-contents")

	waitPostingAmount := s.extractBatchDates(false)

	s.selenium.HandleSeleniumError(false, fmt.Errorf("Transactions in tracking(Posting): %v \nNew transactions to be processed(Scheduled Sub Batch): %v", waitPostingAmount, subBatchAmount))

	log.Printf("Transactions in Tracking: %v \nNew Transactions: %v", waitPostingAmount, subBatchAmount)

	s.selenium.LogOut()

	return nil
}

func (s *service) navigateToBatchDates() {

	s.selenium.ClickByClassName("dh-navigation-tabs-current-tab-button")

	s.selenium.ClickByCSSSelector("#main-content > div.dh-main-container.ng-scope > div > div > div:nth-child(2) > div.dh-navigation-tabs > div.dialer-container > ul > li:nth-child(1) > button")

	s.selenium.WaitFor(selenium.ByXPATH, "//*[contains(text(), 'Individual Messages (')]")

	s.selenium.ClickByXPath("//*[contains(text(), 'Individual Messages (')]")

	s.selenium.ClickByXPath("//*[contains(text(), 'Waiting')]")
}

func (s *service) extractBatchDates(isFollowingDay bool) int {

	Success := 0

	dates, err := s.selenium.Driver().FindElements(selenium.ByClassName, "ui-grid-cell-contents")
	if err != nil {
		panic(err)
	}

	for _, date := range dates {

		Success += s.extractionLoopBatch(date, isFollowingDay)
	}
	return Success
}

func (s *service) extractionLoopBatch(date selenium.WebElement, isFollowingDay bool) int {
	sp, dValue := s.extract(date)

	if len(sp) != 1 {
		success := dateConfirmSubBatch(dValue, isFollowingDay)
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

func dateConfirmSubBatch(d1 string, isFollowingDay bool) int {

	currentDate := time.Now()
	tomorrowDate := currentDate.AddDate(0, 0, 1)

	td := tomorrowDate.Format("02/01/2006")

	if isFollowingDay == true {
		t := strings.Compare(d1, td)
		if t == 0 {
			return 1
		}
		return 0
	}
	return 1
}

func (s *service) ConfirmWaitSchedSubBatch() {
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		err = s.ConfirmWaitSchedSubBatchMethod()
		if err != nil {
			log.Println("next attempt in 2 minutes")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return attempt < 5, err //5 attempts
	})
	if err != nil {
		log.Fatalln(err)
	}
}
