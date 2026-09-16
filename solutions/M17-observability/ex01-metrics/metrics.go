package main

import (
	"sort"
	"strconv"
	"sync"
)

// KV is a name/value pair, Val rendered as a string.
type KV struct {
	Key, Val string
}

// Metrics is a concurrency-safe registry of named counters.
type Metrics struct {
	mu   sync.RWMutex
	data map[string]int64
}

// NewMetrics initializes an empty registry.
func NewMetrics() *Metrics {
	return &Metrics{data: make(map[string]int64)}
}

// Inc increments the named counter by 1.
func (m *Metrics) Inc(name string) {
	m.Add(name, 1)
}

// Add increments the named counter by n — a write, so it takes the exclusive Lock.
func (m *Metrics) Add(name string, n int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[name] += n
}

// Get returns the current value of a counter and whether it exists — a read, so RLock.
func (m *Metrics) Get(name string) (int64, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[name]
	return v, ok
}

// Snapshot returns all counters as name/value pairs sorted by name — a read, so RLock.
func (m *Metrics) Snapshot() []KV {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, 0, len(m.data))
	for k := range m.data {
		names = append(names, k)
	}
	sort.Strings(names)
	out := make([]KV, 0, len(names))
	for _, k := range names {
		out = append(out, KV{Key: k, Val: strconv.FormatInt(m.data[k], 10)})
	}
	return out
}