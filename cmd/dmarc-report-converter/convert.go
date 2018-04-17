package main

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

func convertFile(f, o string) error {
	mt := mime.TypeByExtension(filepath.Ext(f))
	log.Printf("convert file %v, mimetype %v", f, mt)

	switch mt {
	case "application/gzip":
		file, err := os.Open(f)
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

		d, err := dmarc.Parse(data)
		if err != nil {
			return err
		}

		err = textOutput(d, o)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown mimetype %v for file %v", mt, f)
	}

	return nil
}
