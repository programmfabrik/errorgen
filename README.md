# errors
Go error generator with simple localization support

## Usage

```bash
./errors -i simple.yml -o test.go
```

## Example .yml file

```yaml
package: main
import:
  - "net/url"
defaults:
  HTTPError:
    Lang:
      t: string
    Status:
      t: int
      a: |
        json:"s"
    Url:
      t: "*url.URL"
prefix: Err
errors:
  FileNotFound:
    d: Description of the error
    p:
      File:
        d: Description of the parameter
        t: string
      Url:
        t: "*url.URL"
      Status:
        v: 400
    o: |
      The {{ .Lang }} file "{{ .File }}" could not be found. {{ FuncTest "Test" }} Host: {{ .Url.Host }} Url: {{ .Url }}
```

## Generated code

```go
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
	Unwrap() error
	Params() interface{} // Returns params used for rendering
	ErrorCode() string   // Name / code of the error
	Stack() []byte       // Returns stack of the error
}

// ErrHTTPError
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

type ErrFileNotFoundError struct {
	parent error
	params ErrFileNotFoundParams
	stack  []byte
}

type ErrFileNotFoundParams struct {
	File   string `json:"file"` // Description of the parameter
	Lang   string `json:"lang"` // Lang
	Status int    `json:"s"
` // Status
	Url *url.URL `json:"url"` // Url

}

// ErrFileNotFound returns a new instance of ErrFileNotFound with default values
func ErrFileNotFound() *ErrFileNotFoundError {
	e := &ErrFileNotFoundError{}
	e.stack = debug.Stack()

	e.Status(400)

	return e
}

func (e ErrFileNotFoundError) ErrorCode() string {
	return "FileNotFound"
}

func (e ErrFileNotFoundError) templ() (*template.Template, error) {
	return template.New("n").
		Funcs(ErrTemplateFuncs).
		Parse(strings.Trim("The {{ .Lang }} file \"{{ .File }}\" could not be found. {{ FuncTest \"Test\" }} Host: {{ .Url.Host }} Url: {{ .Url }}\n", " \n"))
}

func (e ErrFileNotFoundError) Stack() []byte {
	return e.stack
}

func (e ErrFileNotFoundError) Error() string {
	t, err := e.templ()
	if err != nil {
		panic(fmt.Sprintf("Error compiling template: %q", err))
	}
	buf := bytes.Buffer{}
	err = t.Execute(&buf, e.Params())
	if err != nil {
		panic(fmt.Sprintf("Error executing template: %q", err))
	}
	return buf.String()
}

// GetFile returns the value of the key
func (e ErrFileNotFoundError) GetFile() string {
	return e.params.File
}

// SetFile set the value of the key
func (e *ErrFileNotFoundError) SetFile(v string) {
	e.params.File = v
}

// File sets the value and returns a copy of the error (use for chaining)
func (e *ErrFileNotFoundError) File(v string) *ErrFileNotFoundError {
	e.params.File = v
	return e
}

// GetLang returns the value of the key
func (e ErrFileNotFoundError) GetLang() string {
	return e.params.Lang
}

// SetLang set the value of the key
func (e *ErrFileNotFoundError) SetLang(v string) {
	e.params.Lang = v
}

// Lang sets the value and returns a copy of the error (use for chaining)
func (e *ErrFileNotFoundError) Lang(v string) *ErrFileNotFoundError {
	e.params.Lang = v
	return e
}

// GetStatus returns the value of the key
func (e ErrFileNotFoundError) GetStatus() int {
	return e.params.Status
}

// SetStatus set the value of the key
func (e *ErrFileNotFoundError) SetStatus(v int) {
	e.params.Status = v
}

// Status sets the value and returns a copy of the error (use for chaining)
func (e *ErrFileNotFoundError) Status(v int) *ErrFileNotFoundError {
	e.params.Status = v
	return e
}

// GetUrl returns the value of the key
func (e ErrFileNotFoundError) GetUrl() *url.URL {
	return e.params.Url
}

// SetUrl set the value of the key
func (e *ErrFileNotFoundError) SetUrl(v *url.URL) {
	e.params.Url = v
}

// Url sets the value and returns a copy of the error (use for chaining)
func (e *ErrFileNotFoundError) Url(v *url.URL) *ErrFileNotFoundError {
	e.params.Url = v
	return e
}

// Params returns all parameters as map
func (e ErrFileNotFoundError) Params() interface{} {
	return e.params
}

// Wrap given error
func (e *ErrFileNotFoundError) Wrap(err error) *ErrFileNotFoundError {
	if e.parent != nil {
		panic("Unable to wrap ErrFileNotFound with already existing parent.")
	}
	e.parent = err
	return e
}

func (e *ErrFileNotFoundError) Unwrap() error {
	return e.parent
}
```
## Use the generated error

```go
err := ErrFileNotFound().
  File("/tmp/henk").
  Url(u)
```