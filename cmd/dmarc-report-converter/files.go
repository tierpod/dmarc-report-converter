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

	log.Printf("found input files: %v", inFiles)
	for _, f := range inFiles {
		err = convertFile(f, cfg)
		if err != nil {
			log.Printf("[ERROR] %v, skip", err)
			continue
		}

		if cfg.Input.Delete {
			log.Printf("delete %v after converting", f)
			err = os.Remove(f)
			if err != nil {
				log.Printf("[ERROR] %v, skip", err)
				continue
			}
		}
	}
}
