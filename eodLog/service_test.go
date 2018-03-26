package eodLog

import (
	"context"
	"testing"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/callout"
	"github.com/golang/mock/gomock"
	"github.com/weAutomateEverything/gppMonitor/sftp"
	"github.com/kyokomi/emoji"
	"os"
	"fmt"
	"github.com/weAutomateEverything/gppMonitor/transactionCountLog"
	"time"
)

func TestDateConvert(t *testing.T) {
	dateResult, timeResult := dateTimeConvert("06032018", "11:49:30")

	if dateResult != "06/03/2018" {
		t.Errorf("Date was not propperly formatted, got: %s, want: %s.", dateResult, "06/03/2018")
	}
	if timeResult != "11:49AM" {
		t.Errorf("Time was not propperly formatted, got: %s, want: %s.", timeResult, "11:49AM")
	}
}
func TestResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	
	mockAlert := alert.NewMockService(ctrl)
	mockCallout := callout.NewMockService(ctrl)
	mockSFTP := sftp.NewMockService(ctrl)
	mockTrans := transactionCountLog.NewMockService(ctrl)

	os.Setenv("EDO_LOCATION", "testData/")
	
	svc := NewService(mockCallout, mockSFTP, mockAlert)

	mockCallout.EXPECT().InvokeCallout(context.TODO(), "EDO Posting request file send failed", fmt.Sprintf("EDO Posting request file '%s' send failed on the: %s at %s", "EDO_POSTING_REQ_I03F05535IGU3V1C_.txt", time.Now().Format("02/01/2006"), "1:00PM"))
	mockAlert.EXPECT().SendHeartbeatGroupAlert(context.TODO(), emoji.Sprintf(":rotating_light: EDO Posting request file '%s' send failed on the: %s at %s", "EDO_POSTING_REQ_I03F05535IGU3V1C_.txt", "15/03/2018", "1:00PM"))
	mockSFTP.EXPECT().RetrieveFile("testData", "EDO.log")
	mockTrans.EXPECT().RetrieveNightFileTransactions("EDO.log")
	
	svc.response("sending file EDO_POSTING_REQ_I03F05535IGU3V1C_.txt to EDO sending file EDO_POSTING_REQ_I03F05535IGU3V1C_.txt to EDO failed", "EDO_POSTING_REQ_I03F05535IGU3V1C_.txt", time.Now().Format("02/01/2006"), "1:00PM")
}
