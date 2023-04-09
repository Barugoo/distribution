package products

import (
	"fmt"
	"sort"

	"github.com/shopspring/decimal"

	"github.com/barugoo/distribution/core"
)

const distributionPrecision = 2

type Layout struct {
	*core.Layout[float64]
	productsTotal decimal.Decimal
}

type Product struct {
	Price   decimal.Decimal
	VatRate float64
}

func MakeLayout(products []Product) (*Layout, error) {
	// groupping products by vat rate
	m := make(map[float64]decimal.Decimal) // vat rate -> sum of all products with that vat rate
	for _, p := range products {
		m[p.VatRate] = m[p.VatRate].Add(p.Price)
	}

	// making buckets
	buckets := make(core.BucketSlice[float64], 0, len(m))
	for rate, amount := range m {
		bucket := core.Bucket[float64]{
			Key:   rate,
			Value: amount,
		}
		buckets = append(buckets, bucket)
	}

	// sorting buckets by vat rate
	sort.Slice(buckets, func(i, j int) bool {
		return buckets[i].Key < buckets[j].Key
	})

	// emainder will be added to the last bucket
	buckets[len(buckets)-1].ShouldAddRemainder = true

	coreLayout, err := core.MakeLayout(buckets)
	if err != nil {
		return nil, fmt.Errorf("unable to make core layout: %w", err)
	}
	return &Layout{
		Layout:        coreLayout,
		productsTotal: buckets.Total(),
	}, nil
}

func (l *Layout) DistributeDecimal(d decimal.Decimal) Value {
	v := l.Layout.DistributeDecimal(d, distributionPrecision)
	return Value{v}
}

func (l *Layout) GetProductsTotal() decimal.Decimal {
	return l.productsTotal
}

type Value struct {
	*core.Value[float64]
}
