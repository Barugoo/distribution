// Example of how to use the distribution core package
package main

import (
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/barugoo/distribution/examples/products"
)

func main() {
	ps := []products.Product{
		{
			Price:   decimal.NewFromFloat(10),
			VatRate: 5,
		},
		{
			Price:   decimal.NewFromFloat(15),
			VatRate: 5,
		},
		{
			Price:   decimal.NewFromFloat(25),
			VatRate: 7,
		},
	}

	productsLayout, err := products.MakeLayout(ps)
	if err != nil {
		panic(err)
	}

	someDecimalValue := decimal.NewFromFloat(55.55)
	v := productsLayout.DistributeDecimal(someDecimalValue)

	fmt.Printf("someDecimalValue distributed according to products VAT groups: \n\t%s\n", v)
	// Output:
	// someDecimalValue distributed according to products VAT groups:
	// 		[5: 27.78, 7: 27.77]
}
