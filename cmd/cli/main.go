package main

import (
	"os"
	"github.com/Lepovirta/keruu/feeds2html"
)

func main() {
	if err := feeds2html.FromStream(os.Stdin, os.Stdout); err != nil {
		panic(err)
	}
}
