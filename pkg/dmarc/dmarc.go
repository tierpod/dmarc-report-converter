// Package dmarc contains xml parser
package dmarc

import (
	"encoding/xml"
	"sort"
)

// Report represents root of dmarc report struct
type Report struct {
	XMLName         xml.Name        `xml:"feedback"`
	ReportMetadata  ReportMetadata  `xml:"report_metadata"`
	PolicyPublished PolicyPublished `xml:"policy_published"`
	Record          []Record        `xml:"record"`
}

// ReportMetadata represents feedback>report_metadata section
type ReportMetadata struct {
	OrgName   string    `xml:"org_name"`
	Email     string    `xml:"email"`
	ReportID  string    `xml:"report_id"`
	DateRange DateRange `xml:"date_range"`
}

// DateRange represents feedback>report_metadata>date_range section
type DateRange struct {
	Begin int `xml:"begin"` // TODO: time
	End   int `xml:"end"`   // TODO: time
}

// PolicyPublished represents feedback>policy_published section
type PolicyPublished struct {
	Domain string `xml:"domain"`
	ADKIM  string `xml:"adkim"`
	ASPF   string `xml:"aspf"`
	Policy string `xml:"p"`
	Pct    string `xml:"pct"`
}

// Record represents feedback>record section
type Record struct {
	Row         Row         `xml:"row"`
	Identifiers Identifiers `xml:"identifiers"`
	AuthResults AuthResults `xml:"auth_results"`
}

// Row represents feedback>record>row section
type Row struct {
	SourceIP        string          `xml:"source_ip"` // TODO: ip
	Count           int             `xml:"count"`
	PolicyEvaluated PolicyEvaluated `xml:"policy_evaluated"`
}

// PolicyEvaluated represents feedback>record>row>policy_evaluated section
type PolicyEvaluated struct {
	Disposition string `xml:"disposition"`
	DKIM        string `xml:"dkim"`
	SPF         string `xml:"spf"`
}

// Identifiers represents feedback>record>identifiers section
type Identifiers struct {
	HeaderFrom string `xml:"header_from"`
}

// AuthResults represents feedback>record>auth_results section
type AuthResults struct {
	DKIM AuthResult `xml:"dkim"`
	SPF  AuthResult `xml:"spf"`
}

// AuthResult represnets feedback>record>auth_results>dkim and feedback>record>auth_results>spf sections
type AuthResult struct {
	Domain string `xml:"domain"`
	Result string `xml:"result"`
}

// Parse parses input xml bytes to DMARCReport struct
func Parse(b []byte) (Report, error) {
	var result Report
	err := xml.Unmarshal(b, &result)
	if err != nil {
		return Report{}, err
	}

	sort.Slice(result.Record, func(i, j int) bool {
		return result.Record[i].Row.Count > result.Record[j].Row.Count
	})

	return result, nil
}
