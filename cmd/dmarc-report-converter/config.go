package main

import (
	"fmt"
	"io/ioutil"
)

type config struct {
	tmpl       string
	outDir     string
	outFormat  string
	lookupAddr bool
}

func newConfig(outDir, outFormat string, lookupAddr bool) (*config, error) {
	var t string
	var err error

	switch outFormat {
	case "txt":
		t, err = loadTemplate("./templates/txt.gotmpl")
		if err != nil {
			return nil, err
		}
	case "html":
	case "json":
	default:
		return nil, fmt.Errorf("unknown template for format %v", outFormat)
	}

	c := &config{
		tmpl:       t,
		outDir:     outDir,
		outFormat:  outFormat,
		lookupAddr: lookupAddr,
	}

	return c, nil
}

func loadTemplate(s string) (string, error) {
	data, err := ioutil.ReadFile(s)
	if err != nil {
		return "", err
	}

	return string(data), err
}
