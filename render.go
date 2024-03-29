package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	htemplate "html/template"
	ttemplate "text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"

	_ "embed"
)

//go:embed templates/doc.html.tmpl
var docTmpl string

//go:embed templates/codegen.go.tmpl
var codeTmpl string

func runGen(c *cli.Context, conf Config) (err error) {
	var in io.ReadCloser

	switch conf.InputFile {
	case "-", "":
		in = os.Stdin
		conf.inputDebug = "<stdin>"
	default:

		wDir, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "Unable to get workdir")
		}

		in, err = os.Open(conf.InputFile)
		conf.inputDebug = filepath.Join(wDir, conf.InputFile)
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

	err = conf.Validate()
	if err != nil {
		return err
	}

	rendered, err := render(conf)
	if conf.Server > 0 {
		if err != nil && len(rendered) == 0 {
			return err
		}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			bytes, err := render(conf)
			if err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", err.Error()))
				w.Write([]byte("<pre>" + err.Error() + "</pre>"))
				return
			}
			w.Write(bytes)
		})
		fmt.Fprintf(os.Stderr, "Starting webserver on :%d\n", conf.Server)
		http.ListenAndServe(":"+strconv.Itoa(conf.Server), nil)
		return
	}

	// Even with error, output bytes. In case go format fails, we want to see the code
	if len(rendered) > 0 {
		switch conf.OutputFile {
		case "-", "":
			io.Copy(os.Stdout, bytes.NewReader(rendered))
		default:
			err = ioutil.WriteFile(conf.OutputFile, rendered, 0644)
			if err != nil {
				return errors.Wrapf(err, "Unable to write output file %q", conf.OutputFile)
			}
		}
	}
	if err != nil {
		return err
	}

	if conf.CsvFile != "" {
		conf.D.writeCSV(conf.CsvFile)
	}

	return nil

}

func render(conf Config) (parsedBytes []byte, err error) {
	buf := bytes.Buffer{}
	var out string

	if conf.HTML {
		// Render html/template
		tmpl := htemplate.New("n").Funcs(sprig.HtmlFuncMap())
		if conf.Template != "" {
			tmpl, err = tmpl.ParseFiles(conf.Template)
			if err != nil {
				return nil, errors.Wrap(err, "Unable to parse template")
			}
			out = fmt.Sprintf("Rendered html/template %q", conf.Template)
		} else {
			tmpl, err = tmpl.Parse(docTmpl)
			if err != nil {
				return nil, errors.Wrap(err, "Unable to parse template")
			}
			out = fmt.Sprintf("Rendered html/template built-in")
		}

		err = tmpl.Execute(&buf, conf)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to execute template")
		}
		fmt.Fprintf(os.Stderr, "%s: %d bytes. Error: %v\n", out, buf.Len(), err)
		return buf.Bytes(), nil
	}

	// Render text/template

	tmpl := ttemplate.New("n").Funcs(sprig.TxtFuncMap())
	if conf.Template != "" {
		tmpl, err = tmpl.ParseFiles(conf.Template)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to parse template")
		}
		out = fmt.Sprintf("text/template %q", conf.Template)
	} else {
		tmpl, err = tmpl.Parse(codeTmpl)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to parse template")
		}
		out = fmt.Sprintf("text/template built-in")
	}
	err = tmpl.Execute(&buf, conf)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to execute template")
	}

	codeBytes, err := format.Source(buf.Bytes())

	if err != nil {
		fmt.Fprintf(os.Stderr, "[errorgen] %s: %s: %d bytes. Error: %v\n", conf.inputDebug, out, len(codeBytes), err)
	} else {
		fmt.Fprintf(os.Stderr, "[errorgen] %s: %s: %d bytes.\n", conf.inputDebug, out, len(codeBytes))
	}

	if err != nil {
		return codeBytes, errors.Wrapf(err, "Unable to format Go code")
	}
	return codeBytes, nil

}
