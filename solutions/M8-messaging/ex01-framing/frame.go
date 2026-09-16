package main

import (
	"encoding/binary"
	"io"
)

// PutFrame writes one length-prefixed frame to w: a 4-byte big-endian length
// followed by the payload. It loops so short writes cannot split a frame.
func PutFrame(w io.Writer, msg []byte) error {
	var hdr [4]byte
	binary.BigEndian.PutUint32(hdr[:], uint32(len(msg)))
	if _, err := writeFull(w, hdr[:]); err != nil {
		return err
	}
	_, err := writeFull(w, msg)
	return err
}

// writeFull loops over w.Write until everything is flushed or an error stops us.
func writeFull(w io.Writer, b []byte) (int, error) {
	total := 0
	for total < len(b) {
		n, err := w.Write(b[total:])
		total += n
		if err != nil {
			return total, err
		}
		if n == 0 {
			return total, io.ErrShortWrite
		}
	}
	return total, nil
}

// GetFrame reads one frame back. io.ReadFull loops internally, so a reader
// that trickles bytes in is no threat. EOF at a clean boundary comes back as
// io.EOF; EOF halfway through a frame comes back as io.ErrUnexpectedEOF.
func GetFrame(r io.Reader) ([]byte, error) {
	var hdr [4]byte
	if _, err := io.ReadFull(r, hdr[:]); err != nil {
		return nil, err
	}
	n := binary.BigEndian.Uint32(hdr[:])
	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}