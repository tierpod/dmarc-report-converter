package main

import (
	"bytes"
	"html/template"
	"testing"
	"time"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

func TestExternalTemplate(t *testing.T) {
	r := `<?xml version="1.0"?>
<feedback>
  <report_metadata>
    <org_name>Org 1</org_name>
    <email>foo@bar.baz</email>
    <report_id>1712279633.907274</report_id>
    <date_range>
      <begin>1712188800</begin>
      <end>1712275199</end>
    </date_range>
  </report_metadata>
  <policy_published>
    <domain>report.test</domain>
    <adkim>r</adkim>
    <aspf>r</aspf>
    <p>none</p>
    <pct>100</pct>
  </policy_published>
  <record>
    <row>
      <source_ip>1.2.3.4</source_ip>
      <count>1</count>
      <policy_evaluated>
        <disposition>none</disposition>
        <dkim>pass</dkim>
        <spf>fail</spf>
      </policy_evaluated>
    </row>
    <identifiers>
      <header_from>headerfrom.test</header_from>
    </identifiers>
    <auth_results>
      <dkim>
        <domain>auth.test</domain>
        <selector>1000073432</selector>
        <result>pass</result>
      </dkim>
      <dkim>
        <domain>cust.test</domain>
        <selector>2020263919</selector>
        <result>pass</result>
      </dkim>
      <spf>
        <domain>spf.test</domain>
        <result>pass</result>
      </spf>
    </auth_results>
  </record>
</feedback>
`

	report, err := dmarc.Parse([]byte(r), false, 1)
	if err != nil {
		t.Errorf("unexpected error parsing XML: %s", err)
	}

	tmpl := `AssetsPath: {{ .AssetsPath }}
# Report
XMLName: {{ .Report.XMLName.Local }}
ReportMetadata: {{ .Report.ReportMetadata }}
PolicyPublished: {{ .Report.PolicyPublished }}
## Records
{{- range .Report.Records }}
- {{ . }}
{{ end -}}
## MessagesStats
{{ .Report.MessagesStats }}

// Deprecated
XMLName: {{ .XMLName.Local }}
ReportMetadata: {{ .ReportMetadata }}
PolicyPublished: {{ .PolicyPublished }}
{{ .MessagesStats }}
`

	conf := config{
		Output: Output{
			AssetsPath: "/foo",
			template: template.Must(template.New("report").Funcs(
				template.FuncMap{
					"now": func(fmt string) string {
						return time.Now().Format(fmt)
					},
				},
			).Parse(tmpl)),
		},
	}

	var buf bytes.Buffer
	out := newOutput(&conf)
	out.w = &buf

	err = out.template(report)
	if err != nil {
		t.Errorf("unexpected error building template: %s", err)
	}

	expect := `AssetsPath: /foo
# Report
XMLName: feedback
ReportMetadata: {Org 1 foo@bar.baz  1712279633.907274 {2024-04-03 19:00:00 -0500 CDT 2024-04-04 18:59:59 -0500 CDT}}
PolicyPublished: {report.test r r none  100}
## Records
- {{1.2.3.4 1 {none pass fail} } {headerfrom.test } {[{auth.test pass 1000073432} {cust.test pass 2020263919}] [{spf.test pass }]}}
## MessagesStats
{1 0 1 100}

// Deprecated
XMLName: feedback
ReportMetadata: {Org 1 foo@bar.baz  1712279633.907274 {2024-04-03 19:00:00 -0500 CDT 2024-04-04 18:59:59 -0500 CDT}}
PolicyPublished: {report.test r r none  100}
{1 0 1 100}
`

	if buf.String() != expect {
		t.Errorf("Oops!\nWANT:\n%s\nGOT:\n%s", expect, buf.String())
	}
}
