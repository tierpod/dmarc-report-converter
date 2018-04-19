package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

type output struct {
	cfg *config
	w   io.Writer
}

func newOutput(cfg *config) *output {
	w := os.Stdout
	return &output{w: w, cfg: cfg}
}

func (o *output) do(d dmarc.Report) error {
	var err error

	// if -out is set, choose output file name and open file for writing
	if o.cfg.outDir != "" {
		filepath := filepath.Join(o.cfg.outDir, d.ReportMetadata.Email+"!"+d.PolicyPublished.Domain+"!"+strconv.Itoa(d.ReportMetadata.DateRange.Begin)+"-"+strconv.Itoa(d.ReportMetadata.DateRange.Begin)+"."+o.cfg.outFormat)
		log.Printf("output to file %v\n", filepath)
		f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		o.w = f
		defer f.Close()
	}

	switch o.cfg.outFormat {
	case "txt":
		err = o.txt(d)
	case "json":
		err = o.json(d)
	case "html":
		err = o.html(d)
	default:
		return fmt.Errorf("unknown output format %v", o.cfg.outFormat)
	}

	return err
}

func (o *output) txt(d dmarc.Report) error {
	fmt.Printf("%+v\n", o.cfg)
	t := template.Must(template.New("report").Parse(o.cfg.tmpl))
	err := t.Execute(o.w, d)
	if err != nil {
		return err
	}

	return nil
}

func (o *output) json(d dmarc.Report) error {
	js, err := json.Marshal(d)
	if err != nil {
		return err
	}

	fmt.Fprint(o.w, string(js))
	return nil
}

func (o *output) html(d dmarc.Report) error {
	return nil
}
