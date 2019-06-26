package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"net/url"
)

type StreamError struct {
	PageID string
	Err    error
}

func (se StreamError) Error() string {
	return fmt.Sprintf(`unable to fetch page "%s": %s`, se.PageID, se.Err)
}

func New(l *rate.Limiter) StashStream {
	var sa StashStream
	sa.Limiter = l
	sa.Stashes = make(chan Stash, 5)
	sa.Err = make(chan error)
	return sa
}

type StashStream struct {
	Limiter *rate.Limiter
	NextID  string
	Stashes chan Stash
	Err     chan error
}

func (sa *StashStream) Start(ctx context.Context) error {
	sa.Limiter.Wait(ctx)

	BASE := url.URL{
		Scheme: "http",
		Host:   "www.pathofexile.com",
		Path:   "/api/public-stash-tabs",
	}

	c := &http.Client{}

	var e Envelope

	for {
		log.Printf(`requesting page "%s"`, sa.NextID)
		v := url.Values{}
		v.Set("id", sa.NextID)
		BASE.RawQuery = v.Encode()

		resp, err := c.Get(BASE.String())

		if err != nil {
			sa.Err <- StreamError{
				PageID: sa.NextID,
				Err:    err,
			}
			log.Printf("error requesting page: %s", err)
			continue
		}
		defer resp.Body.Close()

		d := json.NewDecoder(resp.Body)
		err = d.Decode(&e)
		if err != nil {
			sa.Err <- StreamError{
				PageID: sa.NextID,
				Err:    err,
			}
			log.Printf("error decoding envelope: %s", err)
			continue
		}
		log.Printf("next page ID: %s", e.NextChangeID)
		sa.NextID = e.NextChangeID

		for _, stash := range e.Stashes {
			sa.Stashes <- stash
		}

		select {
		case <-ctx.Done():
			return err
		}

	}
}
