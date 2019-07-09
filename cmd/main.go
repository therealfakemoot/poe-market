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
	"github.com/therealfakemoot/pom/price"
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

	go func() {
		log.Fatalf("error on errchan: %s", <-stream.Err)
	}()

	go func() {
		var gs metrics.GaugeSet
		gs.GaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "market",
			Name:      "price_chaos",
		},
			[]string{
				"name",
				"sockets",
				"links",
				"frametype",
			})
		prometheus.MustRegister(gs.GaugeVec)

		gs.Gauges = make(map[poe.GaugeKey]prometheus.Gauge)
		for item := range stream.Items {

			_, ok := gs.Gauges[item.Key()]
			if !ok {
				gs.Add(item)
			}
			if item.Note != "" {

				ip, err := price.ParsePrice(item.Note)
				if err != nil {
					log.Printf("bad parse: %s", item.Note)
					continue
				}

				gs.Gauges[item.Key()].Set(ip.Cost)
			}
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe("0.0.0.0:9092", nil)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}
