package core

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type Bucket[T comparable] struct {
	Key                T
	Value              decimal.Decimal
	ShouldAddRemainder bool
}

type BucketSlice[T comparable] []Bucket[T]

func (s BucketSlice[T]) copy() BucketSlice[T] {
	cloned := make(BucketSlice[T], len(s))
	copy(cloned, s)
	return cloned
}

func (s BucketSlice[T]) Total() (sum decimal.Decimal) {
	for _, v := range s {
		sum = sum.Add(v.Value)
	}
	return sum
}

func (s BucketSlice[T]) String() string {
	var str strings.Builder
	str.WriteString("[")
	for idx, v := range s {
		str.WriteString(fmt.Sprintf("%v: %s", v.Key, v.Value))
		if idx != len(s)-1 {
			str.WriteString(", ")
		}
	}
	str.WriteString("]")
	return str.String()
}
