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
	log.Printf("convert file %v, extension %v", i, ext)

	switch ext {
	case ".gz":
		report, err = readGZIP(r, cfg)
		if err != nil {
			return err
		}

	case ".zip":
		report, err = readZIP(r, cfg)
		if err != nil {
			return err
		}

	case ".xml":
		report, err = readXML(r, cfg)
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
