package main

import (
	"net/url"
	"testing"
	"text/template"

	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
)

func TestFileNotFound(t *testing.T) {

	ErrTemplateFuncs = template.FuncMap{
		"FuncTest": func(s string) string {
			return "FuncTest: " + s
		},
	}

	u, _ := url.Parse("http://slashdot.org")

	err := ErrFileNotFound().
		File("/tmp/henk").
		Url(u)

	var e interface{} = &err
	errHTTP := e.(ErrHTTPError)
	errHTTP.SetLang("Deutsch")

	if !assert.Equal(t, `The Deutsch file "/tmp/henk" could not be found. FuncTest: Test Host: slashdot.org`, err.Error()) {
		return
	}

	if !assert.Equal(t, 400, errHTTP.GetStatus()) {
		return
	}

	if !assert.Equal(t, "FileNotFound", errHTTP.Name()) {
		return
	}

	pp.Println(errHTTP.Data())

	if !assert.Equal(t, "slashdot.org", errHTTP.GetUrl().Host) {
		return
	}

	// println(string(errHTTP.Stack()))

}
