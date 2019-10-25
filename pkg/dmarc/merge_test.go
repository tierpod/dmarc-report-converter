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

func TestReport_MergeReport(t *testing.T) {
	f, err := os.Open("testdata/test.xml")
	if err != nil {
		t.Fatalf("Report_MergeReport: %v", err)
	}
	defer f.Close()

	rep1, err := ReadParseXML(f, false)
	if err != nil {
		t.Fatalf("Report_MergeReport: %v", err)
	}

	// test2.xml must be merged into test1.xml
	f, err = os.Open("testdata/test2.xml")
	if err != nil {
		t.Fatalf("Report_MergeReport: %v", err)
	}
	defer f.Close()

	rep2, err := ReadParseXML(f, false)
	if err != nil {
		t.Fatalf("Report_MergeReport: %v", err)
	}

	rep1.MergeReport(rep2)

	// check merged report dates.
	// Begin date must be the earliest date
	d1 := rep1.ReportMetadata.DateRange.Begin
	d2 := rep2.ReportMetadata.DateRange.Begin
	if d1 != d2 {
		t.Errorf("Report_MergeReport: got begin date %v, want %v", d1, d2)
	}

	// End date must be the oldest date
	d1 = rep1.ReportMetadata.DateRange.End
	d2 = rep2.ReportMetadata.DateRange.End
	if d1 != d2 {
		t.Errorf("Report_MergeReport: got end date %v, want %v", d1, d2)
	}

	// chech total amount of records
	if len(rep1.Records) != 3 {
		t.Errorf("Report_MergeReport: got Records length %v, want %v", len(rep1.Records), 3)
	}

	// check records
	var tests = []struct {
		idx       int
		outSource string
		outCount  int
	}{
		{0, "192.168.1.1", 15},
		{1, "10.1.1.1", 2},
		{2, "172.16.1.1", 2},
	}

	for idx, tt := range tests {
		inCount := rep1.Records[idx].Row.Count
		if inCount != tt.outCount {
			t.Errorf("Report_MergeReport: got Count %v, want %v", inCount, tt.outCount)
		}

		inSource := rep1.Records[idx].Row.SourceIP
		if inSource != tt.outSource {
			t.Errorf("Report_MergeReport: got Source %v, want %v", inSource, tt.outCount)
		}
	}

	// check that messages statistic is updated
	rep1.CalculateStats()
	inStats := rep1.MessagesStats
	outStats := MessagesStats{
		All:           19,
		Failed:        2,
		Passed:        17,
		PassedPercent: 89,
	}

	if !reflect.DeepEqual(inStats, outStats) {
		t.Errorf("Report_MergeReport: got MessagesStats %v, want %v", inStats, outStats)
	}
}
