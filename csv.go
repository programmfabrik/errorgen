package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (d Declaration) writeCSV(fn string) (err error) {

	fn, err = filepath.Abs(fn)
	if err != nil {
		return err
	}

	errs := []string{}
	for key := range d.Errors {
		errs = append(errs, key)
	}
	recs := [][]string{}
	sort.Strings(errs)
	for _, err := range errs {
		def := d.Errors[err]
		params := []string{}
		for param := range def.Params {
			params = append(params, param)
		}
		sort.Strings(params)
		recs = append(recs, []string{
			d.Package,                  // package
			err,                        // key
			strings.Join(params, ", "), // params
			strings.TrimSpace(def.O),   // output
		})
	}

	writeHeader := false

	_, err = os.Stat(fn)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeHeader = true
		} else {
			return err
		}
	}

	f, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = '\u0009' // tab

	if writeHeader {
		err = w.Write([]string{
			"package", "key", "params", "error",
		})
		if err != nil {
			return
		}
	}

	err = w.WriteAll(recs)
	if err != nil {
		return
	}
	w.Flush()

	fmt.Fprintf(os.Stderr, "[errorgen] wrote csv to %q: %d keys\n", fn, len(recs))

	return nil

}
