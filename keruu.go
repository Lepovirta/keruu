package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"gitlab.com/lepovirta/keruu/feeds2html"
)

var configPath string
var outputPath string
var config feeds2html.Config

func init() {
	// Setup CLI flags
	flag.StringVar(&configPath, "config", "STDIN", "Location of the configuration file")
	flag.StringVar(&outputPath, "output", "STDOUT", "Location of the HTML output file")

	// Init config
	config.Init()

	// Setup custom usage function
	defaultUsage := flag.Usage
	flag.Usage = func() {
		defaultUsage()
		fmt.Fprintf(flag.CommandLine.Output(), "\nConfiguration format:\n\n%s", feeds2html.ConfigTemplateYAML)
	}
}

func main() {
	flag.Parse()

	if err := readConfig(); err != nil {
		log.Fatalf("failed to read config from STDIN: %s", err)
	}

	if err := writeOutput(func(w io.Writer) error {
		return feeds2html.Run(&config, w)
	}); err != nil {
		log.Panicf("feed aggregation failed: %s", err)
	}
}

func readConfig() error {
	if isSTDIN() {
		return config.FromSTDIN()
	}
	return config.FromYAMLFile(configPath)
}

func writeOutput(f func(io.Writer) error) error {
	var writer io.Writer
	if isSTDOUT() {
		writer = os.Stdout
	} else {
		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("failed to close output file: %s", err)
			}
		}()
		writer = file
	}

	bufWriter := bufio.NewWriter(writer)
	if err := f(bufWriter); err != nil {
		return err
	}
	return bufWriter.Flush()
}

func isSTDIN() bool {
	switch configPath {
	case "", "-", "STDIN":
		return true
	}
	return false
}

func isSTDOUT() bool {
	switch outputPath {
	case "", "-", "STDOUT":
		return true
	}
	return false
}
