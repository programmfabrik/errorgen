package main

//go:generate ../errorgen -i simple.yml -o example.go

import (
	"encoding/json"
	"errors"
	"net/url"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestFileNotFound(t *testing.T) {

	ErrTemplateFuncs = template.FuncMap{
		"FuncTest": func(s string) string {
			return "FuncTest: " + s
		},
	}

	u, _ := url.Parse("http://slashdot.org")

	pErr := errors.New("Testerror")

	err := ErrFileNotFound().
		File("/tmp/henk").
		Url(u).Wrap(pErr)

	var e1 error = err

	if !assert.Equal(t, true, errors.Is(err, e1)) {
		return
	}

	if !assert.Equal(t, true, errors.Is(err, ErrFileNotFound())) {
		return
	}

	if !assert.Equal(t, true, errors.Is(err, pErr)) {
		return
	}

	var errHTTP ErrHTTPError

	if !assert.Equal(t, true, errors.As(e1, &errHTTP)) {
		return
	}

	errHTTP.SetLang("Deutsch")

	if !assert.Equal(t, `The Deutsch file "/tmp/henk" could not be found. FuncTest: Test Host: slashdot.org Url: http://slashdot.org`, err.Error()) {
		return
	}

	if !assert.Equal(t, 400, errHTTP.GetStatus()) {
		return
	}

	if !assert.Equal(t, "FileNotFound", errHTTP.ErrorCode()) {
		return
	}

	// pp.Println(errHTTP.Data())

	if !assert.Equal(t, "slashdot.org", errHTTP.GetUrl().Host) {
		return
	}

	json, _ := json.Marshal(errHTTP.Params())
	if !assert.Equal(t, `{"file":"/tmp/henk","lang":"Deutsch","s":400,"url":{"Scheme":"http","Opaque":"","User":null,"Host":"slashdot.org","Path":"","RawPath":"","ForceQuery":false,"RawQuery":"","Fragment":"","RawFragment":""}}`, string(json)) {
		return
	}

	// println(string(errHTTP.Stack()))

}
