package core

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrBucketsTotalValueIsZero = errors.New("buckets total value is zero")
	ErrDuplicateBucketKey      = errors.New("duplicate bucket key")
)

// Layout serves as a distribution reference. Sum of all numerators divided by denominator is always equal to 1. Is immutable
type Layout[T comparable] struct {
	fractions []fraction[T]
	idxMap    map[T]int // for faster access to fraction by key
}

// fraction describes a ratio of 'bucket value / buckets total value' and additional info
type fraction[T comparable] struct {
	key T // bucket key
	numerator,
	denominator decimal.Decimal

	shouldAddRemainder bool // if true, the remainder will be added to resulting bucket
}

// DistributeDecimal distribute v based on a Layout. Sum of all buckets in the returned Value always equal to v
func (dl *Layout[T]) DistributeDecimal(v decimal.Decimal, precision int32) *Value[T] {
	bucketSlice := make(BucketSlice[T], 0, len(dl.fractions))

	var bucketSum decimal.Decimal

	// expectily setting out of bound index to check later if a new value was assigned in the loop
	remaindableBucketIdx := len(dl.fractions)
	for _, f := range dl.fractions {
		if f.denominator.IsZero() {
			// if somehow layout was created with zero denominator, we panic
			panic("layout fraction divisor cannot be zero")
		}

		bucket := Bucket[T]{
			Key: f.key,
			// bucket value = v * fraction
			Value:              v.Mul(f.numerator).Div(f.denominator).RoundBank(precision),
			ShouldAddRemainder: f.shouldAddRemainder,
		}
		// check if we should add the remainder to this bucket
		if bucket.ShouldAddRemainder {
			// saving idx of the new bucket
			remaindableBucketIdx = len(bucketSlice)
		}
		bucketSum = bucketSum.Add(bucket.Value)

		bucketSlice = append(bucketSlice, bucket)
	}

	// calculating remainder as v - sum of all buckets
	remainder := v.Sub(bucketSum)
	// if no fraction of the Layout was marked as remaindable
	if remaindableBucketIdx >= len(dl.fractions) {
		// add the remainder to the last bucket
		remaindableBucketIdx = len(bucketSlice) - 1
	}
	// adding the result to the remaindable bucket
	bucketSlice[remaindableBucketIdx].Value = bucketSlice[remaindableBucketIdx].Value.Add(remainder)

	c := Value[T]{
		bucketSlice: bucketSlice,
		layout:      dl,
	}
	return &c
}

// MakeLayout is a generic Layout constuctor. It takes a slice of buckets + optional LayoutOptions and returns a Layout. If sum of all bucket values zero, it returns an error. Duplicate bucket keys are not allowed
func MakeLayout[T comparable](buckets BucketSlice[T]) (*Layout[T], error) {
	total := buckets.Total()
	if total.IsZero() {
		// handling zero denominator
		return nil, ErrBucketsTotalValueIsZero
	}
	idxMap := make(map[T]int, len(buckets))
	fs := make([]fraction[T], 0, len(buckets))
	for _, b := range buckets {
		f := fraction[T]{
			key:         b.Key,
			numerator:   b.Value,
			denominator: total,

			shouldAddRemainder: b.ShouldAddRemainder,
		}
		if _, ok := idxMap[b.Key]; ok {
			// handling duplicate bucket keys
			return nil, ErrDuplicateBucketKey
		}
		idxMap[b.Key] = len(fs) // saving idx of the new bucket

		fs = append(fs, f)
	}
	l := Layout[T]{
		fractions: fs,
		idxMap:    idxMap,
	}
	return &l, nil
}

func (dl *Layout[T]) Keys() []T {
	keys := make([]T, 0, len(dl.fractions))
	for _, f := range dl.fractions {
		keys = append(keys, f.key)
	}
	return keys
}
