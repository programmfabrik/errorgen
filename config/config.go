package config

// Config contains all information needed to run errors.
type Config struct {
	InputFile  string
	OutputFile string

	Prefix  string // Defaults to "", could be "Err" if errors are used as part of another package
	Package string

	D Declaration // Read and parsed from InputFile
}
