package eodLog

import "testing"

func TestDateConvert(t *testing.T) {
	dateResult, timeResult := dateTimeConvert("06032018", "11:49:30")

	if dateResult != "06/03/2018" {
		t.Errorf("Date was not propperly formatted, got: %s, want: %s.", dateResult, "06/03/2018")
	}
	if timeResult != "11:49AM" {
		t.Errorf("Time was not propperly formatted, got: %s, want: %s.", timeResult, "11:49AM")
	}
}
