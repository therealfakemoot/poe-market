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

	pdb, err := price.NewLiveDB()
	if err != nil {
		log.Printf("error creating pricedb: %s", err)
		return
	}

	go func() {
		var gs metrics.HistogramSet
		gs.HistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "market",
			Name:      "item_price_chaos",
		},
			[]string{
				"name",
				"sockets",
				"links",
				"frametype",
			},
		)
		prometheus.MustRegister(gs.HistogramVec)

		gs.Histograms = make(map[poe.HistoKey]prometheus.Observer)
		for item := range stream.Items {

			_, ok := gs.Histograms[item.Key()]
			if !ok {
				gs.Add(item)
			}
			if item.Note != "" {

				ip, err := pdb.Price(item)
				if err != nil {
					// log.Printf("bad parse: %s", item.Note)
					continue
				}

				gs.Histograms[item.Key()].Observe(ip.Cost)
			}
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe("0.0.0.0:9092", nil)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}
