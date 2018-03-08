package eodLog

import "testing"

func TestDateConvert(t *testing.T) {
	dateResult, timeResult := dateTimeConvert("2018-03-05 14:08:49.723242043 +0200 SAST m=+0.000201828", "11:49:30")

	if dateResult != "01/01/0001" && timeResult == "11:49AM" {
		t.Errorf("Date was not propperly formatted, got: %s, want: %s.", dateResult, "05/03/2018")
	}
}
