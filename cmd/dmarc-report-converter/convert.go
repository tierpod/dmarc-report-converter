package main

import (
	"io"
	"log"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

// readParse is a helper function that passes r, lookupAddr, and lookupLimit to
// dmarc.ReadParse.
//
// fname is the file name associated with r and is only used for debug logging.
func readParse(r io.Reader, fname string, lookupAddr bool, lookupLimit int) (dmarc.Report, error) {
	log.Printf("[DEBUG] parse: %v", fname)
	return dmarc.ReadParse(r, lookupAddr, lookupLimit)
}
