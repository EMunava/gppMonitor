package postingException

import (
	"fmt"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/gppSelenium"
	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
	"github.com/zamedic/go2hal/alert"
)

/*
	Service interface exposes the ConfirmPostingException method
 */
type Service interface {
	ConfirmPostingException()
}

type service struct {
	selenium gppSelenium.Service
	alert    alert.Service
}

func NewService(alert alert.Service, selenium gppSelenium.Service) Service {
	return &service{alert: alert, selenium: selenium}
}

func (s *service) ConfirmPostingException() {
	s.selenium.NewClient()

	driver := s.selenium.Driver()

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

	shot, err := driver.Screenshot()
	if err != nil {
		panic("Failed to take screenshot")
	}
	s.alert.SendImageToHeartbeatGroup(shot)

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
