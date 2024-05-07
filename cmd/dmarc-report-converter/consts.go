package main

const htmlMain = `
    <p></p>
    <div class="container">
        <div class="row">
            <div class="col">
                <div class="card">
                    <div class="card-header">
                        DMARC Report, id {{.Report.ReportMetadata.ReportID}}
                    </div>
                    <div class="card-body">
                        <table class="table table-sm borderless bottomless">
                            <tbody>
                                <tr>
                                    <td>Organization</td>
                                    <td><a href="{{.Report.ReportMetadata.ExtraContactInfo}}">{{.Report.ReportMetadata.OrgName}}</a> ({{.Report.ReportMetadata.Email}})</td>
                                </tr>
                                <tr>
                                    <td>Date range</td>
                                    <td>since {{.Report.ReportMetadata.DateRange.Begin.UTC}} until {{.Report.ReportMetadata.DateRange.End.UTC}}</td>
                                </tr>
                                <tr>
                                    <td>Policy published</td>
                                    <td>{{.Report.PolicyPublished.Domain}}: p={{.Report.PolicyPublished.Policy}} sp={{.Report.PolicyPublished.SPolicy}} pct={{.Report.PolicyPublished.Pct}} adkim={{.Report.PolicyPublished.ADKIM}} aspf={{.Report.PolicyPublished.ASPF}}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
        <p></p>
        <div class="row">
            <div class="col">
                <canvas id="hosts-chart"></canvas>
            </div>
            <div class="col">
                <canvas id="stats-chart"></canvas>
                <center><span class="badge bg-success">passed {{ .Report.MessagesStats.Passed }}</span> <span class="badge bg-danger">failed {{ .Report.MessagesStats.Failed }}</span> <span class="badge text-bg-light">total {{ .Report.MessagesStats.All }}</span></center>
            </div>
        </div>
        <p></p>
        <div class="row">
            <div class="col">
                <table class="table table-bordered table-sm" id="items-table">
                    <thead>
                        <tr>
                        <th scope="col" colspan="3"></th>
                        <th scope="col" colspan="3" class="left-border">policy evaluated</th>
                        <th scope="col" colspan="4" class="left-border">auth results</th>
                        </tr>
                        <tr>
                        <th scope="col">ip</th>
                        <th scope="col" data-toggle="tooltip" title="ptr records">hostname</th>
                        <th scope="col" data-toggle="tooltip" title="messages count">msgs</th>
                        <th scope="col" class="left-border">disposition</th>
                        <th scope="col">DKIM</th>
                        <th scope="col">SPF</th>
                        <th scope="col" class="left-border">DKIM domain</th>
                        <th scope="col">result</th>
                        <th scope="col">SPF domain</th>
                        <th scope="col">result</th>
                        </tr>
                    </thead>
                    <tbody id="items-table">
                        {{- range .Report.Records }}
                        {{- if .IsPassed }}
                        <tr class="policy-evaluated-result-pass table-success">
                        {{- else }}
                        <tr class="policy-evaluated-result-fail">
                        {{- end }}
                        <td id="ip">{{.Row.SourceIP}}</td>
                        <td id="hostname">{{.Row.SourceHostname}}</td>
                        <td id="msgc">{{.Row.Count}}</td>
                        <td class="left-border" title="identifiers&#13;header_from: {{.Identifiers.HeaderFrom}}&#13;envelope_from: {{.Identifiers.EnvelopeFrom}}">{{.Row.PolicyEvaluated.Disposition}}</td>
                        <td>
                            {{- if eq .Row.PolicyEvaluated.DKIM "fail" }}
                            <span class="badge bg-danger">{{.Row.PolicyEvaluated.DKIM}}</span>
                            {{- else }}
                            <span class="badge bg-success">{{.Row.PolicyEvaluated.DKIM}}</span>
                            {{- end}}
                        </td>
                        <td>
                            {{- if eq .Row.PolicyEvaluated.SPF "fail" }}
                            <span class="badge bg-danger">{{.Row.PolicyEvaluated.SPF}}</span>
                            {{- else }}
                            <span class="badge bg-success">{{.Row.PolicyEvaluated.SPF}}</span>
                            {{- end }}
                        </td>
                        <td class="left-border">
                        {{- $len := len .AuthResults.DKIM }}
                        {{- range $j, $foo := .AuthResults.DKIM }}
                            {{- .Domain }}
                            {{- if lt $j $len }}<br/>{{ end }}
                        {{- end -}}
                        </td>
                        <td>
                        {{- range $j, $d := .AuthResults.DKIM }}
                            <span title="selector: {{.Selector}}" class="badge bg-
                            {{- if eq .Result "pass" }}success
                            {{- else if eq .Result "fail" }}danger
                            {{- else }}warning
                            {{- end -}}
                            ">{{.Result}}</span>
                            {{- if lt $j $len }}<br/>{{ end }}
                        {{- end -}}
                        </td>
                        <td>
                        {{- $len = len .AuthResults.SPF }}
                        {{- range $j, $foo := .AuthResults.SPF }}
                            {{- .Domain }}
                            {{- if lt $j $len }}<br/>{{ end }}
                        {{- end -}}
                        <td>
                        {{- range $j, $d := .AuthResults.SPF }}
                            <span title="scope: {{.Scope}}" class="badge bg-
                            {{- if eq .Result "pass" }}success
                            {{- else if eq .Result "fail" }}danger
                            {{- else }}warning
                            {{- end -}}
                            ">{{.Result}}</span>
                            {{- if lt $j $len }}<br/>{{ end }}
                        {{- end -}}
                        </td>
                        </tr>
                        {{- end }}
                    </tbody>
                </table>
            </div>
        </div>
    </div>
`

const htmlTmpl = `
<!doctype html>
<html lang="en">
<head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap core CSS -->
    <link href="{{.AssetsPath}}/css/bootstrap.min.css" rel="stylesheet">
<style>
td.left-border, th.left-border
{
    border-left: 1px solid #dee2e6;
}

.borderless td, .borderless th
{
    border: none;
}

table.table.bottomless
{
    margin-bottom: 0rem;
}
</style>
</head>

<body>
    <script src="{{.AssetsPath}}/js/jquery-3.3.1.min.js"></script>
    <script src="{{.AssetsPath}}/js/bootstrap.min.js"></script>
    <script src="{{.AssetsPath}}/js/chart.umd.min.js"></script>
    <script src="{{.AssetsPath}}/js/dmarc-report-converter.js"></script>
` + htmlMain + `
</body>
</html>`

const htmlStaticTmpl = `
<!doctype html>
<html lang="en">
<head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap core CSS -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-gH2yIJqKdNHPEq0n4Mqa/HGKIhSkIHeL5AyhkYV8i59U5AR6csBvApHHNl/vI1Bx" crossorigin="anonymous">
<style>
td.left-border, th.left-border
{
    border-left: 1px solid #dee2e6;
}

.borderless td, .borderless th
{
    border: none;
}

table.table.bottomless
{
    margin-bottom: 0rem;
}
</style>
</head>

<body>` + htmlMain + `
</body>
</html>`

// NOTE(moorereason): This template assumes only one SPF result will be present even though the DMARC spec allows for multiple.
const txtTmpl = `
DMARC report with id {{.ReportMetadata.ReportID}}
  Organization:     {{.ReportMetadata.ExtraContactInfo}} ({{.ReportMetadata.Email}}))
  Date range:       since {{.ReportMetadata.DateRange.Begin.UTC}} until {{.ReportMetadata.DateRange.End.UTC}}
  Policy published: {{.PolicyPublished.Domain}}: p={{.PolicyPublished.Policy}} sp={{.PolicyPublished.SPolicy}} pct={{.PolicyPublished.Pct}} adkim={{.PolicyPublished.ADKIM}} aspf={{.PolicyPublished.ASPF}}
--------------------------------------------------------------------------------------------------------------------------
{{printf "%24v | %32v | %58v |" "" "policy evaluated" "auth results" }}
--------------------------------------------------------------------------------------------------------------------------
{{printf "%17v | %4v | %10v | %8v | %8v | %16v | %9v | %16v | %8v | %v" "ip"          "msgs"     "disp"                           "dkim"                    "spf"                    "dkim domain"            "dkim res"               "spf domain"            "spf res"               "hostname"         }}
--------------------------------------------------------------------------------------------------------------------------
{{- range .Records }}
	{{- $prefix := "  " }}{{ if .IsPassed }}{{ $prefix = "* " }}{{ end }}
	{{- $dkimLen := len .AuthResults.DKIM }}
	{{- if eq $dkimLen 0 }}
{{ printf "%2v%15v | %4v | %10v | %8v | %8v | %16v | %9v | %16v | %8v | %v" $prefix .Row.SourceIP .Row.Count .Row.PolicyEvaluated.Disposition .Row.PolicyEvaluated.DKIM .Row.PolicyEvaluated.SPF "" "" (index .AuthResults.SPF 0).Domain (index .AuthResults.SPF 0).Result .Row.SourceHostname }}
	{{- else if eq $dkimLen 1 }}
		{{- $dkimDomain := (index .AuthResults.DKIM 0).Domain }}
		{{- $dkimResult := (index .AuthResults.DKIM 0).Result }}
{{ printf "%2v%15v | %4v | %10v | %8v | %8v | %16v | %9v | %16v | %8v | %v" $prefix .Row.SourceIP .Row.Count .Row.PolicyEvaluated.Disposition .Row.PolicyEvaluated.DKIM .Row.PolicyEvaluated.SPF $dkimDomain $dkimResult (index .AuthResults.SPF 0).Domain (index .AuthResults.SPF 0).Result .Row.SourceHostname }}
	{{- else }}
		{{- $rec := . }}
		{{- range $i, $res := .AuthResults.DKIM }}
			{{- $dkimDomain := $res.Domain }}
			{{- $dkimResult := $res.Result }}
			{{- if eq $i 0 }}
{{ printf "%2v%15v | %4v | %10v | %8v | %8v | %16v | %9v | %16v | %8v | %v" $prefix $rec.Row.SourceIP $rec.Row.Count $rec.Row.PolicyEvaluated.Disposition $rec.Row.PolicyEvaluated.DKIM $rec.Row.PolicyEvaluated.SPF $dkimDomain $dkimResult (index $rec.AuthResults.SPF 0).Domain (index $rec.AuthResults.SPF 0).Result $rec.Row.SourceHostname }}
			{{- else }}
{{ printf "%2v%15v | %4v | %10v | %8v | %8v | %16v | %9v | %16v | %8v | %v" "" "" "" "" "" "" $dkimDomain $dkimResult "" "" "" }}
			{{- end }}
		{{- end }}
	{{- end }}
{{- end }}

{{ printf "Total: %v, passed %v, failed %v" .MessagesStats.All .MessagesStats.Passed .MessagesStats.Failed }}
`
