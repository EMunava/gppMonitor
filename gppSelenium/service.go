package gppSelenium


import (
"github.com/zamedic/go2hal/halSelenium"
"github.com/zamedic/go2hal/alert"
"github.com/tebeka/selenium"
)

type Service interface{
	HandleSeleniumError(err error)
	Driver() selenium.WebDriver
	WaitForWaitFor() error
}

type service struct{
	halSelenium halSelenium.Service
}

func NewService(service alert.Service, seleniumEndpoint string)Service{
	sel := halSelenium.NewChromeService(service,seleniumEndpoint)
	return &service{sel}
}

func(s *service)HandleSeleniumError(err error){
	s.halSelenium.HandleSeleniumError(err)
}

func (s * service)Driver() selenium.WebDriver{
	return s.halSelenium.Driver()
}
\func (s *service)WaitForWaitFor() error {
	return s.halSelenium.Driver().Wait(func(wb selenium.WebDriver) (bool, error) {
		elem, err := wb.FindElement(selenium.ByCSSSelector, "body > div.dh-notification.ng-scope.success > div")
		if err != nil {
			return true, nil
		}
		r, err := elem.IsDisplayed()
		return !r, nil
	})
}



