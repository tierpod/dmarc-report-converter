// Convert dmarc reports from xml to human-readable formats
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/emersion/go-message/charset"
	"github.com/hashicorp/logutils"
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

	setupLog(cfg.LogDebug, cfg.LogDatetime)

	converter, err := newFilesConverter(cfg)
	if err != nil {
		log.Fatalf("[ERROR] newFilesConverter: %v", err)
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

	err = converter.ConvertWrite()
	if err != nil {
		log.Fatalf("[ERROR] processFiles: %v", err)
	}
}

func setupLog(debug, datetime bool) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stderr,
	}

	if debug {
		filter.MinLevel = logutils.LogLevel("DEBUG")
	}

	if datetime {
		log.SetFlags(log.LstdFlags)
	} else {
		log.SetFlags(0)
	}

	log.SetOutput(filter)
}
