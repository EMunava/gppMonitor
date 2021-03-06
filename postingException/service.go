package postingException

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/gppMonitor/gppSelenium"
	"log"
	"strings"
	"time"
)

//Service interface exposes the ConfirmPostingException method
type Service interface {
	ConfirmPostingException()
}

type service struct {
	selenium       gppSelenium.Service
	alert          alert.Service
	previousPostEx postExInfo
}

type postExInfo struct {
	MIDList []string
	Amount  int
}

//NewService function creates instances of required external service structs for local use
func NewService(alert alert.Service, selenium gppSelenium.Service) Service {
	s := &service{alert: alert, selenium: selenium}
	go func() {
		s.schedule()
	}()
	return s
}

func (s *service) schedule() {
	postex := gocron.NewScheduler()
	postexConst := gocron.NewScheduler()

	go func() {
		postex.Every(1).Day().At("08:00").Do(s.ConfirmPostingException)
		<-postex.Start()
	}()

	go func() {
		postexConst.Every(1).Hour().Do(s.ConfirmPostingException)
		<-postexConst.Start()
	}()
}

func (s *service) ConfirmPostingException() {
	s.selenium.NewClient()

	defer s.selenium.Driver().Close()

	defer func() {
		if err := recover(); err != nil {
			s.selenium.HandleSeleniumError(true, errors.New(fmt.Sprint(err)))
			s.selenium.LogOut()
		}
	}()

	currentDate := time.Now()
	if currentDate.Before(time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 8, 0, 0, 0, currentDate.Location())) && currentDate.After(time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 6, 0, 0, 0, currentDate.Location())) {

		s.previousPostEx.resetPrevious()
	}

	s.selenium.LogIn()

	s.navigateToPostingExce()

	s.selenium.WaitFor(selenium.ByCSSSelector, "#main-content > div.dh-main-container.ng-scope > div > div > div.dh-main-right-container.ng-scope > div > div > div > div > div > div.ft-top-grid-action > div.pull-left > div.top-grid-action-section-title > span")

	postEx := s.extractPostEx()

	postExCheck := s.compareCurrentAndPastPostex(postEx, s.previousPostEx)

	if !postExCheck && postEx.Amount != 0 {
		s.selenium.HandleSeleniumError(false, fmt.Errorf("Posting Exception count: %d for %v", postEx.Amount, time.Now().Format("02/01/2006")))

		s.previousPostEx.setPrevious(&postEx)
	} else {
		log.Println("No changes")
	}

	s.selenium.LogOut()
}

func (s *service) navigateToPostingExce() {

	s.selenium.ClickByClassName("dh-navigation-tabs-current-tab-button")

	//s.selenium.ClickByCSSSelector("#main-content > div.dh-main-container.ng-scope > div > div > div:nth-child(2) > div.dh-navigation-tabs > div.dialer-container > ul > li:nth-child(1) > button")
	s.selenium.ClickByCSSSelector("#main-content > div.dh-main-container.ng-scope > div > div > div.dh-main-left-container-wrapper > div:nth-child(2) > div.dh-navigation-tabs > div.dialer-container > ul > li:nth-child(1) > button")

	s.selenium.WaitFor(selenium.ByXPATH, "//*[contains(text(), 'Individual Messages (')]")

	s.selenium.ClickByXPath("//*[contains(text(), 'Individual Messages (')]")

	s.selenium.ClickByXPath("//*[contains(text(), 'Exception')]")

	s.selenium.ClickByXPath("//*[contains(text(), 'Posting Exception')]")
}

func (s *service) extractPostEx() postExInfo {

	postEx := postExInfo{}

	mids, err := s.selenium.Driver().FindElements(selenium.ByClassName, "ui-grid-cell-contents")
	if err != nil {
		panic(err)
	}

	for _, mid := range mids {

		sp, mv := s.extractionLoop(mid)
		postEx.Amount += sp
		if mv != "" {
			postEx.MIDList = append(postEx.MIDList, mv)
		}
	}
	return postEx
}

func (s *service) extractionLoop(mid selenium.WebElement) (int, string) {
	sp, mv := s.extract(mid)

	if sp {

		success := 1
		return success, mv
	}
	return 0, ""
}

func (s *service) extract(mid selenium.WebElement) (bool, string) {
	mValue, err := mid.GetAttribute("innerText")
	if err != nil {
		s.selenium.HandleSeleniumError(true, err)
	}
	sp := strings.Contains(mValue, "I0")
	return sp, mValue
}

func (s *service) compareCurrentAndPastPostex(currentPostex, pastPostex postExInfo) bool {
	for _, mid := range currentPostex.MIDList {
		if !contains(pastPostex.MIDList, mid) {
			return false
		}
	}
	return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (s *postExInfo) resetPrevious() {
	s.Amount = 0
	s.MIDList = s.MIDList[:0]
}

func (s *postExInfo) setPrevious(st *postExInfo) {
	s.Amount = st.Amount
	s.MIDList = st.MIDList
}
