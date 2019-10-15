package main

import (
	"github.com/Lepovirta/keruu/feeds2html"
	"os"
	"bufio"
	"log"
)

func main() {
	feedUrls := make([]string, 0, 100)
	feedScanner := bufio.NewScanner(os.Stdin)
	for feedScanner.Scan() {
		feedUrls = append(feedUrls, feedScanner.Text())
	}

	if err := feedScanner.Err(); err != nil {
		log.Printf("error reading feed links: %s", err)
	}

	conf := feeds2html.DefaultConfig()
	conf.FeedURLs = feedUrls

	if err := feeds2html.Run(conf, os.Stdout); err != nil {
		panic(err)
	}
}
