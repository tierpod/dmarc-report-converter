package main

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

func readXML(r io.Reader, cfg *config) (dmarc.Report, error) {
	d, err := dmarc.ReadParse(r, cfg.LookupAddr)
	if err != nil {
		return dmarc.Report{}, err
	}

	return d, nil
}

func readGZIP(r io.Reader, cfg *config) (dmarc.Report, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return dmarc.Report{}, err
	}
	defer gr.Close()

	d, err := dmarc.ReadParse(gr, cfg.LookupAddr)
	if err != nil {
		return dmarc.Report{}, err
	}

	return d, nil
}

func readZIP(r io.Reader, cfg *config) (dmarc.Report, error) {
	zipBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return dmarc.Report{}, err
	}

	size := int64(len(zipBytes))
	readerAt := bytes.NewReader(zipBytes)

	zr, err := zip.NewReader(readerAt, size)
	if err != nil {
		return dmarc.Report{}, err
	}

	for _, file := range zr.File {
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

		d, err := dmarc.ReadParse(rr, cfg.LookupAddr)
		if err != nil {
			return dmarc.Report{}, err
		}

		return d, nil
	}

	return dmarc.Report{}, err
}
