package core

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestValueToSlice(t *testing.T) {
	tt := []struct {
		name     string
		value    *Value[string]
		expected BucketSlice[string]
	}{
		{
			name:     "empty",
			value:    &Value[string]{},
			expected: BucketSlice[string]{},
		},
		{
			name: "single",
			value: &Value[string]{
				bucketSlice: BucketSlice[string]{
					{
						Key:   "10",
						Value: toDecimal(10),
					},
				},
			},
			expected: BucketSlice[string]{
				{
					Key:   "10",
					Value: toDecimal(10),
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.value.ToSlice())
		})
	}
}

func TestValueTotal(t *testing.T) {
	tt := []struct {
		name     string
		value    *Value[string]
		expected decimal.Decimal
	}{
		{
			name:     "empty",
			value:    &Value[string]{},
			expected: toDecimal(0),
		},
		{
			name: "single",
			value: &Value[string]{
				bucketSlice: BucketSlice[string]{
					{
						Key:   "10",
						Value: toDecimal(10),
					},
				},
			},
			expected: toDecimal(10),
		},
		{
			name: "multiple",
			value: &Value[string]{
				bucketSlice: BucketSlice[string]{
					{
						Key:   "10",
						Value: toDecimal(10),
					},
					{
						Key:   "20",
						Value: toDecimal(20),
					},
				},
			},
			expected: toDecimal(30),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assertDecimal(t, tc.expected, tc.value.Total())
		})
	}
}

func TestValueGet(t *testing.T) {
	tt := []struct {
		name          string
		value         *Value[string]
		key           string
		expected      decimal.Decimal
		expectedNotOK bool
	}{
		{
			name: "empty",
			value: &Value[string]{
				layout: &Layout[string]{
					idxMap: map[string]int{},
				},
			},
			key:           "10",
			expected:      toDecimal(0),
			expectedNotOK: true,
		},
		{
			name: "single",
			value: &Value[string]{
				layout: &Layout[string]{
					idxMap: map[string]int{
						"10": 0,
					},
				},
				bucketSlice: BucketSlice[string]{
					{
						Key:   "10",
						Value: toDecimal(10),
					},
				},
			},
			key:      "10",
			expected: toDecimal(10),
		},
		{
			name: "multiple",
			value: &Value[string]{
				layout: &Layout[string]{
					idxMap: map[string]int{
						"10": 0,
						"20": 1,
					},
				},
				bucketSlice: BucketSlice[string]{
					{
						Key:   "10",
						Value: toDecimal(10),
					},
					{
						Key:   "20",
						Value: toDecimal(20),
					},
				},
			},
			key:      "20",
			expected: toDecimal(20),
		},
		{
			name: "idx out of range",
			value: &Value[string]{
				layout: &Layout[string]{
					idxMap: map[string]int{
						"10":         0,
						"20":         1,
						"whatisthis": 999,
					},
				},
				bucketSlice: BucketSlice[string]{
					{
						Key:   "10",
						Value: toDecimal(10),
					},
					{
						Key:   "20",
						Value: toDecimal(20),
					},
				},
			},
			key:           "whatisthis",
			expected:      toDecimal(0),
			expectedNotOK: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			res, ok := tc.value.Get(tc.key)
			assert.Equal(t, tc.expectedNotOK, !ok)
			assertDecimal(t, tc.expected, res)
		})
	}
}

func TestValueGetLayout(t *testing.T) {
	tt := []struct {
		name     string
		value    *Value[string]
		expected *Layout[string]
	}{
		{
			name:     "empty",
			value:    &Value[string]{},
			expected: nil,
		},
		{
			name: "single",
			value: &Value[string]{
				layout: &Layout[string]{
					idxMap: map[string]int{
						"10": 0,
					},
				},
			},
			expected: &Layout[string]{
				idxMap: map[string]int{
					"10": 0,
				},
			},
		},
		{
			name: "multiple",
			value: &Value[string]{
				layout: &Layout[string]{
					idxMap: map[string]int{
						"10": 0,
						"20": 1,
					},
				},
			},
			expected: &Layout[string]{
				idxMap: map[string]int{
					"10": 0,
					"20": 1,
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.value.GetLayout())
		})
	}
}
