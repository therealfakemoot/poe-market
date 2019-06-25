package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func FetchStashes(id string) (Envelope, error) {
	var e Envelope

	BASE := url.URL{
		Scheme: "http",
		Host:   "www.pathofexile.com",
		Path:   "/api/public-stash-tabs",
	}

	c := &http.Client{}

	v := url.Values{}
	if id != "" {
		v.Set("id", id)
	}
	BASE.RawQuery = v.Encode()

	resp, err := c.Get(BASE.String())

	if err != nil {
		return e, err
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	err = d.Decode(&e)
	if err != nil {
		return e, err
	}

	return e, nil
}
