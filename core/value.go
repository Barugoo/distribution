package core

import (
	"github.com/shopspring/decimal"
)

// Value is a result of Distribute operations; immutable
type Value[T comparable] struct {
	precision   int32
	bucketSlice BucketSlice[T]
	layout      *Layout[T]
}

func (v Value[T]) ToSlice() BucketSlice[T] {
	return v.bucketSlice.copy()
}

func (v Value[T]) Get(key T) (res decimal.Decimal, ok bool) {
	idx, ok := v.layout.idxMap[key]
	if !ok || idx >= len(v.bucketSlice) {
		// prevent panic if idx is out of bound
		return decimal.Zero, false
	}
	return v.bucketSlice[idx].Value, ok
}

func (v Value[T]) Total() decimal.Decimal {
	return v.bucketSlice.Total()
}

func (v Value[T]) GetLayout() *Layout[T] {
	return v.layout
}

func (v Value[T]) String() string {
	return v.bucketSlice.String()
}
