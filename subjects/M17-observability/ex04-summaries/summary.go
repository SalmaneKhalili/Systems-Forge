package main

import (
	"math"
	"sort"
	"strconv"
)

// KV is a key/value pair, Val rendered as a string.
type KV struct {
	Key, Val string
}

// Summary aggregates observed durations in milliseconds.
type Summary struct {
	vals []int64
}

// NewSummary initializes an empty summary.
func NewSummary() *Summary {
	return &Summary{}
}

// Add records one observed duration in ms.
func (s *Summary) Add(durMs int64) {
	s.vals = append(s.vals, durMs)
}

// Count returns the number of observations.
func (s *Summary) Count() int64 {
	return int64(len(s.vals))
}

// Sum returns the total of all observations.
func (s *Summary) Sum() int64 {
	var total int64
	for _, v := range s.vals {
		total += v
	}
	return total
}

// Max returns the largest observation, or 0 if empty.
func (s *Summary) Max() int64 {
	var best int64
	for _, v := range s.vals {
		if v > best {
			best = v
		}
	}
	return best
}

// Quantile returns the smallest observation at index ceil(q*n)-1 on the
// sorted data (1-indexed), or 0 for empty data.
func (s *Summary) Quantile(q float64) float64 {
	n := len(s.vals)
	if n == 0 {
		return 0
	}
	sorted := append([]int64(nil), s.vals...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	idx := int(math.Ceil(q*float64(n))) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= n {
		idx = n - 1
	}
	return float64(sorted[idx])
}

// SummaryValues returns count, sum, max, p50, p95 as name/value pairs.
func (s *Summary) SummaryValues() []KV {
	out := []KV{
		{Key: "count", Val: strconv.FormatInt(s.Count(), 10)},
		{Key: "sum", Val: strconv.FormatInt(s.Sum(), 10)},
		{Key: "max", Val: strconv.FormatInt(s.Max(), 10)},
		{Key: "p50", Val: formatF(s.Quantile(0.5))},
		{Key: "p95", Val: formatF(s.Quantile(0.95))},
	}
	return out
}

func formatF(v float64) string {
	if v == math.Trunc(v) {
		return strconv.FormatInt(int64(v), 10)
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}