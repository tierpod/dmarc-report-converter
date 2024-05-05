UPGRADING instructions
======================

Migration to 0.8
----------------

There are some breaking changes in 0.8:

* *external_template* now uses the same template context as other kind of templates. It makes all
  templates consistent. If you use external_template, just replace `.` with `.Report`. See const.go
  for example, more details in PR #47, thanks to @moorereason)

* json output format changed, because `feedback>record>auth_results` dkim and spf can have
  multiple results (list of results). More details in PR #51, thanks to @moorereason)

Migration to 0.6
----------------

Since v0.6, all templates moved from templates/*.gotmpl files to consts.go (thanks to @morrowc).
This change makes installation and usage easier - you don't have to install and update templates
folder anymore. Inside external template dmarc.Report struct can be used (see consts.go -> txtTmpl
for example).

If you prefer to use self-written external templates, you can still do this:

```yaml
output:
  format: "external_template"
  external_template: "/path/to/your/txt.gotmpl"
```

* config: added *output -> format -> external_template* format

* config: added *output -> external_template* option

* deleted templates folder

Migration to 0.5
----------------

Since v0.5, dmarc-report-converter can merge similar reports, so it can be used to make weekly and
daily reports. Some internal structures was modified, so you have to update templates. If IMAP was
configured, you have to update configuration file.

* imap: before 0.5, dmarc-report-converter fetched all attachments into memory and then converted
  them. Since 0.5, it downloads attachments to local directory. You have to set *input -> dir*.

* templates: Result.Record was renamed to Result.Records

* templates: Report.Total* was moved to Report.MessagesStats.*

* config: *imap_debug* was moved to *input -> imap -> debug*

* config: added *input -> imap -> delete*

* config: added *log_datetime* and *log_debug* parameters

* config: added {{ .TodayID }} shortcut for filename
