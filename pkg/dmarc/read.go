package dmarc

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	// MimeTypeGZIP is the mimetype for *.gz files
	MimeTypeGZIP = "application/x-gzip"
	// MimeTypeZIP is the mimetype for *.zip files
	MimeTypeZIP = "application/zip"
	// MimeTypeXML is the mimetype for *.xml files
	MimeTypeXML = "text/xml"
)

// ReadParseXML reads xml data from r and parses it to Report struct.
//
// If lookupAddr is true, performs reverse DNS lookups for all
// feedback>record>row>source_ip entries.
//
// lookupLimit is the maximum pool size for doing concurrent DNS lookups. Any
// lookupLimit value less than 1 will disable concurrency by setting the pool
// size to 1.
func ReadParseXML(r io.Reader, lookupAddr bool, lookupLimit int) (Report, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return Report{}, err
	}

	return Parse(data, lookupAddr, lookupLimit)
}

// ReadParseGZIP reads gzipped xml data from r and parses it to Report struct.
//
// If lookupAddr is true, performs reverse DNS lookups for all
// feedback>record>row>source_ip entries.
//
// lookupLimit is the maximum pool size for doing concurrent DNS lookups. Any
// lookupLimit value less than 1 will disable concurrency by setting the pool
// size to 1.
func ReadParseGZIP(r io.Reader, lookupAddr bool, lookupLimit int) (Report, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return Report{}, err
	}
	defer gr.Close()

	// in some cases inside gzip file can be places gzipped xml (gzipped twice)
	// unpack and check mimetype again
	buf := bytes.NewBuffer(nil)
	teer := io.TeeReader(gr, buf)
	data, err := io.ReadAll(teer)
	if err != nil {
		return Report{}, err
	}

	mtype := http.DetectContentType(data)
	if mtype == MimeTypeGZIP {
		log.Printf("[DEBUG] ReadParseGZIP: detected nested %v mimetype", mtype)
		return ReadParseGZIP(buf, lookupAddr, lookupLimit)
	} else if strings.HasPrefix(mtype, MimeTypeXML) {
		return ReadParseXML(buf, lookupAddr, lookupLimit)
	}

	return Report{}, fmt.Errorf("ReadParseGZIP: supported mimetypes not found")
}

// ReadParseZIP reads zipped xml data from r and parses it to Report struct.
//
// If lookupAddr is true, performs reverse DNS lookups for all
// feedback>record>row>source_ip entries.
//
// lookupLimit is the maximum pool size for doing concurrent DNS lookups. Any
// lookupLimit value less than 1 will disable concurrency by setting the pool
// size to 1.
func ReadParseZIP(r io.Reader, lookupAddr bool, lookupLimit int) (Report, error) {
	zipBytes, err := io.ReadAll(r)
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

		return ReadParseXML(rr, lookupAddr, lookupLimit)
	}

	return Report{}, err
}

// ReadParse reads any data from reader r, detects mimetype, and parses it to
// Report struct (if mimetype is supported).
//
// If lookupAddr is true, performs reverse DNS lookups for all
// feedback>record>row>source_ip entries.
//
// lookupLimit is the maximum pool size for doing concurrent DNS lookups. Any
// lookupLimit value less than 1 will disable concurrency by setting the pool
// size to 1.
func ReadParse(r io.Reader, lookupAddr bool, lookupLimit int) (Report, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return Report{}, err
	}

	mtype := http.DetectContentType(data)
	log.Printf("[DEBUG] ReadParse: detected %v mimetype", mtype)

	br := bytes.NewReader(data)
	if mtype == MimeTypeGZIP {
		return ReadParseGZIP(br, lookupAddr, lookupLimit)
	} else if mtype == MimeTypeZIP {
		return ReadParseZIP(br, lookupAddr, lookupLimit)
	} else if strings.HasPrefix(mtype, MimeTypeXML) {
		return ReadParseXML(br, lookupAddr, lookupLimit)
	}

	return Report{}, fmt.Errorf("mimetype %v not supported", mtype)
}
