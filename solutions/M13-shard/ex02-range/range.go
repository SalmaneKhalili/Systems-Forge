package main

// RangeShard returns the shard index owning key. boundaries is sorted ascending;
// the result is the smallest i with key <= boundaries[i], or len(boundaries) when
// the key exceeds every boundary.
func RangeShard(key string, boundaries []string) int {
	for i, b := range boundaries {
		if key <= b {
			return i
		}
	}
	return len(boundaries)
}