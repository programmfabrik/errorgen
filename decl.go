package main

import (
	"errors"
	"fmt"
)

type Declaration struct {
	Package  string
	Import   []string
	Prefix   string // Defaults to "", could be "Err" if errors are used as part of another package
	Defaults map[string]ErrorParams
	Errors   map[string]ErrorDef
}

type declInterfaces map[string]map[string]string

type ErrorDef struct {
	D      string      // (D)escription
	O      string      // (O)utput
	P      ErrorParams // (P)arams
	Params ErrorParams // merged params
}

type ErrorParams map[string]ErrorParam

type ErrorParam struct {
	T string      // Go (T)ype of the param
	D string      // (D)escription of the param
	V interface{} // Default (V)alue of the param
}

// merge sets key in eps to ep, merges info if needed
func (eps *ErrorParams) merge(key string, ep ErrorParam) error {
	ep1, ok := (*eps)[key]
	if !ok {
		(*eps)[key] = ep
		return nil
	}
	// key is already in the map, merge info
	if ep.T != "" {
		if ep1.T == "" {
			ep1.T = ep.T
		} else if ep1.T != ep.T {
			return errors.New(fmt.Sprintf("Param %q has conflicting type %q. Default type: %q", key, ep.T, ep1.T))
		}
	}

	if ep.D != "" {
		ep1.D = ep.D
	}

	if ep.V != nil {
		ep1.V = ep.V
	}
	(*eps)[key] = ep1
	return nil
}

// merge lets ep2 params overwrite ep params and returns a new merged map
func mergeParams(defs map[string]ErrorParams, ep1 ErrorParams) (ep ErrorParams, err error) {
	ep = ErrorParams{}
	for _, errParam := range defs {
		for pname, param := range errParam {
			err = ep.merge(pname, param)
			if err != nil {
				return ep, err
			}
		}
	}
	for pname, param := range ep1 {
		err = ep.merge(pname, param)
		if err != nil {
			return ep, err
		}
	}
	return ep, nil
}

type reservedNames []string

func isReserved(name string) bool {
	for _, n := range []string{"Error", "Wrap", "Unwrap", "Data", "Name", "Stack"} {
		if n == name {
			return true
		}
	}
	return false
}

// Validate returns an error if the declaration is not valid
func (d *Declaration) Validate() (err error) {

	// Check that all defaults have a type, without we cannot create
	// the interfaces
	for defName, errParams := range d.Defaults {
		for pname, errParam := range errParams {
			if isReserved(pname) {
				return errors.New(fmt.Sprintf(`Error %q: Parameter cannot be reversed name %q`, defName, pname))
			}
			if errParam.T == "" {
				return errors.New(fmt.Sprintf(`Default %q.%q requires a type.`, defName, pname))
			}
		}
	}

	for ename, errDef := range d.Errors {
		for pname := range errDef.P {
			if isReserved(pname) {
				return errors.New(fmt.Sprintf(`Error %q: Parameter cannot be reversed name %q`, ename, pname))
			}
		}
		errDef.Params, err = mergeParams(d.Defaults, errDef.P)
		if err != nil {
			return err
		}
		d.Errors[ename] = errDef
	}
	return nil
}
