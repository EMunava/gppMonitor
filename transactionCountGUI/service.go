package transactionCountGUI

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/kyokomi/emoji"
	"github.com/matryer/try"
	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/gppMonitor/gppSelenium"
	"log"
	"regexp"
	"strconv"
	"time"
)

type Service interface {
	ExtractTransactionCount()
}

type service struct {
	selenium gppSelenium.Service
	alert    alert.Service
}

type transactionCount struct {
	waitSchedSubBatch int
	waitPosting       int
	scheduled         int
}

func NewService(alert alert.Service, selenium gppSelenium.Service) Service {
	s := &service{alert: alert, selenium: selenium}
	go func() {
		s.schedule()
	}()
	return s
}

func (s *service) schedule() {
	confirmWaitSched := gocron.NewScheduler()

	go func() {
		confirmWaitSched.Every(1).Day().At("19:00").Do(s.ExtractTransactionCount)
		<-confirmWaitSched.Start()
	}()
}

func (s *service) ExtractTransactionCountMethod() (r error) {
	s.selenium.NewClient()
	defer s.selenium.Driver().Quit()

	defer func() {
		if err := recover(); err != nil {
			s.selenium.HandleSeleniumError(true, errors.New(fmt.Sprint(err)))
			s.selenium.LogOut()
			if e, ok := err.(error); ok {
				r = errors.New(e.Error())
			}
			r = errors.New("Wait Scheduled Sub Batch transaction amount retrieval failed")
		}
	}()
	
	transactions := transactionCount{}

	s.selenium.LogIn()

	s.navigateToIndividualMessages()

	wh, err := s.selenium.Driver().FindElements(selenium.ByXPATH, "//*[contains(text(), 'Warehouse')]")

	transactions.scheduled = s.extractInteger(s.extractString(wh[0]))

	s.selenium.ClickByXPath("//*[contains(text(), 'Waiting')]")

	wp, err := s.selenium.Driver().FindElements(selenium.ByXPATH, "//*[contains(text(), 'Wait Posting')]")
	if err != nil {
		panic(err)
	}

	sb, err := s.selenium.Driver().FindElements(selenium.ByXPATH, "//*[contains(text(), 'Wait Sched Sub Batch')]")
	if err != nil {
		panic(err)
	}

	transactions.waitPosting = s.extractInteger(s.extractString(wp[0]))

	transactions.waitSchedSubBatch = s.extractInteger(s.extractString(sb[0]))

	s.selenium.HandleSeleniumError(false, fmt.Errorf(emoji.Sprintf(":white_check_mark: Transactions in Wait Posting: %v\nTransactions in Wait Scheduled Sub Batch: %v\nTransactions in Warehouse Scheduled: %v", transactions.waitPosting, transactions.waitSchedSubBatch, transactions.scheduled)))

	log.Printf("Transactions in Wait Posting: %v\nTransactions in Wait Scheduled Sub Batch: %v\nTransactions in Warehouse Scheduled: %v", transactions.waitPosting, transactions.waitSchedSubBatch, transactions.scheduled)

	s.selenium.LogOut()

	return nil
}

func (s *service) navigateToIndividualMessages() {

	s.selenium.ClickByClassName("dh-navigation-tabs-current-tab-button")

	s.selenium.ClickByCSSSelector("#main-content > div.dh-main-container.ng-scope > div > div > div:nth-child(2) > div.dh-navigation-tabs > div.dialer-container > ul > li:nth-child(1) > button")

	s.selenium.WaitFor(selenium.ByXPATH, "//*[contains(text(), 'Individual Messages (')]")

	s.selenium.ClickByXPath("//*[contains(text(), 'Individual Messages (')]")

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

func (s *service) ExtractTransactionCount() {
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		err = s.ExtractTransactionCountMethod()
		if err != nil {
			log.Println("next attempt in 2 minutes")
			time.Sleep(2 * time.Minute) // wait 2 minutes
		}
		return attempt < 5, err //5 attempts
	})
	if err != nil {
		log.Println(err)
	}
}
