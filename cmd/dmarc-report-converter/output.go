package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

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

	// if config Output.Dir is set, choose output file name and open file for writing
	if o.cfg.Output.Dir != "" {
		// generate output filename from config filename template
		var buf bytes.Buffer
		err := o.cfg.Output.fileTemplate.Execute(&buf, d)
		if err != nil {
			return err
		}

		file := filepath.Join(o.cfg.Output.Dir, buf.String())
		dir := filepath.Dir(file)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return nil
		}

		log.Printf("output to file %v\n", file)
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		o.w = f
		defer f.Close()
	}

	switch o.cfg.Output.Format {
	case "txt", "html":
		err = o.template(d)
	case "json":
		err = o.json(d)
	default:
		return fmt.Errorf("unknown output format %v", o.cfg.Output.Format)
	}

	return err
}

func (o *output) template(d dmarc.Report) error {
	err := o.cfg.Output.template.Execute(o.w, d)
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
