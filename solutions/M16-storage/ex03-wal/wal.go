package main

import (
	"encoding/binary"
	"os"
	"sort"
)

// KV is a key-value pair, ordered by Key.
type KV struct {
	Key, Val string
}

// LogEntry is one recorded write operation.
type LogEntry struct {
	Key, Val string
}

// AppendLog appends entries to an append-only log, creating it if needed.
// Each record is length-prefixed: uint16 keyLen, key, uint16 valLen, val.
func AppendLog(path string, entries []LogEntry) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	var buf []byte
	for _, e := range entries {
		buf = binary.LittleEndian.AppendUint16(buf, uint16(len(e.Key)))
		buf = append(buf, e.Key...)
		buf = binary.LittleEndian.AppendUint16(buf, uint16(len(e.Val)))
		buf = append(buf, e.Val...)
	}
	if _, err := f.Write(buf); err != nil {
		return err
	}
	return f.Sync()
}

// ReadLog reads back every entry in the log in order.
func ReadLog(path string) ([]LogEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out []LogEntry
	for i := 0; i < len(data); {
		if i+4 > len(data) {
			break
		}
		klen := int(binary.LittleEndian.Uint16(data[i : i+2]))
		i += 2
		if i+klen > len(data) {
			break
		}
		key := string(data[i : i+klen])
		i += klen
		if i+2 > len(data) {
			break
		}
		vlen := int(binary.LittleEndian.Uint16(data[i : i+2]))
		i += 2
		if i+vlen > len(data) {
			break
		}
		out = append(out, LogEntry{Key: key, Val: string(data[i : i+vlen])})
		i += vlen
	}
	return out, nil
}

// Replay reconstructs the key->value mapping by applying every entry in
// order, returning the result sorted by key.
func Replay(entries []LogEntry) []KV {
	m := map[string]string{}
	for _, e := range entries {
		m[e.Key] = e.Val
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]KV, 0, len(keys))
	for _, k := range keys {
		out = append(out, KV{Key: k, Val: m[k]})
	}
	return out
}