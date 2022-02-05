package main

import (
	"io"
	"log"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

func readParse(r io.Reader, fname string, lookupAddr bool) (dmarc.Report, error) {
	var report dmarc.Report
	var err error

	log.Printf("[DEBUG] parse: %v", fname)

	report, err = dmarc.ReadParse(r, lookupAddr)
	if err != nil {
		return dmarc.Report{}, err
	}
	return report, nil
}
