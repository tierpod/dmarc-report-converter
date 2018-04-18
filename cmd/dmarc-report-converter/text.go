package main

import (
	"html/template"
	"os"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

const textTemplate = `
DMARC report with id {{.ReportMetadata.ReportID}}
  from {{.ReportMetadata.OrgName}} ({{.ReportMetadata.Email}})
  period {{.ReportMetadata.DateRange.Begin}}-{{.ReportMetadata.DateRange.End}}

Policy published for {{.PolicyPublished.Domain}}: p={{.PolicyPublished.Policy}} pct={{.PolicyPublished.Pct}} adkim={{.PolicyPublished.ADKIM}} aspf={{.PolicyPublished.ASPF}}
----------------------------------------------------------------------------------------------------
{{printf "%16v|%16v|%4v|%10v|%8v|%8v|%16v|%16v" "ip"          "hostname" "msg"      "disp"                           "dkim r"                 "spf r"            "dkim domain"            "spf domain" }}
----------------------------------------------------------------------------------------------------
{{- range .Record }}
{{printf "%16v|%16v|%4v|%10v|%8v|%8v|%16v|%16v" .Row.SourceIP "TODO"     .Row.Count .Row.PolicyEvaluated.Disposition .AuthResults.DKIM.Result .AuthResults.SPF.Result .AuthResults.DKIM.Domain .AuthResults.DKIM.Domain}}
{{- end }}
`

func textOutput(d dmarc.Report, o string) error {
	t := template.Must(template.New("report").Parse(textTemplate))
	err := t.Execute(os.Stdout, d)
	if err != nil {
		return err
	}

	return nil
}
