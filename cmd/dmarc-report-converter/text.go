package main

import (
	"html/template"
	"os"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

const textTemplate = `
DMARC report from {{.ReportMetadata.OrgName}} ({{.ReportMetadata.Email}})
  with id {{.ReportMetadata.ReportID}}
  period {{.ReportMetadata.DateRange.Begin}}-{{.ReportMetadata.DateRange.End}}
-------------------------------------------------------------------------------

Policy published for {{.PolicyPublished.Domain}}:
  p={{.PolicyPublished.Policy}} pct={{.PolicyPublished.Pct}} adkim={{.PolicyPublished.ADKIM}} aspf={{.PolicyPublished.ASPF}}
-------------------------------------------------------------------------------
{{ range .Record }}
policy evaluated for source ip: {{.Row.SourceIP}} (count: {{.Row.Count}}):
  disposition={{.Row.PolicyEvaluated.Disposition}}
  dkim={{.Row.PolicyEvaluated.DKIM}}
  spf={{.Row.PolicyEvaluated.SPF}}

identifiers:
  header_from={{.Identifiers.HeaderFrom}}

auth_result:
  DKIM={{.AuthResults.DKIM.Result}}
  SPF={{.AuthResults.SPF.Result}}
---{{ end }}
`

func textOutput(d dmarc.Report, o string) error {
	t := template.Must(template.New("report").Parse(textTemplate))
	err := t.Execute(os.Stdout, d)
	if err != nil {
		return err
	}

	return nil
}
