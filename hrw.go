// Package hrw provides an implementation of Highest Random Weight hashing, an
// alternative to consistent hashing which is both simple and fast.
//
// HRW allows you to consistently select the same nodes (or sets of nodes) for a
// given key. When a node is removed from the set of available nodes, only the
// data is was responsible for is at all affected.
//
// For more details on HRW hashing, see
// http://www.eecs.umich.edu/techreports/cse/96/CSE-TR-316-96.pdf or
// http://en.wikipedia.org/wiki/Rendezvous_hashing.
package hrw

import (
	"hash/fnv"
	"sort"
)

// SortByWeight returns the given set of nodes sorted in decreasing order of
// their weight for the given key.
func SortByWeight(nodes []int, key []byte) []int {
	h := fnv.New32a()
	h.Write(key)
	d := int32(h.Sum32())

	entries := make(entryList, len(nodes))

	for i, node := range nodes {
		entries[i] = entry{node: node, weight: weight(int32(node), d)}
	}

	sort.Sort(entries)

	sorted := make([]int, len(entries))
	for i, e := range entries {
		sorted[i] = e.node
	}
	return sorted
}

// TopN returns the top N nodes in decreasing order of their weight for the
// given key.
func TopN(nodes []int, key []byte, n int) []int {
	// BUG(coda): TopN is not optimized.
	return SortByWeight(nodes, key)[:n]
}

func weight(s, d int32) int {
	v := (a * ((a*s + c) ^ d + c))
	if v < 0 {
		v += m
	}
	return int(v)
}

type entry struct {
	node   int
	weight int
}

type entryList []entry

func (l entryList) Len() int {
	return len(l)
}

func (l entryList) Less(a, b int) bool {
	return l[a].weight > l[b].weight
}

func (l entryList) Swap(a, b int) {
	l[a], l[b] = l[b], l[a]
}

const (
	a = 1103515245    // multiplier
	c = 12345         // increment
	m = (1 << 31) - 1 // modulus (2**32-1)
)
