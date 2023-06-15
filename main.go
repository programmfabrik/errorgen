package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	conf := Config{}
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "input",
				Aliases:     []string{"i"},
				DefaultText: "-",
				Usage:       "Load .yml input from `FILE`. Use \"-\" to read from stdin",
				Destination: &conf.InputFile,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				DefaultText: "-",
				Usage:       "Write generated code to `FILE`. Use \"-\" to output to stdout",
				Destination: &conf.OutputFile,
			},
			&cli.StringFlag{
				Name:        "tmpl",
				Aliases:     []string{"t"},
				Usage:       "Use `FILE` as code template for generating",
				Destination: &conf.Template,
			},
			&cli.StringFlag{
				Name:        "csv",
				Aliases:     []string{"c"},
				Usage:       "Append to `FILE` in CSV format: key, parameters, output",
				Destination: &conf.CsvFile,
			},
			&cli.BoolFlag{
				Name:        "html",
				Usage:       "Use html/template instead of text/template for rendering",
				Destination: &conf.HTML,
			},
			&cli.IntFlag{
				Name:        "server",
				Usage:       "Start webserver on port to show html, implies --html",
				Destination: &conf.Server,
			},
		},
		Name:  "errors",
		Usage: "generate Go error definitions from .yml input",
		Action: func(c *cli.Context) (err error) {
			return runGen(c, conf)
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
