package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"bytes"

	"github.com/pkg/errors"
	"github.com/programmfabrik/errors/config"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func main() {
	conf := config.Config{}
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "input",
				Aliases:     []string{"i"},
				DefaultText: "-",
				Usage:       "Load .yml input from `FILE`. Use \"-\" to read from stdin.",
				Destination: &conf.InputFile,
			},
			&cli.StringFlag{
				Name:        "package",
				Aliases:     []string{"p"},
				DefaultText: "errors",
				Usage:       "Package name for the declaration.",
				Destination: &conf.Package,
			},
		},
		Name:  "errors",
		Usage: "generate Go error definitions from .yml input",
		Action: func(c *cli.Context) (err error) {

			var r io.ReadCloser

			switch conf.InputFile {
			case "-", "":
				r = os.Stdin
			default:
				r, err = os.Open(conf.InputFile)
			}
			if err != nil {
				return errors.Errorf("Unable to open %q", conf.InputFile)
			}

			inputBytes, err := ioutil.ReadAll(r)
			if err != nil {
				return errors.Errorf("Unable to read input")
			}

			r.Close()

			// parse the input as YAML
			err = yaml.Unmarshal(inputBytes, &conf.D)
			if err != nil {
				return errors.Wrapf(err, "Unable to parse input")
			}

			tmpl, err := template.ParseFiles("codegen/codegen.go.tmpl")
			if err != nil {
				return errors.Wrapf(err, "Unable to parse template")
			}

			buf := bytes.Buffer{}
			err = tmpl.Execute(&buf, conf)
			if err != nil {
				return errors.Wrapf(err, "Unable to execute template")
			}

			io.Copy(os.Stdout, &buf)

			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
