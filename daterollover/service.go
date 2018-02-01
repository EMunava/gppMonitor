package daterollover

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
	ConfirmDateRollOver()
}

type service struct {
	selenium     gppSelenium.Service
	alertService alert.Service
}

func NewService(alert alert.Service) Service {
	return &service{alertService:alert}
}

func (s *service) ConfirmDateRollOver() {
	s.selenium = gppSelenium.NewService(s.alertService)
	defer s.selenium.Driver().Close()

	defer func() {
		if err := recover(); err != nil {
			s.selenium.HandleSeleniumError(true, errors.New(fmt.Sprint(err)))
			s.selenium.LogOut()
		}
	}()

	s.selenium.LogIn()

	s.navigateToDates()

	Success := s.extractDates()

	switch Success {

	case 2:
		s.selenium.HandleSeleniumError(false, errors.New("Global and ZA date rollovers have completed successfully"))
	case 1:
		s.selenium.HandleSeleniumError(false, errors.New("Global date rollover has completed successfully"))
	case 0:
		s.selenium.HandleSeleniumError(false, errors.New("Global and ZA dates have failed to rollover to the next business day"))
	}
	s.selenium.LogOut()
}

func (s *service) navigateToDates() {

	s.selenium.ClickByXPath("//*[contains(text(), 'Business Setup')]")

	s.selenium.WaitFor(selenium.ByClassName, "ft-grid-click")
}

func (s *service) extractDates() int {

	Success := 0

	dates, err := s.selenium.Driver().FindElements(selenium.ByClassName, "ui-grid-cell-contents")
	if err != nil {
		panic(err)
	}

	for _, date := range dates {

		Success += s.extractionLoop(date)
	}
	return Success
}

func (s *service) extract(date selenium.WebElement) ([]string, string) {
	dValue, err := date.GetAttribute("innerText")
	if err != nil {
		s.selenium.HandleSeleniumError(true, err)
	}
	sp := strings.Split(dValue, "/")
	return sp, dValue
}

func (s *service) extractionLoop(date selenium.WebElement) int {

	sp, dValue := s.extract(date)

	if len(sp) != 1 {
		success := dateConfirm(dValue)
		return success
	}
	return 0
}

func dateConfirm(d1 string) int {

	currentDate := time.Now()
	tomorrowDate := currentDate.AddDate(0, 0, 1)

	cd := currentDate.Format("02/01/2006")
	td := tomorrowDate.Format("02/01/2006")

	if currentDate.Before(time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 24, 0, 0, 0, currentDate.Location())) && currentDate.After(time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 23, 0, 0, 0, currentDate.Location())) {
		t := strings.Compare(d1, td)
		if t == 0 {
			return 1
		}
		return 0
	}
	c := strings.Compare(d1, cd)
	if c == 0 {
		return 1
	}
	return 0
}
