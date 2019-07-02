package pom

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// https://api.poe.watch/get?league=Legion&category=currency

var (
	ErrInvalidCurrencyQuantity = errors.New("unrecognized currency quantity")
	ErrUnrecognizedPriceType   = errors.New("unrecognized price type")
	ErrUnrecognizedCurrency    = errors.New("unrecognized currency type")
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

type PriceStatus int

const (
	Negotiable PriceStatus = iota
	Exact
	BetterOffer
)

type ItemPrice struct {
	PriceStatus PriceStatus
	Cost        float64
}

type PriceDB interface {
	Convert(string, string) (float64, error)
}

func ParsePrice(s string, db PriceDB) (ItemPrice, error) {
	var ip ItemPrice

	fields := strings.Split(s, " ")
	switch fields[0] {
	case "-price":
		ip.PriceStatus = Exact
	case "~price":
		ip.PriceStatus = Negotiable
	case "~b/o":
		ip.PriceStatus = BetterOffer
	default:
		return ip, ErrUnrecognizedPriceType
	}

	cost, err := db.Convert(fields[1], fields[2])
	if err != nil {
		return ip, err
	}
	ip.Cost = cost

	return ip, nil
}
