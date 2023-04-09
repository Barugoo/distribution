# Distribution package

## About 
This package is used to split arbitrary decimal.Decimal values (usually money amounts) based on `distribution.Layout` that describes weights of each bucket, i.e.: 
```
{
    bucket1: 1/3, 
    bucket2: 1/6, 
    bucket3: 1/2
}

sum of all fractions is always 1
```

The result of "splitting" operation is a `distribution.Value` that contains the buckets with calculated shares of initial value (with specific precision) and the remainder.  
For example distribution of `100` using the above `distribution.Layout` will produce this result:
```
{
    precision: 2, 
    bucket1: 33.33, 
    bucket2: 16.66, 
    bucket3: 50, 
    remainder: 0.01
}

sum is always 100
```
The `remainder` will be added automatically based on which `distribution.Bucket` in variadic slice of `MakeLayout` was marked as `remaindable`

External code doesn't have access to internals of `Value` or `Layout` and there are no mutating methods in the public API - which means there's no way for an importing code to break the internal rules of this package.
