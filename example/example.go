package main

import (
	"bytes"
	"fmt"
	"net/url"
	"runtime/debug"
	"strings"
	"text/template"
)

// Overwrite this to include template.FuncMap into the rendering
var ErrTemplateFuncs template.FuncMap

type ErrErrorGen interface {
	Error() string
	Wrap(error) ErrErrorGen
	Unwrap() error
	Data() map[string]interface{} // Returns Data used for rendering
	Name() string                 // Name of the error
	Stack() []byte                // Returns stack of the error
}

type ErrHTTPError interface {
	ErrErrorGen

	GetLang() string
	SetLang(string)

	GetStatus() int
	SetStatus(int)

	GetUrl() *url.URL
	SetUrl(*url.URL)
}

// ErrFileNotFound

type ErrFileNotFoundStruct struct {
	parent error
	params ErrFileNotFoundParams
	stack  []byte
}

type ErrFileNotFoundParams struct {
	File   string   // Description of the parameter
	Lang   string   //
	Status int      //
	Url    *url.URL //

}

// ErrFileNotFound returns a new instance of ErrFileNotFound with default values
func ErrFileNotFound() ErrFileNotFoundStruct {
	e := ErrFileNotFoundStruct{}
	e.stack = debug.Stack()

	e = e.Status(400)

	return e
}

func (e ErrFileNotFoundStruct) Name() string {
	return "FileNotFound"
}

func (e ErrFileNotFoundStruct) templ() (*template.Template, error) {
	return template.New("n").
		Funcs(ErrTemplateFuncs).
		Parse(strings.Trim("The {{ .Lang }} file \"{{ .File }}\" could not be found. {{ FuncTest \"Test\" }} Host: {{ .Url.Host }}\n", " \n"))
}

func (e ErrFileNotFoundStruct) Stack() []byte {
	return e.stack
}

func (e ErrFileNotFoundStruct) Error() string {
	t, err := e.templ()
	if err != nil {
		panic(fmt.Sprintf("Error compiling template: %q", err))
	}
	buf := bytes.Buffer{}
	err = t.Execute(&buf, e.Data())
	if err != nil {
		panic(fmt.Sprintf("Error executing template: %q", err))
	}
	return buf.String()
}

// GetFile returns the value of the key
func (e ErrFileNotFoundStruct) GetFile() string {
	return e.params.File
}

// File sets the value and returns a copy of the error (use for chaining)
func (e ErrFileNotFoundStruct) File(v string) ErrFileNotFoundStruct {
	e.params.File = v
	return e
}

// File sets the value in place
func (e *ErrFileNotFoundStruct) SetFile(v string) {
	e.params.File = v
}

// GetLang returns the value of the key
func (e ErrFileNotFoundStruct) GetLang() string {
	return e.params.Lang
}

// Lang sets the value and returns a copy of the error (use for chaining)
func (e ErrFileNotFoundStruct) Lang(v string) ErrFileNotFoundStruct {
	e.params.Lang = v
	return e
}

// Lang sets the value in place
func (e *ErrFileNotFoundStruct) SetLang(v string) {
	e.params.Lang = v
}

// GetStatus returns the value of the key
func (e ErrFileNotFoundStruct) GetStatus() int {
	return e.params.Status
}

// Status sets the value and returns a copy of the error (use for chaining)
func (e ErrFileNotFoundStruct) Status(v int) ErrFileNotFoundStruct {
	e.params.Status = v
	return e
}

// Status sets the value in place
func (e *ErrFileNotFoundStruct) SetStatus(v int) {
	e.params.Status = v
}

// GetUrl returns the value of the key
func (e ErrFileNotFoundStruct) GetUrl() *url.URL {
	return e.params.Url
}

// Url sets the value and returns a copy of the error (use for chaining)
func (e ErrFileNotFoundStruct) Url(v *url.URL) ErrFileNotFoundStruct {
	e.params.Url = v
	return e
}

// Url sets the value in place
func (e *ErrFileNotFoundStruct) SetUrl(v *url.URL) {
	e.params.Url = v
}

// Data returns all parameters as map
func (e ErrFileNotFoundStruct) Data() map[string]interface{} {
	data := map[string]interface{}{
		"File":   e.GetFile(),
		"Lang":   e.GetLang(),
		"Status": e.GetStatus(),
		"Url":    e.GetUrl(),
	}
	return data
}

// Wrap given error
func (e *ErrFileNotFoundStruct) Wrap(err error) ErrErrorGen {
	if e.parent != nil {
		panic("Unable to wrap ErrFileNotFound with already existing parent.")
	}
	e.parent = err
	return e
}

func (e ErrFileNotFoundStruct) Unwrap() error {
	return e.parent
}
