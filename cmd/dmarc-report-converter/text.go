package main

import (
	"html/template"
	"os"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

func textOutput(d dmarc.Report, c *config) error {
	t := template.Must(template.New("report").Parse(c.tmpl))
	err := t.Execute(os.Stdout, d)
	if err != nil {
		return err
	}

	return nil
}
