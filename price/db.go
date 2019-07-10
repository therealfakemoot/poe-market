package price

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/therealfakemoot/pom/poe"
)

func NewLiveDB() (LiveDB, error) {
	var prices []PricePoint

	ldb := make(LiveDB)

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
		return ldb, err
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	err = d.Decode(&prices)
	if err != nil {
		return ldb, err
	}

	priceMap := make(map[int]float64)

	for _, pp := range prices {
		priceMap[pp.ID] = pp.Median
	}

	for k, v := range IDMap {
		ldb[k] = float64(priceMap[v])
	}

	ldb["chaos"] = 1.0
	return ldb, nil
}

type LiveDB map[string]float64

func (ldb LiveDB) Price(item poe.Item) (ItemPrice, error) {
	var ip ItemPrice

	n := item.Note
	fields := strings.Fields(n)

	if n == "~price" {
		ip.PriceStatus = Negotiable
		return ip, nil
	}

	if len(fields) < 3 {
		var err ErrBadParse
		err.raw = n
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

	ip.BaseCurrency = fields[2]
	marketPrice, ok := ldb[ip.BaseCurrency]
	if !ok {
		return ip, ErrUnrecognizedCurrency
	}
	itemCost, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return ip, err
	}

	ip.Cost = itemCost * marketPrice

	return ip, nil
}
