package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	htemplate "html/template"
	ttemplate "text/template"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func runGen(c *cli.Context, conf Config) (err error) {
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

	return nil

}

func render(conf Config) (parsedBytes []byte, err error) {
	buf := bytes.Buffer{}
	var out string

	if conf.HTML {
		// Render html/template
		var tmpl *htemplate.Template

		if conf.Template != "" {
			tmpl, err = htemplate.ParseFiles(conf.Template)
			if err != nil {
				return nil, errors.Wrap(err, "Unable to parse template")
			}
			out = fmt.Sprintf("Rendered html/template %q", conf.Template)
		} else {
			data, err := _escFSByte(false, "/templates/doc.html.tmpl")
			if err != nil {
				return nil, err
			}
			tmpl, err = htemplate.New("n").Parse(string(data))
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

	var tmpl *ttemplate.Template

	if conf.Template != "" {
		tmpl, err = ttemplate.ParseFiles(conf.Template)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to parse template")
		}
		out = fmt.Sprintf("Rendered text/template %q", conf.Template)
	} else {
		data, err := _escFSByte(false, "/templates/codegen.go.tmpl")
		if err != nil {
			return nil, err
		}
		tmpl, err = ttemplate.New("n").Parse(string(data))
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to parse template")
		}
		out = fmt.Sprintf("Rendered text/template built-in")
	}

	err = tmpl.Execute(&buf, conf)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to execute template")
	}

	codeBytes, err := format.Source(buf.Bytes())

	fmt.Fprintf(os.Stderr, "%s: %d bytes. Error: %v\n", out, len(codeBytes), err)

	if err != nil {
		return codeBytes, errors.Wrapf(err, "Unable to format Go code")
	}
	return codeBytes, nil

}
