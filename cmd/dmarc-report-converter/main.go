package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var version string

func main() {
	var (
		flagVersion bool
		flagConfig  string
	)

	flag.BoolVar(&flagVersion, "version", false, "show version and exit")
	flag.StringVar(&flagConfig, "config", "./config.yaml", "Path to config file")
	flag.Parse()

	if flagVersion {
		fmt.Printf("Version: %v\n", version)
		os.Exit(0)
	}

	cfg, err := loadConfig(flagConfig)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	fmt.Printf("%+v\n", cfg)

	inFiles, err := filepath.Glob(filepath.Join(cfg.Input.Dir, "*.*"))
	if err != nil {
		log.Fatal(err)
	}

	if cfg.LookupAddr {
		log.Printf("performs a reverse lookups, this may take some time")
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
