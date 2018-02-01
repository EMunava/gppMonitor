package gppSelenium


import (
"github.com/zamedic/go2hal/halSelenium"
"github.com/zamedic/go2hal/alert"
"github.com/tebeka/selenium"
	"os"
)

type Service interface{
	WaitForWaitFor()
	LogIn()
	LogOut()

	// Override
	HandleSeleniumError(internal bool, err error)
	Driver() selenium.WebDriver
	ClickByClassName(cn string)
	ClickByXPath(xp string)
	ClickByCSSSelector(cs string)

	WaitFor(findBy, selector string)
}

type service struct{
	halSelenium halSelenium.Service
}

func NewService(alert alert.Service)Service{
	sel := halSelenium.NewChromeService(alert,seleniumServer())
	err := sel.Driver().Get(endpoint())
	if err != nil {
		panic(err)
	}
	return &service{sel}
}

func(s *service)HandleSeleniumError(internal bool, err error){
	s.halSelenium.HandleSeleniumError(internal,err)
}

func (s * service)Driver() selenium.WebDriver{
	return s.halSelenium.Driver()
}
func (s *service)WaitForWaitFor() {
	err := s.halSelenium.Driver().Wait(func(wb selenium.WebDriver) (bool, error) {
		elem, err := wb.FindElement(selenium.ByCSSSelector, "body > div.dh-notification.ng-scope.success > div")
		if err != nil {
			return true, nil
		}
		r, err := elem.IsDisplayed()
		return !r, nil
	})
	if err != nil {
		panic(err)
	}
}

func (s *service)LogIn() {

	s.WaitFor( selenium.ByClassName, "dh-input-field")

	user, err := s.Driver().FindElement(selenium.ByName, "txtUserId")
	if err != nil {
		panic(err)
	}

	pass, err := s.Driver().FindElement(selenium.ByName, "txtPassword")
	if err != nil {
		panic(err)
	}

	loginButton, err := s.Driver().FindElement(selenium.ByXPATH, "//*[contains(text(), 'Sign In')]")

	user.SendKeys(gppUser())
	pass.SendKeys(gppPass())
	loginButton.Submit()

	//Wait for successful login
	s.WaitFor(selenium.ByClassName, "dh-customer-logo")


	s.WaitForWaitFor()
}

func (s *service)LogOut() {

	signOutButton, err := s.halSelenium.Driver().FindElement(selenium.ByXPATH, "//*[contains(text(), 'Sign Out')]")
	if err != nil {
		s.HandleSeleniumError(true,err)
		return
	}
	err = signOutButton.Click()
	if err != nil {
		s.HandleSeleniumError(true,err)
		return
	}

	s.WaitFor(selenium.ByClassName, "dh-input-field")
}

func (s *service)ClickByClassName(cn string){
	s.halSelenium.ClickByClassName(cn)
}
func (s *service)ClickByXPath(xp string){
	s.halSelenium.ClickByXPath(xp)
}
func (s *service)ClickByCSSSelector(cs string){
	s.halSelenium.ClickByCSSSelector(cs)
}

func (s *service)WaitFor(findBy, selector string){
	s.WaitFor(findBy,selector)
}

func gppUser() string {
	return os.Getenv("GPP_USER")
}

func gppPass() string {
	return os.Getenv("GPP_PASSWORD")
}

func seleniumServer() string {
	return os.Getenv("SELENIUM_SERVER")
}

func endpoint() string {
	return os.Getenv("GPP_ENDPOINT")
}


