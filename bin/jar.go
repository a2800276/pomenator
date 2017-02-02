package main

import "pomenator"

import (
	"flag"
	"fmt"
	"os"
)

var fs = []string{}
var jarFn = flag.String("jar", "", "name of jar file to produce")
var dir = flag.String("dir", "", "directory to pack into jar")

func main() {

	flag.Parse()
	if *jarFn == "" || *dir == "" {
		usage("")
	}
	err := pomenator.GenerateJarFromDirs(*jarFn, *dir)
	if err != nil {
		usage(fmt.Sprintf("%v\n", err))
	}
}

func usage(m string) {
	fmt.Fprintf(os.Stderr, m)
	flag.Usage()
	os.Exit(1)
}
