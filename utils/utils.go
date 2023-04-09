package utils

import "github.com/shopspring/decimal"

// ToDecimal is a shorthand for decimal.NewFromFloat(f)
func ToDecimal(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}
