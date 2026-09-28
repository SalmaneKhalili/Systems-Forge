# M8-ex01 — message framing

## Goal

A message layer needs a wire format: a stream of bytes must be cut back into
the exact messages that were written. You implement a **length-prefixed frame
codec** and a decoder that survives the Internet's favourite trick: a reader
that hands you your bytes **one at a time**.

## Shape

Harness-style (like M3/M5). The exercise ships `main.go` (the driver), `go.mod`
and a `Makefile`. You write **`frame.go`** in the same package and implement two
functions:

```go
func PutFrame(w io.Writer, msg []byte) error            // WriteFrame
func GetFrame(r io.Reader) ([]byte, error)              // ReadFrame
```

- **`PutFrame(w, msg)`** writes one frame to `w`: a 4-byte **big-endian** length
  followed by the payload bytes. It must loop instead of trusting a single
  `Write` to flush everything.
- **`GetFrame(r)`** reads one frame back, tolerating `Read` calls that return
  almost nothing. It must return:
  - the payload and `nil` for a complete frame;
  - `nil, io.EOF` when the stream ends exactly at a frame boundary;
  - `nil, io.ErrUnexpectedEOF` when the stream ends mid-header or mid-payload.

## Acceptance criteria (all graded)

`make all`, then `./test` must exit 0 and print exactly:

```text
roundtrip: hello|world|x
split: alpha|beta
truncated: unexpected eof
```

- `roundtrip` — three messages packed and read back from an ordinary buffer.
- `split` — the **same** messages read back from a reader that returns **one byte
  per `Read`**. If your decoder trusts a single header read or a single payload
  read, this line comes out wrong (or the program dies).
- `truncated` — a frame truncated halfway must surface as `io.ErrUnexpectedEOF`,
  never as a silently short message and never as a hang.

Stderr must stay empty. A clean exit (0) is part of the grade.

## Readings

- Go `io` docs: https://pkg.go.dev/io#Reader, https://pkg.go.dev/io#ReadFull,
  https://pkg.go.dev/io#EOF
- Go `encoding/binary` docs → `BigEndian`.
- Go `encoding/binary` docs → `BigEndian`.
- A. Tanenbaum, *Computer Networks* — framing/byte stuffing chapter.

## Quiz

1. Which 4-byte integer in each frame tells the reader how many payload bytes follow?
2. Which error must `GetFrame` return when the stream ends halfway through a frame?
3. Which `io` helper loops internally so the codec survives one-byte-at-a-time readers?