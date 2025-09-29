package core

import "github.com/shopspring/decimal"

// toDecimal is a shorthand for decimal.NewFromFloat(f)
func toDecimal(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}
