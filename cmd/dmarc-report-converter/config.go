package main

import (
	"fmt"
	"html/template"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type config struct {
	Input      Input  `yaml:"input"`
	Output     Output `yaml:"output"`
	LookupAddr bool   `yaml:"lookup_addr"`
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
}

// Output is the output section of config
type Output struct {
	File         string `yaml:"file"`
	fileTemplate *template.Template
	Format       string `yaml:"format"`
	Template     string
	template     *template.Template
}

func (o *Output) isStdout() bool {
	if o.File == "" || o.File == "stdout" {
		return true
	}

	return false
}

func loadTemplate(s string) *template.Template {
	data, err := ioutil.ReadFile(s)
	if err != nil {
		panic(err)
	}

	t := template.Must(template.New("report").Parse(string(data)))
	return t
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

	// load and parse output file template
	switch c.Output.Format {
	case "txt", "html":
		t := loadTemplate("./templates/" + c.Output.Format + ".gotmpl")
		c.Output.template = t
	case "json":
	default:
		return nil, fmt.Errorf("unknown template for format %v", c.Output.Format)
	}

	if !c.Output.isStdout() {
		// load and parse output filename template
		ft := template.Must(template.New("filename").Parse(c.Output.File))
		c.Output.fileTemplate = ft
	}

	return &c, nil
}
