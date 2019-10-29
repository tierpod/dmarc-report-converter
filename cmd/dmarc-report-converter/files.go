package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

func processFiles(cfg *config) error {
	var reports []dmarc.Report

	inFiles, err := filepath.Glob(filepath.Join(cfg.Input.Dir, "*.*"))
	if err != nil {
		return err
	}

	log.Printf("[DEBUG]: files: found %v input files: %v", len(inFiles), inFiles)
	for _, f := range inFiles {
		file, err := os.Open(f)
		if err != nil {
			log.Printf("[ERROR] files: %v", err)
			continue
		}

		report, err := readParse(file, f, cfg.LookupAddr)
		if err != nil {
			file.Close()
			log.Printf("[ERROR] files: %v, skip", err)
			continue
		}
		file.Close()

		reports = append(reports, report)
	}

	mergedReports, err := groupMergeReports(reports)
	if err != nil {
		return err
	}

	for _, report := range mergedReports {
		o := newOutput(cfg)
		err = o.do(report)
		if err != nil {
			return err
		}
	}

	if cfg.Input.Delete {
		for _, f := range inFiles {
			log.Printf("[INFO] files: delete %v after converting", f)
			err = os.Remove(f)
			if err != nil {
				log.Printf("[ERROR] files: %v, skip", err)
				continue
			}
		}
	}

	return nil
}
