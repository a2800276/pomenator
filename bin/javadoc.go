package main

import "pomenator"
import (
	"flag"
	"os"
)

var src = flag.String("source", "", "source directory")
var out = flag.String("out", "", "output directory")

func main() {

	flag.Parse()
	if *src == "" || *out == "" {
		usage()
	}
	if err := pomenator.GenerateJavadoc([]string{*src}, *out); err != nil {
		println(err.Error())
		usage()
	}
}

func usage() {

	flag.Usage()
	os.Exit(1)
}
