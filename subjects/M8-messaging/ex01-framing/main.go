package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

// chunkReader serves its data one byte per Read call, so any codec that
// trusts a single read to fill a header or a payload will misbehave.
type chunkReader struct{ data []byte }

func (c *chunkReader) Read(p []byte) (int, error) {
	if len(c.data) == 0 {
		return 0, io.EOF
	}
	p[0] = c.data[0]
	c.data = c.data[1:]
	return 1, nil
}

func encode(msgs []string) []byte {
	var buf bytes.Buffer
	for _, m := range msgs {
		if err := PutFrame(&buf, []byte(m)); err != nil {
			panic(err)
		}
	}
	return buf.Bytes()
}

func decode(r io.Reader, count int) [][]byte {
	out := make([][]byte, 0, count)
	for i := 0; i < count; i++ {
		m, err := GetFrame(r)
		if err != nil {
			panic(err)
		}
		out = append(out, m)
	}
	return out
}

func join(msgs [][]byte) string {
	parts := make([]string, len(msgs))
	for i, m := range msgs {
		parts[i] = string(m)
	}
	return strings.Join(parts, "|")
}

func main() {
	blob := encode([]string{"hello", "world", "x"})

	got := join(decode(bytes.NewReader(blob), 3))
	fmt.Printf("roundtrip: %s\n", got)

	blob2 := encode([]string{"alpha", "beta"})
	got = join(decode(&chunkReader{data: blob2}, 2))
	fmt.Printf("split: %s\n", got)

	trunc := encode([]string{"gamma"})
	half := trunc[:len(trunc)/2+1] // header plus one payload byte: mid-frame EOF
	m, err := GetFrame(&chunkReader{data: half})
	if err == nil {
		fmt.Printf("truncated: got %q\n", string(m))
	} else if errors.Is(err, io.ErrUnexpectedEOF) {
		fmt.Println("truncated: unexpected eof")
	} else {
		fmt.Printf("truncated: %v\n", err)
	}
}