package main

import (
	"fmt"
	"io"
	"log"
	"path/filepath"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

func readParse(r io.Reader, fname string, lookupAddr bool) (dmarc.Report, error) {
	var report dmarc.Report
	var err error

	ext := filepath.Ext(fname)
	log.Printf("[DEBUG] parse: %v", fname)

	switch ext {
	case ".gz":
		report, err = dmarc.ReadParseGZIP(r, lookupAddr)
		if err != nil {
			return dmarc.Report{}, err
		}

	case ".zip":
		report, err = dmarc.ReadParseZIP(r, lookupAddr)
		if err != nil {
			return dmarc.Report{}, err
		}

	case ".xml":
		report, err = dmarc.ReadParseXML(r, lookupAddr)
		if err != nil {
			return dmarc.Report{}, err
		}

	default:
		return dmarc.Report{}, fmt.Errorf("extention %v not supported for file %v", ext, fname)
	}

	return report, nil
}
