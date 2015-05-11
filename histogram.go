// Simple Histogram lib, with support for interpolation and percentile calculation.
//
// The Histogram stores a count of the number of occurrences for each value. This count
// can later be used to interpolate expected counts for unknown values, and to
// calculate the percentile position of a value (number of equal or smaller values).
// The percentile calculation also does implicit interpolation if needed.
//
// Sample usage
//  h := NewHistogram()
//  for _ , v := range values {
//      h.Add(v)
//  }
//  p := GetPercentile(x)
//
package histogram

import (
	"errors"
	"sort"
)

// This library does not support extrapolation of histogram values outside the range
// of known values. An ExtrapolationError is used if such extrapolation is attempted.
var ExtrapolationError error = errors.New("Extrapolation of histogram values not supported")

type Histogram struct {
	values       map[int]float64
	count        int
	minV         *int
	maxV         *int
	accHist      map[int]float64
	sortedValues []int
}

// Create new empty Histogram.
func NewHistogram() *Histogram {
	return &Histogram{
		values: map[int]float64{},
	}
}

// Add value to the Histogram.
func (h *Histogram) Add(v int) {
	h.accHist = nil
	h.sortedValues = nil

	h.updateBoundaryValues(v)

	count := h.values[v] + 1.0
	h.values[v] = count

	h.count += 1
}

func (h *Histogram) updateBoundaryValues(v int) {
	if h.minV == nil || *h.minV > v {
		h.minV = &v
	}

	if h.maxV == nil || *h.maxV < v {
		h.maxV = &v
	}
}

// Numnber of non unique values inserted in the Histogram.
func (h *Histogram) Len() int {
	return h.count
}

func (h *Histogram) Sum() int {
	sum := 0
	for value, _ := range h.values {
		sum += value
	}

	return sum
}

func (h *Histogram) checkInitialized() {
	if h.minV == nil || h.maxV == nil {
		panic("Uninitialized Histogram")
	}
}

// Get count for value v. Automatically performs interpolation if needed.
func (h *Histogram) Get(v int) (float64, error) {
	h.checkInitialized()

	if v < *h.minV || v > *h.maxV {
		return 0, ExtrapolationError
	}

	count, ok := h.values[v]
	if !ok {
		count = h.interpolateValue(v, h.values)
	}

	return count, nil
}

// Get percentile position for value v. Automatically performs interpolation if needed.
func (h *Histogram) GetPercentile(v int) float64 {
	if v < *h.minV {
		return 0.0
	}

	if v > *h.maxV {
		return 1.0
	}

	if h.accHist == nil {
		h.initAccHist()
	}

	count, ok := h.accHist[v]
	if !ok {
		count = h.interpolateValue(v, h.accHist)
	}

	return count / float64(h.Len())
}

func (h *Histogram) initAccHist() {
	h.initSortedValues()
	h.accHist = map[int]float64{}

	total := 0.0
	for _, v := range h.sortedValues {
		count, _ := h.values[v]
		total += float64(count)
		h.accHist[v] = total
	}
}

func (h *Histogram) initSortedValues() {
	if h.sortedValues != nil {
		return
	}
	values := make([]int, 0, len(h.values))
	for v, _ := range h.values {
		values = append(values, v)
	}

	sort.Ints(values)
	h.sortedValues = values
}

func (h *Histogram) interpolateValue(v int, m map[int]float64) float64 {
	h.initSortedValues()

	x0, x1 := h.neighbours(v)
	y0 := m[x0]
	y1 := m[x1]

	result := y0 + (y1-y0)*(float64(v-x0)/float64(x1-x0))

	return result
}

func (h *Histogram) neighbours(v int) (prev, next int) {
	l := len(h.sortedValues)
	i := l / 2

	for {
		prev = h.sortedValues[i]
		if prev <= v {
			next = h.sortedValues[i+1]
			if next >= v {
				return
			} else {
				i += (l - i) / 2
			}
		} else {
			i = i / 2
		}
	}
}
