package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var (
		input string
		live  bool
	)
	flag.StringVar(&input, "file", "tabs.json", "input file")
	flag.BoolVar(&live, "api", false, "use live api data")

	flag.Parse()

	if !live {
		f, err := os.Open(input)
		defer f.Close()

		if err != nil {
			log.Fatalf(`Unable to open "%s": %s`, input, err)
		}

		var e Envelope

		d := json.NewDecoder(f)
		err = d.Decode(&e)

		if err != nil {
			log.Fatalf(`Unable to decode payload: %s`, err)
		}

		log.Printf("Next page ID: %s", e.NextChangeID)

		return
	}

	fmt.Printf("Loading %s as datasource\n", input)
}
