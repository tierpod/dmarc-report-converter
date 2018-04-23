package main

import (
	"log"
	"os"
	"path/filepath"
)

func processFiles(cfg *config) {
	inFiles, err := filepath.Glob(filepath.Join(cfg.Input.Dir, "*.*"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[DEBUG]: files: found input files: %v", inFiles)
	for _, f := range inFiles {
		file, err := os.Open(f)
		if err != nil {
			log.Printf("[ERROR] files: %v", err)
			continue
		}
		defer file.Close()

		err = readConvert(file, f, cfg)
		if err != nil {
			log.Printf("[ERROR] files: %v, skip", err)
			continue
		}
		file.Close()

		if cfg.Input.Delete {
			log.Printf("[INFO] files: delete %v after converting", f)
			err = os.Remove(f)
			if err != nil {
				log.Printf("[ERROR] files: %v, skip", err)
				continue
			}
		}
	}
}
