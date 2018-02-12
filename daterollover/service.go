package daterollover

import (
	"fmt"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/gppSelenium"
	"github.com/matryer/try"
	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
	"github.com/zamedic/go2hal/alert"
	"gopkg.in/kyokomi/emoji.v1"
	"log"
	"strings"
	"time"
)

type Service interface {
	ConfirmDateRollOver()
}

type service struct {
	selenium     gppSelenium.Service
	alertService alert.Service
}

func NewService(alert alert.Service, selenium gppSelenium.Service) Service {
	return &service{alertService: alert, selenium: selenium}
}

func (s *service) ConfirmDateRollOverMethod() (r error) {
	s.selenium.NewClient()

	defer s.selenium.Driver().Quit()

	defer func() {
		if err := recover(); err != nil {
			s.selenium.HandleSeleniumError(true, errors.New(fmt.Sprint(err)))
			s.selenium.LogOut()
			r = errors.New("Date confirmation failed")
		}
	}()

	s.selenium.LogIn()

	s.navigateToDates()

	Success := s.extractDates()

	_, _, cd, td := date()

	switch Success {

	case 2:
		s.selenium.HandleSeleniumError(false, errors.New(emoji.Sprintf(":white_check_mark: Global and ZA dates have successfully rolled over to: %s", cd)))
	case 1:
		s.selenium.HandleSeleniumError(false, errors.New(emoji.Sprintf(":white_check_mark: Global date has successfully roled over to: %s", td)))
	case 0:
		s.selenium.HandleSeleniumError(false, errors.New(emoji.Sprintf(":rotating_light: Global and ZA dates have failed to roll over to : %s", cd)))
	}
	s.selenium.LogOut()

	return nil
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

	currentDate, _, cd, td := date()

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

func (s *service) ConfirmDateRollOver() {
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		err = s.ConfirmDateRollOverMethod()
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

func date() (time.Time, time.Time, string, string) {

	currentDate := time.Now()
	tomorrowDate := currentDate.AddDate(0, 0, 1)

	cd := currentDate.Format("02/01/2006")
	td := tomorrowDate.Format("02/01/2006")

	return currentDate, tomorrowDate, cd, td
}