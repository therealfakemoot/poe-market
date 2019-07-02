package price

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func PriceCheck() ([]PricePoint, error) {
	var prices []PricePoint

	queryValues := url.Values{}
	queryValues.Set("league", "Legion")
	queryValues.Set("category", "currency")

	PricesEndpoint := url.URL{
		Host:     "api.poe.watch",
		Scheme:   "https",
		Path:     "get",
		RawQuery: queryValues.Encode(),
	}

	c := http.Client{}

	resp, err := c.Get(PricesEndpoint.String())
	if err != nil {
		return prices, err
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	err = d.Decode(&prices)
	if err != nil {
		return prices, err
	}

	return prices, nil
}

type PricePoint struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Group     string    `json:"group"`
	Frame     int       `json:"frame"`
	StackSize int       `json:"stackSize,omitempty"`
	Icon      string    `json:"icon"`
	Mean      float64   `json:"mean"`
	Median    float64   `json:"median"`
	Mode      float64   `json:"mode"`
	Min       float64   `json:"min"`
	Max       float64   `json:"max"`
	Exalted   float64   `json:"exalted"`
	Total     int       `json:"total"`
	Daily     int       `json:"daily"`
	Current   int       `json:"current"`
	Accepted  int       `json:"accepted"`
	Change    float64   `json:"change"`
	History   []float64 `json:"history"`
	Type      string    `json:"type,omitempty"`
}
