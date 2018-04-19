package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

func textOutput(d dmarc.Report, cfg *config) error {
	var err error

	t := template.Must(template.New("report").Parse(cfg.tmpl))

	o := os.Stdout
	if cfg.outDir != "" {
		filepath := filepath.Join(cfg.outDir, d.ReportMetadata.Email+"!"+d.PolicyPublished.Domain+"!"+strconv.Itoa(d.ReportMetadata.DateRange.Begin)+"-"+strconv.Itoa(d.ReportMetadata.DateRange.Begin)+".txt")
		log.Printf("output to file %v\n", filepath)
		o, err = os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer o.Close()
	}

	err = t.Execute(o, d)
	if err != nil {
		return err
	}

	return nil
}
