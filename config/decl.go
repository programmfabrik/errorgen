package config

import "strings"

type Declaration struct {
	Errors map[string]Error
}

type Error struct {
	D string           // Description
	O string           // Output
	P map[string]Param // Params
}

func (e Error) Oescaped() string {
	return "`" + strings.ReplaceAll(e.O, "`", "``") + "`"
}

type Param struct {
	T string // Go type of the param
	D string // Description of the param
}
