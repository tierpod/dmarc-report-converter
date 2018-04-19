package main

import "io/ioutil"

type config struct {
	tmpl       string
	lookupAddr bool
}

func newConfig(outFormat string, lookupAddr bool) (*config, error) {
	var t string
	var err error
	switch outFormat {
	case "text":
		t, err = loadTemplate("./templates/text.gotmpl")
		if err != nil {
			return nil, err
		}
	}

	c := &config{
		tmpl:       t,
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
