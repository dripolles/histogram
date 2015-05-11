# histogram

Simple Histogram library, with support for interpolation and percentile calculation.

The Histogram stores a count of the number of occurrences for each value. This count
can later be used to interpolate expected counts for unknown values, and to
calculate the percentile position of a value (number of equal or smaller values).
The percentile calculation also does implicit interpolation if needed.

Sample usage
```
 h := NewHistogram()
 for _ , v := range values {
     h.Add(v)
 }
 p := GetPercentile(x)
```
