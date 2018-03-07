package extractFooterTransactions

import (
	"github.com/golang/mock/gomock"
	"github.com/weAutomateEverything/go2hal/alert"
	"testing"

	"github.com/weAutomateEverything/gppMonitor/sftp"
	"os"
)

func TestLastLines(t *testing.T) {
	result := lastLines("testData/", "2018-02-22_00-01-36.366_ACDEBIT.RESPONSE.LEG.SAP.20180222.000134")

	if result != "9900108121234320000000000000269600000000000000005" {
		t.Errorf("Transaction amount was not successfully extracted. Expected: %v, Extracted: %v", "9900108121234320000000000000269600000000000000005", result)
	}
}

func TestExtractTransactionAmount(t *testing.T) {
	result := extractTransactionAmount("9900108121234320000000000000269600000000000000005")

	if result != 5 {
		t.Errorf("Transaction string was not successfully parsed to integer amount. Expected: %v, Extracted: %v", "5", result)
	}
}

func TestRetrieveSAPTransactions(t *testing.T) {

	ctrl := gomock.NewController(t)

	mockAlert := alert.NewMockService(ctrl)
	mockSFTP := sftp.NewMockService(ctrl)

	os.Setenv("TRANSACTION_LOCATION", "testData/")

	svc := NewService(mockSFTP, mockAlert)

	mockSFTP.EXPECT().RetrieveFile("testData/", "RESPONSE.SAP")
	mockSFTP.EXPECT().GetFilesInPath("testData/")

	svc.retreiveTransactions("RESPONSE.SAP")

}
