package transactionCountLog

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/callout"
	"github.com/weAutomateEverything/gppMonitor/sftp"
	"golang.org/x/net/context"
	"os"
	"testing"
	"time"
)

func TestLastLines(t *testing.T) {
	result, ts := lastLines("testData/", "2018-02-22_00-01-36.366_ACDEBIT.RESPONSE.LEG.SAP.20180222.000134")

	if result != "9900108121234320000000000000269600000000000000005                                                                                 " && ts.processed != 5 {
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
	mockCallout := callout.NewMockService(ctrl)

	os.Setenv("TRANSACTION_LOCATION", "testData/")

	svc := NewService(mockCallout, mockSFTP, mockAlert)

	mockSFTP.EXPECT().RetrieveFile("testData/", "RESPONSE.SAP")
	mockSFTP.EXPECT().GetFilesInPath("testData/")
	mockCallout.EXPECT().InvokeCallout(context.TODO(), fmt.Sprintf("%v file has not yet arrived from EDO at: %v", "RESPONSE.SAP", time.Now().Format("3:04PM")), fmt.Sprintf("%v file has not yet arrived from EDO at: %v", "RESPONSE.SAP", time.Now().Format("3:04PM")))

	svc.retreiveTransactions("RESPONSE.SAP")

}
