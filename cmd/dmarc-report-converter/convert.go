package main

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

func convertFile(i string, cfg *config) error {
	var d dmarc.Report

	ext := filepath.Ext(i)
	log.Printf("convert file %v, extension %v", i, ext)

	switch ext {
	case ".gz":
		file, err := os.Open(i)
		if err != nil {
			return err
		}
		defer file.Close()

		r, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer r.Close()

		data, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		d, err = dmarc.Parse(data, cfg.lookupAddr)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("extention %v not supported for file %v", ext, i)
	}

	o := newOutput(cfg)
	err := o.do(d)
	if err != nil {
		return err
	}

	return nil
}
