package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"gitlab.com/lepovirta/keruu/internal/aggregation"
	"gitlab.com/lepovirta/keruu/internal/config"
	"gitlab.com/lepovirta/keruu/internal/fetch"
	"gitlab.com/lepovirta/keruu/internal/file"
)

var cfgPath string
var outPath string
var cfg config.Config

func init() {
	// Setup CLI flags
	flag.StringVar(&cfgPath, "config", "STDIN", "Location of the configuration file")
	flag.StringVar(&outPath, "output", "STDOUT", "Location of the HTML output file")

	// Init config
	cfg.Init()

	// Setup custom usage function
	defaultUsage := flag.Usage
	flag.Usage = func() {
		defaultUsage()
		fmt.Fprintf(flag.CommandLine.Output(), "\nConfiguration format:\n\n%s", config.TemplateYAML)
	}
}

func main() {
	flag.Parse()

	if err := readConfig(); err != nil {
		log.Fatalf("failed to read config from STDIN: %s", err)
	}

	if err := writeOutput(func(w io.Writer) error {
		// Error checking is intentionally skipped here to report it later
		posts, err := fetch.Run(&cfg.Fetch, cfg.Feeds, cfg.Links)

		aggregation := aggregation.New(&cfg.Aggregation, posts)
		if err := aggregation.ToHTML(w); err != nil {
			return err
		}
		return err
	}); err != nil {
		log.Panicf("feed aggregation failed: %s", err)
	}
}

func readConfig() error {
	if isSTDIN() {
		return cfg.FromSTDIN()
	}
	return cfg.FromYAMLFile(cfgPath)
}

func writeOutput(f func(io.Writer) error) error {
	if isSTDOUT() {
		writer := bufio.NewWriter(os.Stdout)
		if err := f(writer); err != nil {
			return err
		}
		return writer.Flush()
	} else {
		return file.WithFileWriter(outPath, f)
	}
}

func isSTDIN() bool {
	switch cfgPath {
	case "", "-", "STDIN":
		return true
	}
	return false
}

func isSTDOUT() bool {
	switch outPath {
	case "", "-", "STDOUT":
		return true
	}
	return false
}
