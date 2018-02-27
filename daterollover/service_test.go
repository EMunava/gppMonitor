package daterollover

import (
	"fmt"
	"github.com/CardFrontendDevopsTeam/GPPMonitor/gppSelenium/mock_gppSekenium"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
	"github.com/zamedic/go2hal/alert"
	gomock2 "github.com/zamedic/go2hal/gomock"
	"github.com/zamedic/go2hal/halSelenium/mock_selenium"
	"github.com/zamedic/go2hal/remoteTelegramCommands"
	"github.com/zamedic/go2hal/remoteTelegramCommands/mock_remoteTelegramCommands"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestService_ConfirmDateRollOver(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockAlert := alert.NewMockService(ctrl)
	mockGppSelenium := mock_gppSekenium.NewMockService(ctrl)
	mockRemoteTelegramCommandClient := mock_remoteTelegramCommands.NewMockRemoteCommandClient(ctrl)
	mockDriver := mock_selenium.NewMockWebDriver(ctrl)
	mockRemoteClient := mock_remoteTelegramCommands.NewMockRemoteCommand_RegisterCommandClient(ctrl)

	mockRemoteTelegramCommandClient.EXPECT().RegisterCommand(context.Background(), &remoteTelegramCommands.RemoteCommandRequest{Name: "GPPDateRolloverCheck", Description: "Execute GPP Date Roll Over"}).Return(mockRemoteClient, nil)
	mockRemoteClient.EXPECT().Recv().Return(nil, errors.New("Out if scope for this test"))

	svc := NewService(mockAlert, mockGppSelenium, mockRemoteTelegramCommandClient)

	mockGppSelenium.EXPECT().NewClient()
	mockGppSelenium.EXPECT().Driver().Times(2).Return(mockDriver)
	mockGppSelenium.EXPECT().LogIn()
	mockDriver.EXPECT().Quit()

	mockGppSelenium.EXPECT().ClickByXPath("//*[contains(text(), 'Business Setup')]")
	mockGppSelenium.EXPECT().WaitFor(selenium.ByClassName, "ft-grid-click")

	mockDriver.EXPECT().FindElements(selenium.ByClassName, "ui-grid-cell-contents")

	currentDate := time.Now()
	cd := currentDate.Format("02/01/2006")
	mockGppSelenium.EXPECT().HandleSeleniumError(false, gomock2.ErrorMsgMatches(fmt.Errorf("🚨  Global and ZA dates have failed to roll over to : %v", cd)))

	mockGppSelenium.EXPECT().LogOut()

	svc.ConfirmDateRollOver()

}
