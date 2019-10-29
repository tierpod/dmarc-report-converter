// Convert dmarc reports from xml to human-readable formats
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
		log.Fatalf("[ERROR] loadConfig: %v", err)
	}

	if cfg.LookupAddr {
		log.Printf("[INFO] performs a reverse lookups, this may take some time")
	}

	if cfg.Input.IMAP.IsConfigured() {
		err = fetchIMAPAttachments(cfg)
		if err != nil {
			log.Fatalf("[ERROR] fetchIMAPAttachments: %v", err)
		}
	}

	processFiles(cfg)
}
