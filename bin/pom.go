package main

import "pomenator"
import (
	"fmt"
)

var fs = []string{}

func main() {
	cfg, err := pomenator.LoadConfig("test/bootstrap.json")
	if err != nil {
		println(err.Error())
	}
	fmt.Printf("%v\n", cfg)

	cfg, err = pomenator.LoadConfig("test/bootstrap2.json")
	if err != nil {
		println(err.Error())
	}
	fmt.Printf("%v\n", cfg)

	err = pomenator.GeneratePOM("test/test.delete.pom", cfg[0])
	if err != nil {
		println(err.Error())
	}

}
