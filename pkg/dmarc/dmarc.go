// Package dmarc contains reader and parser for DMARC xml reports.
package dmarc

import (
	"encoding/xml"
	"fmt"
	"math"
	"net"
	"sort"
	"time"
)

// ReportIDDateTime is the DateTime format for Report.ID
const ReportIDDateTime = "2006-01-02"

// Report represents root of dmarc report struct
type Report struct {
	XMLName            xml.Name        `xml:"feedback" json:"feedback"`
	ReportMetadata     ReportMetadata  `xml:"report_metadata" json:"report_metadata"`
	PolicyPublished    PolicyPublished `xml:"policy_published" json:"policy_published"`
	Records            []Record        `xml:"record" json:"records"`
	Total              int             `json:"_total"`
	TotalFailed        int             `json:"_total_failed"`
	TotalPassed        int             `json:"_total_passed"`
	TotalPassedPercent float64         `json:"_total_passed_percent"`
}

// TotalMessages calculates total amount of messages
func (r *Report) TotalMessages() int {
	total := 0
	for _, record := range r.Records {
		total = total + record.Row.Count
	}

	return total
}

// ID returns report identifier in format YEAR-MONTH-DAY-DOMAIN/EMAIL-ID (can be used in config to
// calculate filename)
func (r Report) ID() string {
	d := r.ReportMetadata.DateRange.Begin.Format(ReportIDDateTime)
	return fmt.Sprintf("%v-%v/%v-%v", d, r.PolicyPublished.Domain, r.ReportMetadata.Email, r.ReportMetadata.ReportID)
}

// ReportMetadata represents feedback>report_metadata section
type ReportMetadata struct {
	OrgName          string    `xml:"org_name" json:"org_name"`
	Email            string    `xml:"email" json:"email"`
	ExtraContactInfo string    `xml:"extra_contact_info" json:"extra_contact_info"`
	ReportID         string    `xml:"report_id" json:"report_id"`
	DateRange        DateRange `xml:"date_range" json:"date_range"`
}

// DateRange represents feedback>report_metadata>date_range section
type DateRange struct {
	Begin DateTime `xml:"begin" json:"begin"`
	End   DateTime `xml:"end" json:"end"`
}

// DateTime is the custom time for DateRange.Begin and DateRange.End values
type DateTime struct {
	time.Time
}

// UnmarshalXML unmarshals unix timestamp to time.Time
func (t *DateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v int64
	d.DecodeElement(&v, &start)
	datetime := time.Unix(v, 0)
	*t = DateTime{datetime}
	return nil
}

// PolicyPublished represents feedback>policy_published section
type PolicyPublished struct {
	Domain  string `xml:"domain" json:"domain"`
	ADKIM   string `xml:"adkim" json:"adkim"`
	ASPF    string `xml:"aspf" json:"aspf"`
	Policy  string `xml:"p" json:"p"`
	SPolicy string `xml:"sp" json:"sp"`
	Pct     string `xml:"pct" json:"pct"`
}

// Record represents feedback>record section
type Record struct {
	Row         Row         `xml:"row" json:"row"`
	Identifiers Identifiers `xml:"identifiers" json:"identifiers"`
	AuthResults AuthResults `xml:"auth_results" json:"auth_results"`
	IsPassed    bool        `json:"_is_passed"`
}

// IsPass returns true if DKIM or SPF policies are passed
func (r *Record) IsPass() bool {
	return (r.Row.PolicyEvaluated.DKIM == "pass" || r.Row.PolicyEvaluated.SPF == "pass")
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
	HeaderFrom   string `xml:"header_from" json:"header_from"`
	EnvelopeFrom string `xml:"envelope_from" json:"envelope_from"`
}

// AuthResults represents feedback>record>auth_results section
type AuthResults struct {
	DKIM DKIMAuthResult `xml:"dkim" json:"dkim"`
	SPF  SPFAuthResult  `xml:"spf" json:"spf"`
}

// DKIMAuthResult represnets feedback>record>auth_results>dkim sections
type DKIMAuthResult struct {
	Domain   string `xml:"domain" json:"domain"`
	Result   string `xml:"result" json:"result"`
	Selector string `xml:"selector" json:"selector"`
}

// SPFAuthResult represnets feedback>record>auth_results>spf section
type SPFAuthResult struct {
	Domain string `xml:"domain" json:"domain"`
	Result string `xml:"result" json:"result"`
	Scope  string `xml:"scope" json:"scope"`
}

// Parse parses input xml data b to Report struct. If lookupAddr is true, performs a reverse
// lookups for feedback>record>row>source_ip
func Parse(b []byte, lookupAddr bool) (Report, error) {
	var result Report
	err := xml.Unmarshal(b, &result)
	if err != nil {
		return Report{}, err
	}

	sort.Slice(result.Records, func(i, j int) bool {
		return result.Records[i].Row.Count > result.Records[j].Row.Count
	})

	// count all counters
	result.Total = result.TotalMessages()
	for i, record := range result.Records {
		result.Records[i].IsPassed = record.IsPass()
		if result.Records[i].IsPassed {
			result.TotalPassed = result.TotalPassed + record.Row.Count
		}
	}
	result.TotalFailed = result.Total - result.TotalPassed
	result.TotalPassedPercent = math.Round((float64(result.TotalPassed) / float64(result.Total)) * 100)

	if lookupAddr {
		doPTRlookup(&result)
	}

	return result, nil
}

func doPTRlookup(r *Report) error {
	for i, record := range r.Records {
		var hostname string
		hostnames, err := net.LookupAddr(record.Row.SourceIP)
		if err != nil {
			hostname = ""
		} else {
			hostname = hostnames[0]
		}
		r.Records[i].Row.SourceHostname = hostname
	}

	return nil
}
