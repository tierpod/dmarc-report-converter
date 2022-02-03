package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type config struct {
	Input        Input  `yaml:"input"`
	Output       Output `yaml:"output"`
	LookupAddr   bool   `yaml:"lookup_addr"`
	MergeReports bool   `yaml:"merge_reports"`
	LogDebug     bool   `yaml:"log_debug"`
	LogDatetime  bool   `yaml:"log_datetime"`
}

// Input is the input section of config
type Input struct {
	Dir    string `yaml:"dir"`
	IMAP   IMAP   `yaml:"imap"`
	Delete bool   `yaml:"delete"`
}

// IMAP is the input.imap section of config
type IMAP struct {
	Server   string `yaml:"server"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Mailbox  string `yaml:"mailbox"`
	Debug    bool   `yaml:"debug"`
	Delete   bool   `yaml:"delete"`
}

// IsConfigured return true if IMAP is configured
func (i IMAP) IsConfigured() bool {
	return i.Server != ""
}

// Output is the output section of config
type Output struct {
	File             string `yaml:"file"`
	fileTemplate     *template.Template
	Format           string `yaml:"format"`
	Template         string
	template         *template.Template
	AssetsPath       string `yaml:"assets_path"`
	ExternalTemplate string `yaml:"external_template"`
}

func (o *Output) isStdout() bool {
	if o.File == "" || o.File == "stdout" {
		return true
	}

	return false
}

func loadConfig(path string) (*config, error) {
	var c config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	if c.Input.Dir == "" {
		return nil, fmt.Errorf("input.dir is not configured")
	}

	// Determine which template is used based upon Output.Format.
	t := txtTmpl
	switch c.Output.Format {
	case "txt":
	case "html":
		t = htmlTmpl
	case "html_static":
		t = htmlStaticTmpl
	case "external_template":
		if c.Output.ExternalTemplate == "" {
			return nil, fmt.Errorf("'output.external_template' must be configured to use external_template output")
		}
		data, err := ioutil.ReadFile(c.Output.ExternalTemplate)
		if err != nil {
			return nil, err
		}
		t = string(data)
	case "json":
	default:
		return nil, fmt.Errorf("unable to found template for format '%v' in templates folder", c.Output.Format)
	}
	tmplFuncs := template.FuncMap{
		"now": func(fmt string) string {
			return time.Now().Format(fmt)
		},
	}

	c.Output.template = template.Must(template.New("report").Funcs(tmplFuncs).Parse(t))

	if !c.Output.isStdout() {
		// load and parse output filename template
		ft := template.Must(template.New("filename").Funcs(tmplFuncs).Parse(c.Output.File))
		c.Output.fileTemplate = ft
	}

	return &c, nil
}
