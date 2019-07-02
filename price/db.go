package price

import (
	"strconv"
)

type PriceDB interface {
	Convert(string, string) (float64, error)
}

type MapPriceDB struct {
	priceMap map[string]float64
}

func (db MapPriceDB) Convert(q string, t string) (float64, error) {

	c, ok := db.priceMap[t]
	if !ok {
		return 0.0, ErrUnrecognizedCurrency
	}

	v, err := strconv.ParseFloat(q, 64)
	if err != nil {
		return 0.0, ErrInvalidCurrencyQuantity
	}

	return c * v, nil
}
