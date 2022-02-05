package dmarc

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// ReadParseXML reads xml data from r and parses it to Report struct. If lookupAddr is
// true, performs a reverse lookups for feedback>record>row>source_ip
func ReadParseXML(r io.Reader, lookupAddr bool) (Report, error) {
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

func unpackGZIP(b []byte) ([]byte, error) {
	r := bytes.NewReader(b)
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	data, err := ioutil.ReadAll(gr)
	if err != nil {
		return nil, err
	}

	return data, err
}

// ReadParseGZIP reads gzipped xml data from r and parses it to Report struct. If lookupAddr is
// true, performs a reverse lookups for feedback>record>row>source_ip
func ReadParseGZIP(r io.Reader, lookupAddr bool) (Report, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return Report{}, err
	}
	defer gr.Close()

	d, err := ReadParseXML(gr, lookupAddr)
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

		d, err := ReadParseXML(rr, lookupAddr)
		if err != nil {
			return Report{}, err
		}

		return d, nil
	}

	return Report{}, err
}

// ReadParse reads any data from reader r, detects mimetype, and parses it to
// Report struct (if mimetype is supported).
// If lookupAddr is true, performs reverse lookups for feedback>record>row>source_ip
func ReadParse(r io.Reader, lookupAddr bool) (Report, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return Report{}, err
	}

	var report Report
	mtype := http.DetectContentType(data)
	log.Printf("[DEBUG] ReadParse: detected %v mimetype", mtype)

	br := bytes.NewReader(data)
	if mtype == "application/x-gzip" {
		// in some cases inside gzip file can be places gzipped xml (gzipped twice)
		// unpack and check mimetype again
		gzdata, err := unpackGZIP(data)
		if err != nil {
			return Report{}, err
		}
		gzmtype := http.DetectContentType(gzdata)
		if gzmtype == "application/x-gzip" {
			log.Printf("[DEBUG] ReadParse: detected nested %v mimetype", gzmtype)
			br = bytes.NewReader(gzdata)
		}
		report, err = ReadParseGZIP(br, lookupAddr)
		if err != nil {
			return Report{}, err
		}
	} else if mtype == "application/zip" {
		report, err = ReadParseZIP(br, lookupAddr)
		if err != nil {
			return Report{}, err
		}
	} else if strings.HasPrefix(mtype, "text/xml") {
		report, err = ReadParseXML(br, lookupAddr)
		if err != nil {
			return Report{}, err
		}
	} else {
		return Report{}, fmt.Errorf("mimetype %v not supported", mtype)
	}

	return report, nil
}
