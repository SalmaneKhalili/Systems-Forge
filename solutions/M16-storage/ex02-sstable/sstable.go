package main

import (
	"encoding/binary"
	"os"
	"sort"
)

// KV is a key-value pair. The key is used to order records.
type KV struct {
	Key, Val string
}

// blockSize is the number of data records between sparse index entries.
const blockSize = 3

// indexEntry records the first key of a block and its byte offset in the
// data section.
type indexEntry struct {
	key    string
	offset int64
}

// SSTable is an open sstable: the raw file bytes plus the in-memory sparse
// index.
type SSTable struct {
	data     []byte
	dataOff  int
	index    []indexEntry
}

// WriteSSTable serializes sorted KV pairs to path.
func WriteSSTable(path string, items []KV) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var index []indexEntry
	buf := make([]byte, 0, 1<<16)
	for i, it := range items {
		if i%blockSize == 0 {
			index = append(index, indexEntry{key: it.Key, offset: int64(len(buf))})
		}
		buf = append(buf, it.Key...)
		buf = append(buf, 0)
		buf = append(buf, it.Val...)
		buf = append(buf, 0)
	}
	if _, err := f.Write(buf); err != nil {
		return err
	}

	// Tail: index records (key\0 + 4-byte little-endian offset) followed by a
	// 2-byte little-endian index offset at the very end of the file.
	var tail []byte
	for _, e := range index {
		tail = append(tail, e.key...)
		tail = append(tail, 0)
		tail = binary.LittleEndian.AppendUint32(tail, uint32(e.offset))
	}
	tail = binary.LittleEndian.AppendUint16(tail, uint16(len(buf)))
	if _, err := f.Write(tail); err != nil {
		return err
	}
	return f.Sync()
}

// OpenSSTable loads the sparse index from path.
func OpenSSTable(path string) (*SSTable, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) < 2 {
		return nil, os.ErrInvalid
	}
	dataOff := int(binary.LittleEndian.Uint16(data[len(data)-2:]))
	indexBytes := data[dataOff : len(data)-2]
	s := &SSTable{data: data, dataOff: dataOff, index: parseIndex(indexBytes)}
	return s, nil
}

func parseIndex(b []byte) []indexEntry {
	var out []indexEntry
	i := 0
	for i < len(b) {
		j := i
		for j < len(b) && b[j] != 0 {
			j++
		}
		key := string(b[i:j])
		if j+4 > len(b) {
			break
		}
		off := int64(binary.LittleEndian.Uint32(b[j+1 : j+5]))
		out = append(out, indexEntry{key: key, offset: off})
		i = j + 5
	}
	return out
}

// Get returns the value for key using the sparse index and binary search.
func (s *SSTable) Get(key string) (string, bool) {
	// Find the last block whose first key <= key.
	n := sort.Search(len(s.index), func(i int) bool { return s.index[i].key > key })
	if n == 0 {
		return "", false
	}
	blk := s.index[n-1]
	// Scan records from blk.offset until the data section ends.
	return s.scanFrom(blk.offset, key)
}

func (s *SSTable) scanFrom(start int64, key string) (string, bool) {
	data := s.data
	i := int(start)
	end := s.dataOff
	for i < end {
		// record: key\0 val\0
		j := i
		for j < end && data[j] != 0 {
			j++
		}
		if j >= end {
			return "", false
		}
		k := string(data[i:j])
		i = j + 1
		// value
		v := j + 1
		for v < end && data[v] != 0 {
			v++
		}
		if v >= end {
			return "", false
		}
		val := string(data[j+1 : v])
		i = v + 1
		if k == key {
			return val, true
		}
	}
	return "", false
}

// Blocks returns the number of sparse index entries.
func (s *SSTable) Blocks() int {
	return len(s.index)
}
