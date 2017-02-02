package main

import "pomenator"

var fs = []string{}

func main() {
	err := pomenator.GenerateJarFromFiles("test.jar", ".", "jar.go")
	println(err)
}
