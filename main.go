package main

//go:generate esc -private -local-prefix-cwd -pkg=main -o=resources.go template/

import (
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"bytes"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func main() {
	conf := Config{}
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
				Name:        "output",
				Aliases:     []string{"o"},
				DefaultText: "-",
				Usage:       "Write generated code to `FILE`. Use \"-\" to output to stdout.",
				Destination: &conf.OutputFile,
			},
			&cli.StringFlag{
				Name:        "tmpl",
				Aliases:     []string{"t"},
				Usage:       "Use `FILE` as code template for generating.",
				Destination: &conf.Template,
			},
		},
		Name:  "errors",
		Usage: "generate Go error definitions from .yml input",
		Action: func(c *cli.Context) (err error) {

			var in io.ReadCloser

			switch conf.InputFile {
			case "-", "":
				in = os.Stdin
			default:
				in, err = os.Open(conf.InputFile)
			}
			if err != nil {
				return errors.Errorf("Unable to open %q", conf.InputFile)
			}

			inputBytes, err := ioutil.ReadAll(in)
			if err != nil {
				return errors.Errorf("Unable to read input")
			}

			in.Close()

			// parse the input as YAML
			err = yaml.Unmarshal(inputBytes, &conf.D)
			if err != nil {
				return errors.Wrapf(err, "Unable to parse input")
			}

			err = conf.D.Validate()
			if err != nil {
				return err
			}

			var tmpl *template.Template

			if conf.Template != "" {
				tmpl, err = template.ParseFiles(conf.Template)
				if err != nil {
					return errors.Wrapf(err, "Unable to parse template")
				}
			} else {
				data, err := _escFSByte(false, "/template/codegen.go.tmpl")
				if err != nil {
					panic(err)
				}
				tmpl, err = template.New("n").Parse(string(data))
				if err != nil {
					return errors.Wrapf(err, "Unable to parse template")
				}
			}

			buf := bytes.Buffer{}
			err = tmpl.Execute(&buf, conf)
			if err != nil {
				return errors.Wrapf(err, "Unable to execute template")
			}

			codeBytes, errFormat := format.Source(buf.Bytes())
			if errFormat != nil {
				codeBytes = buf.Bytes()
			}

			switch conf.OutputFile {
			case "-", "":
				io.Copy(os.Stdout, bytes.NewReader(codeBytes))
			default:
				err = ioutil.WriteFile(conf.OutputFile, codeBytes, 0644)
				if err != nil {
					return errors.Wrapf(err, "Unable to write output file %q", conf.OutputFile)
				}
			}

			if errFormat != nil {
				return errors.Wrapf(err, "Unable to format Go code")
			}

			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
