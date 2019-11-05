package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

type filesConverter struct {
	cfg     *config
	files   []string
	reports []dmarc.Report
}

func newFilesConverter(cfg *config) (*filesConverter, error) {
	if _, err := os.Stat(cfg.Input.Dir); os.IsNotExist(err) {
		err := os.MkdirAll(cfg.Input.Dir, 0775)
		if err != nil {
			return nil, err
		}
	}

	return &filesConverter{cfg: cfg}, nil
}

func (c *filesConverter) ConvertWrite() error {
	err := c.find()
	if err != nil {
		return err
	}

	c.convert()

	if c.cfg.MergeReports {
		err = c.merge()
		if err != nil {
			return err
		}
	}

	c.write()

	if c.cfg.Input.Delete {
		c.delete()
	}

	return nil
}

func (c *filesConverter) find() error {
	files, err := filepath.Glob(filepath.Join(c.cfg.Input.Dir, "*.*"))
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] files: found %v input files", len(files))
	c.files = files
	return nil
}

func (c *filesConverter) convert() {
	var reports []dmarc.Report
	for _, f := range c.files {
		file, err := os.Open(f)
		if err != nil {
			log.Printf("[ERROR] files: %v", err)
			continue
		}

		report, err := readParse(file, f, c.cfg.LookupAddr)
		if err != nil {
			file.Close()
			log.Printf("[ERROR] files: %v, skip", err)
			continue
		}
		file.Close()

		reports = append(reports, report)
	}

	c.reports = reports
}

func (c *filesConverter) merge() error {
	reports, err := groupMergeReports(c.reports)
	if err != nil {
		return err
	}

	c.reports = reports
	return nil
}

func (c *filesConverter) delete() {
	for _, f := range c.files {
		log.Printf("[INFO] files: delete %v after converting", f)
		err := os.Remove(f)
		if err != nil {
			log.Printf("[ERROR] files: %v, skip", err)
			continue
		}
	}
}

func (c *filesConverter) write() error {
	for _, report := range c.reports {
		o := newOutput(c.cfg)
		err := o.do(report)
		if err != nil {
			return err
		}
	}
	return nil
}
