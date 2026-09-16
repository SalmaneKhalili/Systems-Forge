package main

// VC is a vector clock for two processes (indices 0 and 1).
type VC struct {
	v [2]int
}

// Local ticks component i for a local event and returns the new vector.
func (vc *VC) Local(i int) [2]int {
	vc.v[i]++
	return vc.v
}

// Merge folds a received vector in elementwise before the receive ticks.
func (vc *VC) Merge(other [2]int) {
	for i := range vc.v {
		if other[i] > vc.v[i] {
			vc.v[i] = other[i]
		}
	}
}

// Now returns the current vector without changing it.
func (vc *VC) Now() [2]int {
	return vc.v
}