package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"golang.org/x/time/rate"
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
		file string
	)
	flag.StringVar(&file, "file", "stash-data.db", "database file")

	flag.Parse()
	var (
		err error
	)

	APILimit := rate.Limit(1.0)
	l := rate.NewLimiter(APILimit, 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	stream := New(l)

	log.Println("starting stream")
	go stream.Start(ctx)
	var count int
	for _ = range stream.Stashes {
		count++
		log.Printf("%d total stashes received", count)
	}
	if err != nil {
		log.Fatalf("error fetching API data: %s", err)
	}
}
