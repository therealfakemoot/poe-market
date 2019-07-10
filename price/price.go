package price

import (
	"errors"
	"fmt"
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
