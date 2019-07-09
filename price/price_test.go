package price

import (
	"testing"
)

func Test_ParsePrice(t *testing.T) {
	var db MapPriceDB
	db.priceMap = make(map[string]float64)
	db.priceMap["chaos"] = 1.0
	db.priceMap["exa"] = 100.0
	db.priceMap["alt"] = 0.01

	t.Run("fixed", func(t *testing.T) {
		cases := []struct {
			in       string
			expected ItemPrice
		}{
			{"-price 1 chaos", ItemPrice{Exact, "chaos", 1.0}},
			{"-price 1.5 chaos", ItemPrice{Exact, "chaos", 1.5}},
			{"-price 1.5 exa", ItemPrice{Exact, "exa", 150.0}},
			{"-price 25 alt", ItemPrice{Exact, "alt", 0.25}},
		}

		for _, tt := range cases {
			actual, err := ParsePrice(tt.in)
			if err != nil {
				t.Logf("unable to parse price: %s", err)
				t.Fail()
			}
			if tt.expected != actual {
				t.Logf("Expected `%#v`, received `%#v`", tt.expected, actual)
				t.Fail()
			}
		}
	})

	t.Run("negotiable", func(t *testing.T) {
		cases := []struct {
			in       string
			expected ItemPrice
		}{
			{"-price 1 chaos", ItemPrice{Negotiable, "chaos", 1.0}},
			{"-price 1.5 chaos", ItemPrice{Negotiable, "chaos", 1.5}},
			{"-price 1.5 exa", ItemPrice{Negotiable, "exa", 150.0}},
			{"-price 25 alt", ItemPrice{Negotiable, "alt", 0.25}},
		}

		for _, tt := range cases {
			actual, err := ParsePrice(tt.in)
			if err != nil {
				t.Logf("unable to parse price: %s", err)
				t.Fail()
			}
			if tt.expected != actual {
				t.Logf("Expected `%#v`, received `%#v`", tt.expected, actual)
				t.Fail()
			}
		}
	})

	t.Run("best offer", func(t *testing.T) {
		cases := []struct {
			in       string
			expected ItemPrice
		}{
			{"-price 1 chaos", ItemPrice{BetterOffer, "chaos", 1.0}},
			{"-price 1.5 chaos", ItemPrice{BetterOffer, "chaos", 1.5}},
			{"-price 1.5 exa", ItemPrice{BetterOffer, "exa", 150.0}},
			{"-price 25 alt", ItemPrice{BetterOffer, "alt", 0.25}},
		}

		for _, tt := range cases {
			actual, err := ParsePrice(tt.in)
			if err != nil {
				t.Logf("unable to parse price: %s", err)
				t.Fail()
			}
			if tt.expected != actual {
				t.Logf("Expected `%#v`, received `%#v`", tt.expected, actual)
				t.Fail()
			}
		}
	})
}
