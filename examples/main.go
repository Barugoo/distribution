package main

import (
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/barugoo/distribution/examples/products"
	"github.com/barugoo/distribution/utils"
)

func main() {
	ps := []products.Product{
		{
			Price:   utils.ToDecimal(10),
			VatRate: 5,
		},
		{
			Price:   utils.ToDecimal(15),
			VatRate: 5,
		},
		{
			Price:   utils.ToDecimal(25),
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
