package price

import (
	"errors"
	"fmt"
	"strings"
)

// https://api.poe.watch/get?league=Legion&category=currency

var (
	ErrInvalidCurrencyQuantity = errors.New("unrecognized currency quantity")
	ErrUnrecognizedPriceType   = errors.New("unrecognized price type")
	ErrUnrecognizedCurrency    = errors.New("unrecognized currency type")
)

type ErrBadParse struct {
	raw string
}

func (ebp ErrBadParse) Error() string {
	return fmt.Sprintf("unable to parse %s", ebp.raw)
}

type PriceStatus int

const (
	Negotiable PriceStatus = iota
	Exact
	BetterOffer
)

type ItemPrice struct {
	PriceStatus  PriceStatus
	BaseCurrency string
	Cost         float64
}

func ParsePrice(s string) (ItemPrice, error) {
	var ip ItemPrice

	fields := strings.Split(s, " ")

	if s == "~price" {
		ip.PriceStatus = Negotiable
		return ip, nil
	}

	if len(fields) < 3 {
		var err ErrBadParse
		err.raw = s
		return ip, err
	}

	switch fields[0] {
	case "-price":
		ip.PriceStatus = Exact
	case "~price":
		ip.PriceStatus = Negotiable
	case "~b/o", "b/o":
		ip.PriceStatus = BetterOffer
	default:
		return ip, ErrUnrecognizedPriceType
	}

	/*
		cost, err := db.Convert(fields[1], fields[2])
		if err != nil {
			return ip, err
		}
		ip.Cost = cost
	*/
	ip.BaseCurrency = fields[2]

	return ip, nil
}
