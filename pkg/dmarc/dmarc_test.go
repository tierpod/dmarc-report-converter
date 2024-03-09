package dmarc

import (
	"encoding/xml"
	"os"
	"reflect"
	"testing"
	"time"
)

var xmlReportMetadata = ReportMetadata{
	OrgName:          "Test Inc.",
	Email:            "postmaster@test",
	ExtraContactInfo: "http://test/help",
	ReportID:         "1.id.0",
	DateRange: DateRange{
		Begin: DateTime{time.Unix(1524182400, 0)},
		End:   DateTime{time.Unix(1524268799, 0)},
	},
}

var xmlPolicyPublished = PolicyPublished{
	Domain:  "test.net",
	ADKIM:   "r",
	ASPF:    "r",
	Policy:  "none",
	SPolicy: "",
	Pct:     "100",
}

var xmlRecord1 = Record{
	Row: Row{
		SourceIP: "192.168.1.1",
		Count:    5,
		PolicyEvaluated: PolicyEvaluated{
			Disposition: "none",
			DKIM:        "pass",
			SPF:         "pass",
		},
	},
	Identifiers: Identifiers{
		HeaderFrom:   "test.net",
		EnvelopeFrom: "",
	},
	AuthResults: AuthResults{
		DKIM: DKIMAuthResult{
			Domain:   "test.net",
			Result:   "pass",
			Selector: "selector",
		},
		SPF: SPFAuthResult{
			Domain: "test.net",
			Result: "pass",
			Scope:  "mfrom",
		},
	},
}

var xmlRecord2 = Record{
	Row: Row{
		SourceIP: "10.1.1.1",
		Count:    2,
		PolicyEvaluated: PolicyEvaluated{
			Disposition: "none",
			DKIM:        "fail",
			SPF:         "fail",
		},
	},
	Identifiers: Identifiers{
		HeaderFrom:   "test2.net",
		EnvelopeFrom: "",
	},
	AuthResults: AuthResults{
		DKIM: DKIMAuthResult{
			Domain:   "test2.net",
			Result:   "fail",
			Selector: "selector",
		},
		SPF: SPFAuthResult{
			Domain: "test2.net",
			Result: "softfail",
			Scope:  "mfrom",
		},
	},
}

var xmlMessagesStats = MessagesStats{
	All:           7,
	Failed:        2,
	Passed:        5,
	PassedPercent: 71,
}

var xmlReport = Report{
	XMLName: xml.Name{
		Space: "",
		Local: "feedback",
	},
	ReportMetadata:  xmlReportMetadata,
	PolicyPublished: xmlPolicyPublished,
	Records:         []Record{xmlRecord1, xmlRecord2},
	MessagesStats:   xmlMessagesStats,
}

func TestRecord_IsPassed(t *testing.T) {
	tests := []struct {
		in  PolicyEvaluated
		out bool
	}{
		{
			PolicyEvaluated{DKIM: "pass", SPF: "pass"},
			true,
		},
		{
			PolicyEvaluated{DKIM: "pass", SPF: "fail"},
			true,
		},
		{
			PolicyEvaluated{DKIM: "", SPF: ""},
			false,
		},
	}

	for _, tt := range tests {
		r := Record{Row: Row{PolicyEvaluated: tt.in}}
		got := r.IsPassed()
		if got != tt.out {
			t.Errorf("Record_IsPassed: got %v, want %v", got, tt.out)
		}
	}
}

func TestReadParseXML(t *testing.T) {
	f, err := os.Open("testdata/test.xml")
	if err != nil {
		t.Fatalf("ReadParseXML: %v", err)
	}
	defer f.Close()

	out, err := ReadParseXML(f, false)
	if err != nil {
		t.Fatalf("ReadParseXML: %v", err)
	}

	if !reflect.DeepEqual(out, xmlReport) {
		t.Errorf("ReadParseXML: parsed structs are invalid: %+v", out)
	}
}

func TestReadParseGZIP(t *testing.T) {
	f, err := os.Open("testdata/test.xml.gz")
	if err != nil {
		t.Fatalf("ReadParseGZIP: %v", err)
	}
	defer f.Close()

	out, err := ReadParseGZIP(f, false)
	if err != nil {
		t.Fatalf("ReadParseGZIP: %v", err)
	}

	if !reflect.DeepEqual(out, xmlReport) {
		t.Errorf("ReadParseGZIP: parsed structs are invalid: %+v", out)
	}
}

func TestReadParseZIP(t *testing.T) {
	f, err := os.Open("testdata/test.xml.zip")
	if err != nil {
		t.Fatalf("ReadParseZIP: %v", err)
	}
	defer f.Close()

	out, err := ReadParseZIP(f, false)
	if err != nil {
		t.Fatalf("ReadParseZIP: %v", err)
	}

	if !reflect.DeepEqual(out, xmlReport) {
		t.Errorf("ReadParseZIP: parsed structs are invalid: %+v", out)
	}
}

func TestReadParse(t *testing.T) {
	testFiles := []string{"testdata/test.xml", "testdata/test.xml.gz", "testdata/test.xml.zip", "testdata/test.xml.gz.gz"}
	for _, testFile := range testFiles {
		f, err := os.Open(testFile)
		if err != nil {
			t.Fatalf("ReadParse(%v): %v", testFile, err)
		}
		defer f.Close()

		out, err := ReadParse(f, false)
		if err != nil {
			t.Fatalf("ReadParse(%v): %v", testFile, err)
		}

		if !reflect.DeepEqual(out, xmlReport) {
			t.Errorf("ReadParse(%v): parsed structs are invalid: %+v", testFile, out)
		}
	}
}

// Test empty data
var xmlEmptyReport = Report{
	XMLName: xml.Name{
		Space: "",
		Local: "feedback",
	},
	ReportMetadata:  xmlReportMetadata,
	PolicyPublished: xmlPolicyPublished,
	Records: []Record{
		{
			Row: Row{
				SourceIP: "",
				Count:    0,
				PolicyEvaluated: PolicyEvaluated{
					Disposition: "",
					DKIM:        "",
					SPF:         "",
				},
			},
			Identifiers: Identifiers{
				HeaderFrom:   "",
				EnvelopeFrom: "",
			},
			AuthResults: AuthResults{
				DKIM: DKIMAuthResult{
					Domain:   "",
					Result:   "",
					Selector: "",
				},
				SPF: SPFAuthResult{
					Domain: "",
					Result: "",
					Scope:  "",
				},
			},
		},
	},
	MessagesStats: MessagesStats{
		All:           0,
		Failed:        0,
		Passed:        0,
		PassedPercent: 0,
	},
}

func TestReadParse_Empty(t *testing.T) {
	testFile := "testdata/test_empty.xml"
	f, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("ReadParse(%v): %v", testFile, err)
	}
	defer f.Close()
	out, err := ReadParse(f, false)
	if err != nil {
		t.Fatalf("ReadParse(%v): %v", testFile, err)
	}

	if !reflect.DeepEqual(out, xmlEmptyReport) {
		t.Errorf("ReadParse(%v): parsed structs are invalid: %+v", testFile, out)
	}
}
