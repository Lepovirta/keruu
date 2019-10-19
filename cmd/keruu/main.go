package main

import (
	"github.com/Lepovirta/keruu/feeds2html"
	"log"
	"os"
)

func main() {
	conf, err := feeds2html.ConfigFromYAML(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read config from STDIN: %s", err)
	}
	if err := feeds2html.Run(conf, os.Stdout); err != nil {
		log.Panicf("feed aggregation failed: %s", err)
	}
}
