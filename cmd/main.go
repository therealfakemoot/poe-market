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
		var (
			gear     metrics.HistogramSet
			gems     metrics.HistogramSet
			currency metrics.HistogramSet
			div      metrics.HistogramSet
			quest    metrics.HistogramSet
			prophecy metrics.HistogramSet
			relic    metrics.HistogramSet
		)

		gear.Histograms = make(map[poe.HistoKey]prometheus.Observer)
		gear.HistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "market",
			Name:      "gear_price_chaos",
		},
			[]string{
				"name",
				"sockets",
				"links",
			},
		)
		prometheus.MustRegister(gear.HistogramVec)

		gems.Histograms = make(map[poe.HistoKey]prometheus.Observer)
		gems.HistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "market",
			Name:      "gems_price_chaos",
		},
			[]string{
				"name",
			},
		)
		prometheus.MustRegister(gems.HistogramVec)

		currency.Histograms = make(map[poe.HistoKey]prometheus.Observer)
		currency.HistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "market",
			Name:      "currency_price_chaos",
		},
			[]string{
				"name",
			},
		)
		prometheus.MustRegister(currency.HistogramVec)

		div.Histograms = make(map[poe.HistoKey]prometheus.Observer)
		div.HistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "market",
			Name:      "div_price_chaos",
		},
			[]string{
				"name",
			},
		)
		prometheus.MustRegister(div.HistogramVec)

		quest.Histograms = make(map[poe.HistoKey]prometheus.Observer)
		quest.HistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "market",
			Name:      "quest_price_chaos",
		},
			[]string{
				"name",
			},
		)
		prometheus.MustRegister(quest.HistogramVec)

		prophecy.Histograms = make(map[poe.HistoKey]prometheus.Observer)
		prophecy.HistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "market",
			Name:      "prophecy_price_chaos",
		},
			[]string{
				"name",
			},
		)
		prometheus.MustRegister(prophecy.HistogramVec)

		relic.Histograms = make(map[poe.HistoKey]prometheus.Observer)
		relic.HistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "market",
			Name:      "relic_price_chaos",
		},
			[]string{
				"name",
			},
		)
		prometheus.MustRegister(relic.HistogramVec)

		for item := range stream.Items {
			if item.Note != "" {
				ip, err := pdb.Price(item)
				if err != nil {
					// log.Printf("bad parse: %s", item.Note)
					continue
				}

				switch item.FrameType {
				case 0, 1, 2, 3:
					_, ok := gear.Histograms[item.Key()]
					if !ok {
						gear.Add(item)
					}
					gear.Histograms[item.Key()].Observe(ip.Cost)
				case 4:
					_, ok := gems.Histograms[item.Key()]
					if !ok {
						gems.Add(item)
					}
					gems.Histograms[item.Key()].Observe(ip.Cost)
				case 5:
					_, ok := currency.Histograms[item.Key()]
					if !ok {
						currency.Add(item)
					}
					currency.Histograms[item.Key()].Observe(ip.Cost)
				case 6:
					_, ok := div.Histograms[item.Key()]
					if !ok {
						div.Add(item)
					}
					div.Histograms[item.Key()].Observe(ip.Cost)
				case 7:
					_, ok := quest.Histograms[item.Key()]
					if !ok {
						quest.Add(item)
					}
					quest.Histograms[item.Key()].Observe(ip.Cost)
				case 8:
					_, ok := prophecy.Histograms[item.Key()]
					if !ok {
						prophecy.Add(item)
					}
					prophecy.Histograms[item.Key()].Observe(ip.Cost)
				case 9:
					_, ok := relic.Histograms[item.Key()]
					if !ok {
						relic.Add(item)
					}
					relic.Histograms[item.Key()].Observe(ip.Cost)

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
