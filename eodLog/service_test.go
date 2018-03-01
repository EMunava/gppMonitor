package eodLog

import "testing"

func testDateConvert(t *testing.T){
	dateResult:= dateConvert("Mon Jan _2 15:04:05 MST 2006")

	if dateResult != "02/01/2006"{
		t.Errorf("Date was not propperly formatted, got: %d, want: %d.", dateResult, "02/01/2006")
	}
}
