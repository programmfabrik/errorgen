package main

// Config contains all information needed to run errors.
type Config struct {
	InputFile  string
	OutputFile string

	Template string // Template file for rendering
	CsvFile  string // Append to CSV file
	HTML     bool   // True, if rendering should use template/html instead of template/text for rendering
	Server   int    // Start server on port

	D Declaration // Read and parsed from InputFile

	inputDebug string // set to stdin or filename for debug output
}

func (c *Config) Validate() error {
	err := c.D.Validate()
	if err != nil {
		return err
	}
	if c.Server > 0 {
		c.HTML = true
	}
	return nil
}
