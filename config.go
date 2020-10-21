package main

// Config contains all information needed to run errors.
type Config struct {
	InputFile  string
	OutputFile string

	Template string

	D Declaration // Read and parsed from InputFile
}
