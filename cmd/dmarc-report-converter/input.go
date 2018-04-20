package main

import (
	"archive/zip"
	"compress/gzip"
	"log"
	"os"
	"path/filepath"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

// func readXML(i string) {

// }

func readGZIP(i string) (dmarc.Report, error) {
	file, err := os.Open(i)
	if err != nil {
		return dmarc.Report{}, err
	}
	defer file.Close()

	r, err := gzip.NewReader(file)
	if err != nil {
		return dmarc.Report{}, err
	}
	defer r.Close()

	d, err := dmarc.ReadParse(r, cfg.lookupAddr)
	if err != nil {
		return dmarc.Report{}, err
	}

	return d, nil
}

func readZIP(i string) (dmarc.Report, error) {
	r, err := zip.OpenReader(i)
	if err != nil {
		return dmarc.Report{}, err
	}
	defer r.Close()

	for _, file := range r.File {
		if filepath.Ext(file.Name) != ".xml" {
			log.Printf("[WARN] skip %v from zip: unknown extension\n", file.Name)
			continue
		}

		log.Printf("read file %v from zip\n", file.Name)

		rr, err := file.Open()
		if err != nil {
			return dmarc.Report{}, err
		}
		defer rr.Close()

		d, err := dmarc.ReadParse(rr, cfg.lookupAddr)
		if err != nil {
			return dmarc.Report{}, err
		}

		return d, nil
	}

	return dmarc.Report{}, err
}
