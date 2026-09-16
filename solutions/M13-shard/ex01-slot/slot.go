package main

// Hash returns the FNV-1a 32-bit hash of key.
func Hash(key string) uint32 {
	const (
		offset = uint32(2166136261)
		prime  = uint32(16777619)
	)
	h := offset
	for i := 0; i < len(key); i++ {
		h ^= uint32(key[i])
		h *= prime
	}
	return h
}

// SlotOf assigns key to a slot in [0, nSlots).
func SlotOf(key string, nSlots int) int {
	return int(Hash(key) % uint32(nSlots))
}