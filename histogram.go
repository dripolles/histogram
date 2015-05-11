package histogram

import (
	"errors"
	"sort"
)

var ExtrapolationError error = errors.New("Extrapolation of histogram values not supported")

type Histogram struct {
	values       map[int]float64
	count        int
	minV         *int
	maxV         *int
	accHist      map[int]float64
	sortedValues []int
}

func NewHistogram() *Histogram {
	return &Histogram{
		values: map[int]float64{},
	}
}

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
