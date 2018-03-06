package extractFooterTransactions

import (
	"testing"
)

func TestLastLines(t *testing.T) {
	result := lastLines("testData/","2018-02-22_00-01-36.366_ACDEBIT.RESPONSE.LEG.SAP.20180222.000134")

	if result != "9900108121234320000000000000269600000000000000005"{
		t.Errorf("Transaction amount was not successfully extracted. Expected: %v, Extracted: %v", "9900108121234320000000000000269600000000000000005", result)
	}
}

func TestExtractTransactionAmount(t *testing.T) {
	result := extractTransactionAmount("9900108121234320000000000000269600000000000000005")

	if result != 5{
		t.Errorf("Transaction string was not successfully parsed to integer amount. Expected: %v, Extracted: %v", "5", result)
	}
}
