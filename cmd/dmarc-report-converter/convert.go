package main

import (
	"fmt"
	"io"
	"log"
	"path/filepath"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

func readConvert(r io.Reader, i string, cfg *config) error {
	var report dmarc.Report
	var err error

	ext := filepath.Ext(i)
	log.Printf("[INFO] convert: file %v", i)

	switch ext {
	case ".gz":
		report, err = dmarc.ReadParseGZIP(r, cfg.LookupAddr)
		if err != nil {
			return err
		}

	case ".zip":
		report, err = dmarc.ReadParseZIP(r, cfg.LookupAddr)
		if err != nil {
			return err
		}

	case ".xml":
		report, err = dmarc.ReadParseXML(r, cfg.LookupAddr)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("extention %v not supported for file %v", ext, i)
	}

	o := newOutput(cfg)
	err = o.do(report)
	if err != nil {
		return err
	}

	return nil
}
