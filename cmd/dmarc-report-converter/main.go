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
		//	flagConfig  string
		flagInputDir  string
		flagOutputDir string
	)

	flag.BoolVar(&flagVersion, "version", false, "Show version and exit")
	flag.StringVar(&flagInputDir, "i", "", "Path to input directory")
	flag.StringVar(&flagOutputDir, "o", "", "Path to output directory")
	//flag.StringVar(&flagConfig, "config", "./config.yaml", "Path to config file")
	flag.Parse()

	if flagVersion {
		fmt.Printf("Version: %v\n", version)
		os.Exit(0)
	}

	if flagInputDir == "" || flagOutputDir == "" {
		log.Fatal("-i or -o is not set")
	}

	inFiles, err := filepath.Glob(filepath.Join(flagInputDir, "*.*"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("found input files: %v", inFiles)
	for _, f := range inFiles {
		err = convertFile(f, flagOutputDir)
		if err != nil {
			log.Printf("[ERROR] %v", err)
		}
	}
}
