package histogram

import (
	"math"
	"sort"
)

type Histogram struct {
	values       map[int]int
	count        int
	minV         *int
	maxV         *int
	accHist      map[int]float64
	sortedValues []int
}

func NewHistogram() *Histogram {
	return &Histogram{
		values: map[int]int{},
	}
}

func (h *Histogram) Add(v int) {
	h.accHist = nil

	h.updateBoundaryValues(v)

	count := h.values[v] + 1
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

func (h *Histogram) GetPercentile(v int) float64 {

	if h.minV == nil || h.maxV == nil {
		panic("Uninitialized Histogram")
	}

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
		count = h.interpolateValue(v)
	}

	return count / float64(h.Len())
}

func (h *Histogram) initAccHist() {
	h.accHist = map[int]float64{}

	values := make([]int, 0, len(h.values))
	for v, _ := range h.values {
		values = append(values, v)
	}

	sort.Ints(values)
	h.sortedValues = values
	total := 0.0
	for _, v := range values {
		count, _ := h.values[v]
		total += float64(count)
		h.accHist[v] = total
	}
}

func (h *Histogram) interpolateValue(v int) float64 {
	x0, x1 := h.neighbours(v)
	y0 := h.accHist[x0]
	y1 := h.accHist[x1]

	result := y0 + (y1-y0)*(float64(v-x0)/float64(x1-x0))
	h.accHist[v] = result

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

func (h *Histogram) previous(v int) (pos int, val float64) {
	for pos = v; pos >= *h.minV; pos-- {
		if val, ok := h.accHist[pos]; ok {
			return pos, val
		}
	}

	return 0, math.NaN()
}

func (h *Histogram) next(v int) (pos int, val float64) {
	for pos = v; pos <= *h.maxV; pos++ {
		if val, ok := h.accHist[pos]; ok {
			return pos, val
		}
	}

	return 0, math.NaN()
}
