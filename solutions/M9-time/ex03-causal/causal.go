package main

// CausallyBefore reports that a happened before b: strictly componentwise <=.
func CausallyBefore(a, b [3]int) bool {
	strict := false
	for i := range a {
		if a[i] > b[i] {
			return false
		}
		if a[i] < b[i] {
			strict = true
		}
	}
	return strict
}

// Concurrent reports that neither a happened before b nor b before a.
func Concurrent(a, b [3]int) bool {
	return !CausallyBefore(a, b) && !CausallyBefore(b, a)
}