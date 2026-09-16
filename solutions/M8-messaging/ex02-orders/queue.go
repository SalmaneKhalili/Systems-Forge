package main

import (
	"errors"
)

type queue struct {
	arrived   []entry
	delivered []bool
	next      map[int]int
}

func NewQueue() *queue {
	return &queue{next: map[int]int{}}
}

// Put records one arrival: an entry. Ordering is decided at Drain time, so
// Put is free to append in arrival order.
func (q *queue) Put(partition, offset int) {
	q.arrived = append(q.arrived, entry{partition, offset})
	q.delivered = append(q.delivered, false)
}

// Drain emits the allowed delivery order:
//   - per-partition FIFO — an entry goes only when its predecessor went;
//   - gap holding — an entry whose predecessor has not arrived is skipped;
//   - earliest arrival wins — deliverable entries are swept from the head of
//     the arrival list.
// Repeatedly restarting the sweep from index 0 enforces "earliest arrival
// first": at every step the head-of-line deliverable entry is the one taken.
// If the sweep ends with entries still held by an unrecoverable gap, Drain
// reports the error instead of silently reordering.
func (q *queue) Drain() ([]entry, error) {
	var out []entry
	for {
		take := -1
		for i := range q.arrived {
			if q.delivered[i] {
				continue
			}
			e := q.arrived[i]
			if e.offset == q.next[e.partition] {
				take = i
				break
			}
		}
		if take < 0 {
			for _, done := range q.delivered {
				if !done {
					return nil, errors.New("per-partition gap")
				}
			}
			return out, nil
		}
		e := q.arrived[take]
		q.delivered[take] = true
		q.next[e.partition]++
		out = append(out, e)
	}
}