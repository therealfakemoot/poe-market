package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"

	"github.com/therealfakemoot/pom/metrics"
	"github.com/therealfakemoot/pom/poe"
)

func LoadFile(filename string) (poe.Envelope, error) {
	var e poe.Envelope

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

	APILimit := rate.Limit(2.0)
	l := rate.NewLimiter(APILimit, 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	stream := poe.New(l)

	log.Println("starting stream")
	go stream.Start(ctx)

	var gs metrics.GaugeSet
	gs.Gauges = make(map[string]prometheus.Gauge)

	go func() {
		for stash := range stream.Stashes {
			if stash.Public {
				for _, item := range stash.Items {
					gs.RegisterItem(item)
				}
			}
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe("0.0.0.0:9092", nil)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}
