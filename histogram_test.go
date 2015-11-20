package histogram

import (
	"math/rand"
	"sort"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type HistogramSuite struct{}

var _ = Suite(&HistogramSuite{})

func (s *HistogramSuite) TestHistogram_Len(c *C) {
	h := NewHistogram()
	h.Add(1)
	h.Add(1)
	h.Add(2)

	c.Assert(h.Len(), Equals, 3)
}

func (s *HistogramSuite) TestHistogram_Get(c *C) {
	h := createSampleHistogram()

	assertHistogramGet(c, h, 100, 2.0)
	assertHistogramGet(c, h, 300, 0.0)
	assertHistogramGet(c, h, 500, 4.0)
	assertHistogramGet(c, h, 700, 0.0)
	assertHistogramGet(c, h, 900, 3.0)
	assertHistogramGet(c, h, 1000, 1.0)
}

func (s *HistogramSuite) TestHistogram_GetInterpolated(c *C) {
	h := createSampleHistogram()

	assertHistogramGetInterpolated(c, h, 100, 2.0)
	assertHistogramGetInterpolated(c, h, 300, 3.0)
	assertHistogramGetInterpolated(c, h, 500, 4.0)
	assertHistogramGetInterpolated(c, h, 700, 3.5)
	assertHistogramGetInterpolated(c, h, 900, 3.0)
	assertHistogramGetInterpolated(c, h, 1000, 1.0)
}

func (s *HistogramSuite) TestHistogram_GetOutOfBounds(c *C) {
	h := createSampleHistogram()

	_, err := h.Get(99)
	c.Assert(err, DeepEquals, ErrExtrapolation)
	_, err = h.Get(1001)
	c.Assert(err, DeepEquals, ErrExtrapolation)
}

func assertHistogramGet(c *C, h *Histogram, v int, expected float64) {
	res, err := h.Get(v)
	c.Assert(err, IsNil)
	c.Assert(res, Equals, expected)
}

func assertHistogramGetInterpolated(c *C, h *Histogram, v int, expected float64) {
	res, err := h.GetInterpolated(v)
	c.Assert(err, IsNil)
	c.Assert(res, Equals, expected)
}

func (s *HistogramSuite) TestHistogram_Percentile(c *C) {
	h := createSampleHistogram()

	c.Assert(h.GetPercentile(0), Equals, 0.0)
	c.Assert(h.GetPercentile(100), Equals, 0.2)
	c.Assert(h.GetPercentile(300), Equals, 0.4)
	c.Assert(h.GetPercentile(500), Equals, 0.6)
	c.Assert(h.GetPercentile(700), Equals, 0.75)
	c.Assert(h.GetPercentile(900), Equals, 0.9)
	c.Assert(h.GetPercentile(1000), Equals, 1.0)
	c.Assert(h.GetPercentile(1001), Equals, 1.0)
}

func createSampleHistogram() *Histogram {
	values := []int{100, 100, 500, 500, 500, 500, 900, 900, 900, 1000}
	h := NewHistogram()
	for _, v := range values {
		h.Add(v)
	}

	return h
}

func (s *HistogramSuite) TestHistogram_PercentileRandom(c *C) {
	numValues := 10000
	maxValue := 1000 * 1000 * 1000 * 1000
	h, values := makeHistogram(numValues, maxValue)

	sort.Ints(values)
	v := values[0]
	for i := 1; i < numValues; i++ {
		next := values[i]
		if v != next {
			p := float64(i) / float64(numValues)
			c.Assert(h.GetPercentile(v), Equals, p)
			v = next
		}
	}
	c.Assert(h.GetPercentile(values[numValues-1]), Equals, 1.0)
}

func (s *HistogramSuite) BenchmarkHistogram_PercentileInterpolate(c *C) {
	numValues := 100000
	maxValue := 1000 * 1000 * 1000 * 1000
	h, _ := makeHistogram(numValues, maxValue)

	for i := 0; i < c.N; i++ {
		v := rand.Intn(maxValue)
		h.GetPercentile(v)
		return
	}
}

func (s *HistogramSuite) TestHistogram_GetAtPercentile(c *C) {
	h := NewHistogram()
	values := []int{1, 1, 2, 2, 3}
	for _, v := range values {
		h.Add(v)
	}

	c.Assert(h.GetAtPercentile(0.1), Equals, 1)
	c.Assert(h.GetAtPercentile(0.25), Equals, 1)
	c.Assert(h.GetAtPercentile(0.4), Equals, 1)
	c.Assert(h.GetAtPercentile(0.45), Equals, 2)
	c.Assert(h.GetAtPercentile(0.50), Equals, 2)
	c.Assert(h.GetAtPercentile(0.80), Equals, 2)
	c.Assert(h.GetAtPercentile(0.85), Equals, 3)
	c.Assert(h.GetAtPercentile(0.99), Equals, 3)
	c.Assert(h.GetAtPercentile(1.0), Equals, 3)

}

func (s *HistogramSuite) TestHistogram_Neighbours(c *C) {
	h := NewHistogram()
	numValues := 100
	for i := 0; i < numValues; i++ {
		h.Add(i * 2)
	}
	h.initSortedValues()

	for i := 0; i < numValues-1; i++ {
		v := i*2 + 1
		p, n := h.neighbours(i*2 + 1)
		c.Assert(p, Equals, v-1)
		c.Assert(n, Equals, v+1)
	}
}

func makeHistogram(num, max int) (h *Histogram, values []int) {
	values = make([]int, 0, num)
	h = NewHistogram()

	for i := 0; i < num; i++ {
		v := rand.Intn(max)
		values = append(values, v)
		h.Add(v)
	}

	return
}
