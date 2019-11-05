UPGRADING instructions
======================

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

* config: added {{. TodayID }} shortcut for filename
