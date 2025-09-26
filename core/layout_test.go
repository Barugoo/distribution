package core

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/barugoo/distribution/utils"
)

func fromStr(str string) decimal.Decimal {
	strDec, _ := decimal.NewFromString(str)
	return strDec
}

func newLayout[T comparable](fractions []fraction[T]) *Layout[T] {
	l := &Layout[T]{
		fractions: fractions,
	}
	m := make(map[T]int, len(fractions))
	for i, f := range fractions {
		m[f.key] = i
	}
	l.idxMap = m
	return l
}

func TestMakeLayout(t *testing.T) {
	tt := []struct {
		input          []Bucket[string]
		expectedLayout *Layout[string]
		expectedErr    error
	}{
		{
			input: []Bucket[string]{
				{
					Key:   "15",
					Value: fromStr("11"),
				},
				{
					Key:                "30",
					Value:              fromStr("1.95"),
					ShouldAddRemainder: true,
				},
			},
			expectedLayout: newLayout([]fraction[string]{
				{
					key:         "15",
					numerator:   fromStr("11"),
					denominator: fromStr("12.95"),
				},
				{
					key:                "30",
					numerator:          fromStr("1.95"),
					denominator:        fromStr("12.95"),
					shouldAddRemainder: true,
				},
			},
			),
		},
		{
			input: []Bucket[string]{
				{
					Key:   "15",
					Value: fromStr("100"),
				},
				{
					Key:   "30",
					Value: fromStr("100"),
				},
				{
					Key:                "45",
					Value:              fromStr("100"),
					ShouldAddRemainder: true,
				},
			},
			expectedLayout: newLayout([]fraction[string]{
				{
					key:         "15",
					numerator:   fromStr("100"),
					denominator: fromStr("300"),
				},
				{
					key:         "30",
					numerator:   fromStr("100"),
					denominator: fromStr("300"),
				},
				{
					key:                "45",
					numerator:          fromStr("100"),
					denominator:        fromStr("300"),
					shouldAddRemainder: true,
				},
			}),
		},
		{
			input: []Bucket[string]{
				{
					Key:   "15",
					Value: fromStr("0"),
				},
			},
			expectedLayout: nil,
			expectedErr:    ErrBucketsTotalValueIsZero,
		},
		{
			input: []Bucket[string]{
				{
					Key:   "15",
					Value: fromStr("10"),
				},
				{
					Key:   "15",
					Value: fromStr("15"),
				},
			},
			expectedLayout: nil,
			expectedErr:    ErrDuplicateBucketKey,
		},
	}

	for _, tc := range tt {
		res, err := MakeLayout(tc.input)
		assert.Equal(t, tc.expectedErr, err)

		assertLayout(t, tc.expectedLayout, res)
	}
}

func assertDecimal(t *testing.T, expected, actual decimal.Decimal) bool {
	return assert.Truef(t, expected.Equal(actual), "expected %s; got %s", expected, actual)
}

func assertLayout[T comparable](t *testing.T, expected, actual *Layout[T]) bool {
	if expected == nil && actual == nil {
		return true
	}
	if expected == nil || actual == nil {
		return false
	}
	for _, f := range expected.fractions {
		actualF := actual.fractions[actual.idxMap[f.key]]
		if !assertDecimal(t, f.numerator, actualF.numerator) {
			return false
		}
		if !assertDecimal(t, f.denominator, actualF.denominator) {
			return false
		}
		assert.Equal(t, f.shouldAddRemainder, actualF.shouldAddRemainder)
	}
	return true
}

func assertValue[T comparable](t *testing.T, expected, actual *Value[T]) {
	assert.Equal(t, expected.bucketSlice, expected.bucketSlice)
	assert.Equal(t, expected.layout, actual.layout)
}

func TestDistributeDecimal(t *testing.T) {
	tt := []struct {
		name          string
		input         decimal.Decimal
		layout        *Layout[string]
		expectedValue *Value[string]
	}{
		{
			name:  "simple",
			input: utils.ToDecimal(12.95),
			layout: newLayout([]fraction[string]{
				{
					key:         "15",
					numerator:   fromStr("11"),
					denominator: fromStr("12.95"),
				},
				{
					key:                "30",
					numerator:          fromStr("1.95"),
					denominator:        fromStr("12.95"),
					shouldAddRemainder: true,
				},
			},
			),
			expectedValue: &Value[string]{
				bucketSlice: BucketSlice[string]{
					{"15", utils.ToDecimal(11), false},
					{"30", utils.ToDecimal(1.95), true},
				},
			},
		},
		{
			name:  "divisor is dividable by 3 and dividend is not",
			input: utils.ToDecimal(10),
			layout: newLayout([]fraction[string]{
				{
					key:         "15",
					numerator:   fromStr("100"),
					denominator: fromStr("300"),
				},
				{
					key:         "30",
					numerator:   fromStr("100"),
					denominator: fromStr("300"),
				},
				{
					key:                "45",
					numerator:          fromStr("100"),
					denominator:        fromStr("300"),
					shouldAddRemainder: true,
				},
			}),
			expectedValue: &Value[string]{
				bucketSlice: BucketSlice[string]{
					{"15", utils.ToDecimal(33.33), false},
					{"30", utils.ToDecimal(33.33), false},
					{"45", utils.ToDecimal(33.34), true},
				},
			},
		},
		{
			name:  "divisor and dividend are dividable by 3",
			input: utils.ToDecimal(15),
			layout: newLayout([]fraction[string]{
				{
					key:         "15",
					numerator:   fromStr("100"),
					denominator: fromStr("300"),
				},
				{
					key:         "30",
					numerator:   fromStr("100"),
					denominator: fromStr("300"),
				},
				{
					key:                "45",
					numerator:          fromStr("100"),
					denominator:        fromStr("300"),
					shouldAddRemainder: true,
				},
			}),
			expectedValue: &Value[string]{
				bucketSlice: BucketSlice[string]{
					{"15", utils.ToDecimal(5), false},
					{"30", utils.ToDecimal(5), false},
					{"45", utils.ToDecimal(5), true},
				},
			},
		},
		{
			name:  "round bank case",
			input: utils.ToDecimal(0.5),
			layout: newLayout([]fraction[string]{
				{
					key:         "15",
					numerator:   fromStr("1"),
					denominator: fromStr("100"),
				},
				{
					key:         "30",
					numerator:   fromStr("99"),
					denominator: fromStr("100"),
				},
			}),
			expectedValue: &Value[string]{
				bucketSlice: BucketSlice[string]{
					{"15", utils.ToDecimal(0.50), false},
					{"30", utils.ToDecimal(0.00), true},
				},
			},
		},
		{
			name:  "division by 2 with remainder",
			input: utils.ToDecimal(325.01),
			layout: newLayout([]fraction[string]{
				{
					key:         "15",
					numerator:   fromStr("1"),
					denominator: fromStr("2"),
				},
				{
					key:         "30",
					numerator:   fromStr("1"),
					denominator: fromStr("2"),
				},
			}),
			expectedValue: &Value[string]{
				bucketSlice: BucketSlice[string]{
					{"15", utils.ToDecimal(162.5), false},
					{"30", utils.ToDecimal(162.51), true},
				},
			},
		},
	}

	for _, tc := range tt {
		tc.expectedValue.layout = tc.layout

		res := tc.layout.DistributeDecimal(tc.input, 2)
		assertValue(t, tc.expectedValue, res)
	}
}
