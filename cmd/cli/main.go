package main

import (
	"github.com/Lepovirta/keruu/feeds2html"
	"os"
)

func main() {
	conf := feeds2html.DefaultConfig()
	if err := feeds2html.New(conf).FromStream(os.Stdin, os.Stdout); err != nil {
		panic(err)
	}
}
