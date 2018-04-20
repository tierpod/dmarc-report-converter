package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

func convertFile(i string, cfg *config) error {
	var report dmarc.Report
	var err error

	ext := filepath.Ext(i)
	log.Printf("convert file %v, extension %v", i, ext)

	switch ext {
	case ".gz":
		report, err = readGZIP(i)
		if err != nil {
			return err
		}

	case ".zip":
		report, err = readZIP(i)
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
