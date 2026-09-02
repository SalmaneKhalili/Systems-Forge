package main

import "sort"

// KV is a single key-value pair, ordered by Key.
type KV struct {
	Key, Val string
}

// MemTable is the write buffer of an LSM-tree: it holds recent writes in
// memory, keeping keys sorted for efficient scans and sequential flush.
type MemTable struct {
	entries map[string]string
}

// NewMemTable initializes an empty memtable.
func NewMemTable() *MemTable {
	return &MemTable{entries: make(map[string]string)}
}

// Put inserts a key-value pair, updating the value if the key exists.
func (m *MemTable) Put(key, val string) {
	m.entries[key] = val
}

// Get retrieves the value for a key, returning ("", false) if absent.
func (m *MemTable) Get(key string) (string, bool) {
	v, ok := m.entries[key]
	return v, ok
}

// All returns a copy of all pairs sorted in lexicographical key order.
func (m *MemTable) All() []KV {
	out := make([]KV, 0, len(m.entries))
	for k, v := range m.entries {
		out = append(out, KV{Key: k, Val: v})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}
