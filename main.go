package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

func LoadFile(filename string) (Envelope, error) {
	var e Envelope

	f, err := os.Open(filename)
	defer f.Close()

	if err != nil {
		return e, err
	}

	d := json.NewDecoder(f)
	err = d.Decode(&e)

	if err != nil {
		return e, err
	}

	log.Printf("Next page ID: %s", e.NextChangeID)

	return e, nil

}

func main() {
	var (
		input string
		live  bool
	)
	flag.StringVar(&input, "file", "tabs.json", "input file")
	flag.BoolVar(&live, "api", false, "use live api data")

	flag.Parse()
	var (
		e   Envelope
		err error
	)

	if !live {
		e, err = LoadFile(input)
		if err != nil {
			log.Fatalf("error loading file: %s", err)
		}
		log.Printf("Loading %s as datasource\n", input)
	}

	e, err = FetchStashes("")
	if err != nil {
		log.Fatalf("error fetching API data: %s", err)
	}

	log.Printf("next-change-id: %s", e.NextChangeID)
	log.Printf("stash count: %d", len(e.Stashes))
}
