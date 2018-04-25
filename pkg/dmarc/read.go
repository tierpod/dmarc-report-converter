package dmarc

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
)

// readParse reads all from reader and parses it to Report struct
func readParse(r io.Reader, lookupAddr bool) (Report, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return Report{}, err
	}

	report, err := Parse(data, lookupAddr)
	if err != nil {
		return Report{}, err
	}

	return report, nil
}

// ReadParseXML reads xml data from r and parses it to Report struct. If lookupAddr is
// true, performs a reverse lookups for feedback>record>row>source_ip
func ReadParseXML(r io.Reader, lookupAddr bool) (Report, error) {
	d, err := readParse(r, lookupAddr)
	if err != nil {
		return Report{}, err
	}

	return d, nil
}

// ReadParseGZIP reads gzipped xml data from r and parses it to Report struct. If lookupAddr is
// true, performs a reverse lookups for feedback>record>row>source_ip
func ReadParseGZIP(r io.Reader, lookupAddr bool) (Report, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return Report{}, err
	}
	defer gr.Close()

	d, err := readParse(gr, lookupAddr)
	if err != nil {
		return Report{}, err
	}

	return d, nil
}

// ReadParseZIP reads zipped xml data from r and parses it to Report struct. If lookupAddr is
// true, performs a reverse lookups for feedback>record>row>source_ip
func ReadParseZIP(r io.Reader, lookupAddr bool) (Report, error) {
	zipBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return Report{}, err
	}

	size := int64(len(zipBytes))
	readerAt := bytes.NewReader(zipBytes)

	zr, err := zip.NewReader(readerAt, size)
	if err != nil {
		return Report{}, err
	}

	for _, file := range zr.File {
		ext := filepath.Ext(file.Name)
		if ext != ".xml" {
			log.Printf("[WARN] ReadParseZIP: skip %v from zip: unknown extension %v", file.Name, ext)
			continue
		}

		log.Printf("[INFO] ReadParseZIP: read file %v from zip", file.Name)

		rr, err := file.Open()
		if err != nil {
			return Report{}, err
		}
		defer rr.Close()

		d, err := readParse(rr, lookupAddr)
		if err != nil {
			return Report{}, err
		}

		return d, nil
	}

	return Report{}, err
}
