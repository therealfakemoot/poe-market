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
			gear     metrics.SummarySet
			gems     metrics.SummarySet
			currency metrics.SummarySet
			div      metrics.SummarySet
			quest    metrics.SummarySet
			prophecy metrics.SummarySet
			relic    metrics.SummarySet
		)

		gear.Summaries = make(map[poe.SummaryKey]prometheus.Observer)
		gear.SummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: "market",
			Name:      "gear_price_chaos",
			Objectives: map[float64]float64{
				0.5:  0.1,
				0.9:  0.1,
				0.99: 0.1,
			},
		},
			[]string{
				"name",
				"sockets",
				"links",
			},
		)
		prometheus.MustRegister(gear.SummaryVec)

		gems.Summaries = make(map[poe.SummaryKey]prometheus.Observer)
		gems.SummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: "market",
			Name:      "gems_price_chaos",
			Objectives: map[float64]float64{
				0.5:  0.1,
				0.9:  0.1,
				0.99: 0.1,
			},
		},
			[]string{
				"name",
			},
		)
		prometheus.MustRegister(gems.SummaryVec)

		currency.Summaries = make(map[poe.SummaryKey]prometheus.Observer)
		currency.SummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: "market",
			Name:      "currency_price_chaos",
			Objectives: map[float64]float64{
				0.5:  0.1,
				0.9:  0.1,
				0.99: 0.1,
			},
		},
			[]string{
				"name",
			},
		)
		prometheus.MustRegister(currency.SummaryVec)

		div.Summaries = make(map[poe.SummaryKey]prometheus.Observer)
		div.SummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: "market",
			Name:      "div_price_chaos",
			Objectives: map[float64]float64{
				0.5:  0.1,
				0.9:  0.1,
				0.99: 0.1,
			},
		},
			[]string{
				"name",
			},
		)
		prometheus.MustRegister(div.SummaryVec)

		quest.Summaries = make(map[poe.SummaryKey]prometheus.Observer)
		quest.SummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: "market",
			Name:      "quest_price_chaos",
			Objectives: map[float64]float64{
				0.5:  0.1,
				0.9:  0.1,
				0.99: 0.1,
			},
		},
			[]string{
				"name",
			},
		)
		prometheus.MustRegister(quest.SummaryVec)

		prophecy.Summaries = make(map[poe.SummaryKey]prometheus.Observer)
		prophecy.SummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: "market",
			Name:      "prophecy_price_chaos",
			Objectives: map[float64]float64{
				0.5:  0.1,
				0.9:  0.1,
				0.99: 0.1,
			},
		},
			[]string{
				"name",
			},
		)
		prometheus.MustRegister(prophecy.SummaryVec)

		relic.Summaries = make(map[poe.SummaryKey]prometheus.Observer)
		relic.SummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: "market",
			Name:      "relic_price_chaos",
			Objectives: map[float64]float64{
				0.5:  0.1,
				0.9:  0.1,
				0.99: 0.1,
			},
		},
			[]string{
				"name",
			},
		)
		prometheus.MustRegister(relic.SummaryVec)

		for item := range stream.Items {
			if item.Note != "" {
				ip, err := pdb.Price(item)
				if err != nil {
					// log.Printf("bad parse: %s", item.Note)
					continue
				}

				switch item.FrameType {
				case 0, 1, 2, 3:
					_, ok := gear.Summaries[item.Key()]
					if !ok {
						gear.Add(item)
					}
					gear.Summaries[item.Key()].Observe(ip.Cost)
				case 4:
					_, ok := gems.Summaries[item.Key()]
					if !ok {
						gems.Add(item)
					}
					gems.Summaries[item.Key()].Observe(ip.Cost)
				case 5:
					_, ok := currency.Summaries[item.Key()]
					if !ok {
						currency.Add(item)
					}
					currency.Summaries[item.Key()].Observe(ip.Cost)
				case 6:
					_, ok := div.Summaries[item.Key()]
					if !ok {
						div.Add(item)
					}
					div.Summaries[item.Key()].Observe(ip.Cost)
				case 7:
					_, ok := quest.Summaries[item.Key()]
					if !ok {
						quest.Add(item)
					}
					quest.Summaries[item.Key()].Observe(ip.Cost)
				case 8:
					_, ok := prophecy.Summaries[item.Key()]
					if !ok {
						prophecy.Add(item)
					}
					prophecy.Summaries[item.Key()].Observe(ip.Cost)
				case 9:
					_, ok := relic.Summaries[item.Key()]
					if !ok {
						relic.Add(item)
					}
					relic.Summaries[item.Key()].Observe(ip.Cost)

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
