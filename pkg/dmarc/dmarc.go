// Package dmarc contains xml parser
package dmarc

import (
	"encoding/xml"
	"net"
	"sort"
)

// Report represents root of dmarc report struct
type Report struct {
	XMLName         xml.Name        `xml:"feedback" json:"feedback"`
	ReportMetadata  ReportMetadata  `xml:"report_metadata" json:"report_metadata"`
	PolicyPublished PolicyPublished `xml:"policy_published" json:"policy_published"`
	Record          []Record        `xml:"record" json:"record"`
}

// ReportMetadata represents feedback>report_metadata section
type ReportMetadata struct {
	OrgName   string    `xml:"org_name" json:"org_name"`
	Email     string    `xml:"email" json:"email"`
	ReportID  string    `xml:"report_id" json:"report_id"`
	DateRange DateRange `xml:"date_range" json:"date_range"`
}

// DateRange represents feedback>report_metadata>date_range section
type DateRange struct {
	Begin int `xml:"begin" json:"begin"` // TODO: time
	End   int `xml:"end" json:"end"`     // TODO: time
}

// PolicyPublished represents feedback>policy_published section
type PolicyPublished struct {
	Domain string `xml:"domain" json:"domain"`
	ADKIM  string `xml:"adkim" json:"adkim"`
	ASPF   string `xml:"aspf" json:"aspf"`
	Policy string `xml:"p" json:"p"`
	Pct    string `xml:"pct" json:"pct"`
}

// Record represents feedback>record section
type Record struct {
	Row         Row         `xml:"row" json:"row"`
	Identifiers Identifiers `xml:"identifiers" json:"identifiers"`
	AuthResults AuthResults `xml:"auth_results" json:"auth_results"`
}

// Row represents feedback>record>row section
type Row struct {
	SourceIP        string          `xml:"source_ip" json:"source_ip"`
	Count           int             `xml:"count" json:"count"`
	PolicyEvaluated PolicyEvaluated `xml:"policy_evaluated" json:"policy_evaluated"`
	SourceHostname  string          `json:"source_hostname"`
}

// PolicyEvaluated represents feedback>record>row>policy_evaluated section
type PolicyEvaluated struct {
	Disposition string `xml:"disposition" json:"disposition"`
	DKIM        string `xml:"dkim" json:"dkim"`
	SPF         string `xml:"spf" json:"spf"`
}

// Identifiers represents feedback>record>identifiers section
type Identifiers struct {
	HeaderFrom string `xml:"header_from" json:"header_from"`
}

// AuthResults represents feedback>record>auth_results section
type AuthResults struct {
	DKIM AuthResult `xml:"dkim" json:"dkim"`
	SPF  AuthResult `xml:"spf" json:"spf"`
}

// AuthResult represnets feedback>record>auth_results>dkim and feedback>record>auth_results>spf sections
type AuthResult struct {
	Domain string `xml:"domain" json:"domain"`
	Result string `xml:"result" json:"result"`
}

// Parse parses input xml bytes to DMARCReport struct. If lookupAddr is true, performs a reverse
// lookups for feedback>record>row>source_ip
func Parse(b []byte, lookupAddr bool) (Report, error) {
	var result Report
	err := xml.Unmarshal(b, &result)
	if err != nil {
		return Report{}, err
	}

	sort.Slice(result.Record, func(i, j int) bool {
		return result.Record[i].Row.Count > result.Record[j].Row.Count
	})

	if lookupAddr {
		reportLookupAddr(&result)
	}

	return result, nil
}

func reportLookupAddr(r *Report) error {
	for i, record := range r.Record {
		var hostname string
		hostnames, err := net.LookupAddr(record.Row.SourceIP)
		if err != nil {
			hostname = ""
		} else {
			hostname = hostnames[0]
		}
		r.Record[i].Row.SourceHostname = hostname
	}

	return nil
}
