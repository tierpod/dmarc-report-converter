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
