package price

import (
	"errors"
	"strings"
)

// https://api.poe.watch/get?league=Legion&category=currency

var (
	ErrInvalidCurrencyQuantity = errors.New("unrecognized currency quantity")
	ErrUnrecognizedPriceType   = errors.New("unrecognized price type")
	ErrUnrecognizedCurrency    = errors.New("unrecognized currency type")
)

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
