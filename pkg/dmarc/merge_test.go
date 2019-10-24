package dmarc

import (
	"os"
	"reflect"
	"testing"
)

func TestReport_MergeRecord(t *testing.T) {
	// we have tested this errors already
	f, _ := os.Open("testdata/test.xml")
	defer f.Close()
	report, _ := ReadParseXML(f, false)

	// this record must be merged with the xmlRecord2
	r1 := xmlRecord2
	report.MergeRecord(r1)

	// check length after merge
	if len(report.Records) != 2 {
		t.Errorf("Report_MergeRecord: got len %v, want %v", len(report.Records), 2)
	}

	// check count after merge
	for _, r := range report.Records {
		if r.Row.SourceIP == r1.Row.SourceIP {
			if r.Row.Count != 4 {
				t.Errorf("Report_MergeRecord: got count %v, want %v", r.Row.Count, 4)
			}
		}
	}

	// this record must be added
	var r2 = Record{
		Row: Row{
			SourceIP: "172.168.1.1",
			Count:    10,
			PolicyEvaluated: PolicyEvaluated{
				Disposition: "none",
				DKIM:        "fail",
				SPF:         "fail",
			},
		},
		Identifiers: Identifiers{
			HeaderFrom:   "test3.net",
			EnvelopeFrom: "",
		},
		AuthResults: AuthResults{
			DKIM: DKIMAuthResult{
				Domain:   "test3.net",
				Result:   "fail",
				Selector: "selector",
			},
			SPF: SPFAuthResult{
				Domain: "test3.net",
				Result: "softfail",
				Scope:  "mfrom",
			},
		},
	}
	report.MergeRecord(r2)

	// check length after merge
	if len(report.Records) != 3 {
		t.Errorf("Report_MergeRecord: got len %v, want %v", len(report.Records), 3)
	}

	// check last record
	lastRecord := report.Records[len(report.Records)-1]
	if !reflect.DeepEqual(lastRecord, r2) {
		t.Errorf("Report_MergeRecord: got last record %v, want %v", lastRecord, r2)
	}
}
