# M8-ex01 · Message framing

Framing is where a message layer starts: the 4-byte length prefix you write
here is the wire format your switch (ex03–ex05) will relay and that the M9+
replicas will speak, byte for byte. This exercise hardens that codec against
the problem that breaks every naive implementation — a reader that hands you
your bytes one at a time. You write a length-prefixed frame codec whose
decoder survives partial reads.

## Shape

Harness-style: the exercise ships `main.go` (the driver), `go.mod` and a
`Makefile`; you write **`frame.go`** in the same package and implement:

```go
func PutFrame(w io.Writer, msg []byte) error            // WriteFrame
func GetFrame(r io.Reader) ([]byte, error)              // ReadFrame
```

- **`PutFrame(w, msg)`** writes one frame: a 4-byte **big-endian** length
  followed by the payload. It must loop instead of trusting a single `Write`
  to flush everything.
- **`GetFrame(r)`** reads one frame back, tolerating `Read` calls that return
  almost nothing. It must return:
  - the payload and `nil` for a complete frame;
  - `nil, io.EOF` when the stream ends exactly at a frame boundary;
  - `nil, io.ErrUnexpectedEOF` when the stream ends mid-header or mid-payload.

## Acceptance

`make all`, then `./test` must exit 0 and print exactly:

```text
roundtrip: hello|world|x
split: alpha|beta
truncated: unexpected eof
```

- `roundtrip` — three messages packed and read back from an ordinary buffer.
- `split` — the **same** messages read back from a reader that returns **one
  byte per `Read`**. A decoder that trusts a single header read or a single
  payload read fails this line (or the program dies).
- `truncated` — a frame cut off halfway must surface as `io.ErrUnexpectedEOF`,
  never as a silently short message and never as a hang.

Stderr must stay empty. A clean exit (0) is part of the grade. Graded
`build` + `stdout` + `quiz`.

## Readings

- Go `io` docs: https://pkg.go.dev/io#Reader, https://pkg.go.dev/io#ReadFull,
  https://pkg.go.dev/io#EOF
- Go `encoding/binary` docs → `BigEndian`.
- A. Tanenbaum, *Computer Networks* — framing / byte stuffing chapter.

## Quiz

1. Which 4-byte integer in each frame tells the reader how many payload bytes follow?
2. Which error must `GetFrame` return when the stream ends halfway through a frame?
3. Which `io` helper loops internally so the codec survives one-byte-at-a-time readers?