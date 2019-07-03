package main

import (
	"fmt"

	"github.com/therealfakemoot/pom/price"
)

func main() {
	prices, err := price.PriceCheck()
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	fmt.Printf("%#v", prices[0])
}
