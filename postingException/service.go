package postingException

import (
	"fmt"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/gppSelenium"
	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
	"github.com/zamedic/go2hal/alert"
	"regexp"
	"strconv"
	"time"
)

//Service interface exposes the ConfirmPostingException method
type Service interface {
	ConfirmPostingException()
}

type service struct {
	selenium gppSelenium.Service
	alert    alert.Service
}

//NewService function creates instances of required external service structs for local use
func NewService(alert alert.Service, selenium gppSelenium.Service) Service {
	return &service{alert: alert, selenium: selenium}
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

	s.selenium.LogIn()

	s.navigateToPostingExce()

	s.selenium.WaitFor(selenium.ByCSSSelector, "#main-content > div.dh-main-container.ng-scope > div > div > div.dh-main-right-container.ng-scope > div > div > div > div > div > div.ft-top-grid-action > div.pull-left > div.top-grid-action-section-title > span")

	px, err := s.selenium.Driver().FindElements(selenium.ByXPATH, "//*[contains(text(), 'Posting Exception')]")
	if err != nil {
		panic(err)
	}

	postEx := s.extractInteger(s.extractString(px[0]))

	s.selenium.HandleSeleniumError(false, fmt.Errorf("Posting Exception count: %d for %v", postEx, time.Now().Format("02/01/2006")))

	s.selenium.LogOut()
}

func (s *service) navigateToPostingExce() {

	s.selenium.ClickByClassName("dh-navigation-tabs-current-tab-button")

	s.selenium.ClickByCSSSelector("#main-content > div.dh-main-container.ng-scope > div > div > div:nth-child(2) > div.dh-navigation-tabs > div.dialer-container > ul > li:nth-child(1) > button")

	s.selenium.WaitFor(selenium.ByXPATH, "//*[contains(text(), 'Individual Messages (')]")

	s.selenium.ClickByXPath("//*[contains(text(), 'Individual Messages (')]")

	s.selenium.ClickByXPath("//*[contains(text(), 'Exception')]")

	s.selenium.ClickByXPath("//*[contains(text(), 'Posting Exception')]")
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
