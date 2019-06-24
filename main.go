package main

import (
	// "encoding/json"
	"flag"
	"fmt"
)

func main() {
	var (
		input string
		live  bool
	)
	flag.StringVar(&input, "file", "tabs.json", "input file")
	flag.BoolVar(&live, "api", false, "use ")

	if !live {
		return
	}

	fmt.Printf("Loading %s as datasource\n", input)
}
