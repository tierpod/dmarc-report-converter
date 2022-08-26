package main

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
    <script src="{{.AssetsPath}}/js/Chart.min.js"></script>
    <script src="{{.AssetsPath}}/js/dmarc-report-converter.js"></script>

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
                <center><span class="badge bg-success">passed {{ .Report.MessagesStats.Passed }}</span> <span class="badge bg-danger">failed {{ .Report.MessagesStats.Failed }}</span> <span class="badge">total {{ .Report.MessagesStats.All }}</span></center>
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
                        <td class="left-border">{{.AuthResults.DKIM.Domain}}</td>
                        <td title="selector: {{.AuthResults.DKIM.Selector}}">
                            {{- if eq .AuthResults.DKIM.Result "pass"}}
                            <span class="badge bg-success">{{.AuthResults.DKIM.Result}}</span>
                            {{- else if eq .AuthResults.DKIM.Result "fail"}}
                            <span class="badge bg-danger">{{.AuthResults.DKIM.Result}}</span>
                            {{- else}}
                            <span class="badge bg-warning">{{.AuthResults.DKIM.Result}}</span>
                            {{- end}}
                        </td>
                        <td>{{.AuthResults.SPF.Domain}}</td>
                        <td title="scope: {{.AuthResults.SPF.Scope}}">
                            {{- if eq .AuthResults.SPF.Result "pass"}}
                            <span class="badge bg-success">{{.AuthResults.SPF.Result}}</span>
                            {{- else if eq .AuthResults.SPF.Result "fail"}}
                            <span class="badge bg-danger">{{.AuthResults.SPF.Result}}</span>
                            {{- else}}
                            <span class="badge bg-warning">{{.AuthResults.SPF.Result}}</span>
                            {{- end}}
                        </td>
                        </tr>
                        {{- end }}
                    </tbody>
                </table>
            </div>
        </div>
    </div>
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
    <link href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
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
                <div class="progress">
                    <div class="progress-bar bg-success" role="progressbar" style="width: {{ .Report.MessagesStats.PassedPercent }}%" aria-valuenow="{{ .Report.MessagesStats.Passed }}" aria-valuemin="0" aria-valuemax="{{ .Report.MessagesStats.All }}">passed {{ .Report.MessagesStats.PassedPercent }}%</div>
                </div>
            </div>
            <div class="col-md-auto">
                <span class="badge badge-success">passed {{ .Report.MessagesStats.Passed }}</span> <span class="badge badge-danger">failed {{ .Report.MessagesStats.Failed }}</span> <span class="badge">total {{ .Report.MessagesStats.All }}</span>
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
                            <span class="badge badge-danger">{{.Row.PolicyEvaluated.DKIM}}</span>
                            {{- else }}
                            <span class="badge badge-success">{{.Row.PolicyEvaluated.DKIM}}</span>
                            {{- end}}
                        </td>
                        <td>
                            {{- if eq .Row.PolicyEvaluated.SPF "fail" }}
                            <span class="badge badge-danger">{{.Row.PolicyEvaluated.SPF}}</span>
                            {{- else }}
                            <span class="badge badge-success">{{.Row.PolicyEvaluated.SPF}}</span>
                            {{- end }}
                        </td>
                        <td class="left-border">{{.AuthResults.DKIM.Domain}}</td>
                        <td title="selector: {{.AuthResults.DKIM.Selector}}">
                            {{- if eq .AuthResults.DKIM.Result "pass"}}
                            <span class="badge badge-success">{{.AuthResults.DKIM.Result}}</span>
                            {{- else if eq .AuthResults.DKIM.Result "fail"}}
                            <span class="badge badge-danger">{{.AuthResults.DKIM.Result}}</span>
                            {{- else}}
                            <span class="badge badge-warning">{{.AuthResults.DKIM.Result}}</span>
                            {{- end}}
                        </td>
                        <td>{{.AuthResults.SPF.Domain}}</td>
                        <td title="scope: {{.AuthResults.SPF.Scope}}">
                            {{- if eq .AuthResults.SPF.Result "pass"}}
                            <span class="badge badge-success">{{.AuthResults.SPF.Result}}</span>
                            {{- else if eq .AuthResults.SPF.Result "fail"}}
                            <span class="badge badge-danger">{{.AuthResults.SPF.Result}}</span>
                            {{- else}}
                            <span class="badge badge-warning">{{.AuthResults.SPF.Result}}</span>
                            {{- end}}
                        </td>
                        </tr>
                        {{- end }}
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</body>
</html>`

const txtTmpl = `
DMARC report with id {{.ReportMetadata.ReportID}}
  Organization:     {{.ReportMetadata.ExtraContactInfo}} ({{.ReportMetadata.Email}}))
  Date range:       since {{.ReportMetadata.DateRange.Begin.UTC}} until {{.ReportMetadata.DateRange.End.UTC}}
  Policy published: {{.PolicyPublished.Domain}}: p={{.PolicyPublished.Policy}} sp={{.PolicyPublished.SPolicy}} pct={{.PolicyPublished.Pct}} adkim={{.PolicyPublished.ADKIM}} aspf={{.PolicyPublished.ASPF}}
------------------------------------------------------------------------------------------------------------------------
{{printf "%23v | %32v | %57v |" "" "policy evaluated" "auth results" }}
------------------------------------------------------------------------------------------------------------------------
{{printf "%16v | %4v | %10v | %8v | %8v | %16v | %8v | %16v | %8v | %v" "ip"          "msgs"     "disp"                           "dkim"                    "spf"                    "dkim domain"            "dkim res"               "spf domain"            "spf res"               "hostname"         }}
------------------------------------------------------------------------------------------------------------------------
{{- range .Records }}
{{- if .IsPassed }}
{{ printf "* %14v | %4v | %10v | %8v | %8v | %16v | %8v | %16v | %8v | %v" .Row.SourceIP .Row.Count .Row.PolicyEvaluated.Disposition .Row.PolicyEvaluated.DKIM .Row.PolicyEvaluated.SPF .AuthResults.DKIM.Domain .AuthResults.DKIM.Result .AuthResults.SPF.Domain .AuthResults.SPF.Result .Row.SourceHostname}}
{{- else }}
{{ printf "%16v | %4v | %10v | %8v | %8v | %16v | %8v | %16v | %8v | %v" .Row.SourceIP .Row.Count .Row.PolicyEvaluated.Disposition .Row.PolicyEvaluated.DKIM .Row.PolicyEvaluated.SPF .AuthResults.DKIM.Domain .AuthResults.DKIM.Result .AuthResults.SPF.Domain .AuthResults.SPF.Result .Row.SourceHostname}}
{{- end}}
{{- end }}

{{ printf "Total: %v, passed %v, failed %v" .MessagesStats.All .MessagesStats.Passed .MessagesStats.Failed }}
`
