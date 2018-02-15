package waitSchduleBatch

import (
	"fmt"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/gppSelenium"
	"github.com/matryer/try"
	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
	"github.com/zamedic/go2hal/alert"
	"log"
	"regexp"
	"strconv"
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
	defer s.selenium.Driver().Quit()

	defer func() {
		if err := recover(); err != nil {
			s.selenium.HandleSeleniumError(true, errors.New(fmt.Sprint(err)))
			s.selenium.LogOut()
			r = errors.New("Wait Scheduled Sub Batch transaction amount retrieval failed")
		}
	}()

	s.selenium.LogIn()

	s.navigateToBatchDates()

	wp, err := s.selenium.Driver().FindElements(selenium.ByXPATH, "//*[contains(text(), 'Wait Posting')]")
	if err != nil {
		panic(err)
	}

	sb, err := s.selenium.Driver().FindElements(selenium.ByXPATH, "//*[contains(text(), 'Wait Sched Sub Batch')]")
	if err != nil {
		panic(err)
	}

	waitPostingAmount := s.extractInteger(s.extractString(wp[0]))

	subBatchAmount := s.extractInteger(s.extractString(sb[0]))

	s.selenium.HandleSeleniumError(false, fmt.Errorf("Transactions in tracking(Posting): %d New transactions to be processed(Scheduled Sub Batch): %d", waitPostingAmount, subBatchAmount))

	log.Printf("Transactions in Tracking: %v\nNew Transactions: %v", waitPostingAmount, subBatchAmount)

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

func (s *service) extractInteger(i string) int {
	re := regexp.MustCompile("[0-9]+")
	ar := re.FindAllString(i, -1)
	s2i, err := strconv.Atoi(ar[0])
	if err != nil {
		s.selenium.HandleSeleniumError(true, err)
	}
	return s2i
}

func (s *service) extractString(date selenium.WebElement) string {
	str, err := date.GetAttribute("innerText")
	if err != nil {
		s.selenium.HandleSeleniumError(true, err)
	}
	return str
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
