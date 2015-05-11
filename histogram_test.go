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
	values := []int{100, 100, 500, 500, 500, 500, 900, 900, 900, 1000}
	h := NewHistogram()
	for _, v := range values {
		h.Add(v)
	}

    assertHistogramGet(c, h, 100, 2.0)
    assertHistogramGet(c, h, 300, 3.0)
    assertHistogramGet(c, h, 500, 4.0)
    assertHistogramGet(c, h, 700, 3.5)
    assertHistogramGet(c, h, 900, 3.0)
    assertHistogramGet(c, h, 1000, 1.0)
}

func assertHistogramGet(c*C, h*Histogram, v int, expected float64) {
    res, err := h.Get(v)
    c.Assert(err, IsNil)
    c.Assert(res, Equals, expected)
}

func (s *HistogramSuite) TestHistogram_Percentile(c *C) {
	values := []int{100, 100, 500, 500, 500, 500, 900, 900, 900, 1000}
	h := NewHistogram()
	for _, v := range values {
		h.Add(v)
	}

	c.Assert(h.GetPercentile(0), Equals, 0.0)
	c.Assert(h.GetPercentile(100), Equals, 0.2)
	c.Assert(h.GetPercentile(300), Equals, 0.4)
	c.Assert(h.GetPercentile(500), Equals, 0.6)
	c.Assert(h.GetPercentile(700), Equals, 0.75)
	c.Assert(h.GetPercentile(900), Equals, 0.9)
	c.Assert(h.GetPercentile(1000), Equals, 1.0)
	c.Assert(h.GetPercentile(1001), Equals, 1.0)
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
	for i := 0; i < c.N; i++ {
		h := NewHistogram()
		h.Add(1)
		h.Add(1000000)
		for v := 100000; v < 1000000; v += 50000 {
			h.GetPercentile(v)
		}
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
