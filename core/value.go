package core

import (
	"github.com/shopspring/decimal"
)

// Value is a result of Distribute operations; immutable
type Value[T comparable] struct {
	bucketSlice BucketSlice[T]
	layout      *Layout[T]
}

// ToSlice returns a copy of the underlying BucketSlice
func (v Value[T]) ToSlice() BucketSlice[T] {
	return v.bucketSlice.copy()
}

// Get returns the value for a given bucket key. If the key does not exist, ok is false and res is zero
func (v Value[T]) Get(key T) (res decimal.Decimal, ok bool) {
	idx, ok := v.layout.idxMap[key]
	if !ok || idx >= len(v.bucketSlice) {
		// prevent panic if idx is out of bound
		return decimal.Zero, false
	}
	return v.bucketSlice[idx].Value, ok
}

// Total returns the total value of all buckets in this Value
func (v Value[T]) Total() decimal.Decimal {
	return v.bucketSlice.Total()
}

// GetLayout returns the Layout used to create this Value
func (v Value[T]) GetLayout() *Layout[T] {
	return v.layout
}

func (v Value[T]) String() string {
	return v.bucketSlice.String()
}
