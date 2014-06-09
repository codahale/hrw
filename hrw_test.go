package hrw

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
)

func Example() {
	// given a set of servers
	servers := map[int]string{
		1: "one.example.com",
		2: "two.example.com",
		3: "three.example.com",
		4: "four.example.com",
		5: "five.example.com",
		6: "six.example.com",
	}

	// which can be mapped to integer values
	ids := make([]int, 0, len(servers))
	for id := range servers {
		ids = append(ids, id)
	}

	// HRW can consistently select a uniformly-distributed set of servers for
	// any given key
	key := []byte("/examples/object-key")
	for _, id := range TopN(ids, key, 3) {
		fmt.Printf("trying GET %d %s%s\n", id, servers[id], key)
	}

	// Output:
	// trying GET 1 one.example.com/examples/object-key
	// trying GET 3 three.example.com/examples/object-key
	// trying GET 5 five.example.com/examples/object-key
}

func TestSortByWeight(t *testing.T) {
	key := []byte("hello, world")
	nodes := []int{1, 2, 3, 4, 5}

	actual := SortByWeight(nodes, key)
	expected := []int{5, 4, 2, 1, 3}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Was %#v, but expected %#v", actual, expected)
	}
}

func TestTopN(t *testing.T) {
	key := []byte("hello, world")
	nodes := []int{1, 2, 3, 4, 5}

	actual := TopN(nodes, key, 3)
	expected := []int{5, 4, 2}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Was %#v, but expected %#v", actual, expected)
	}
}

func TestUniformDistribution(t *testing.T) {
	t.Skip()
	nodes := []int{1, 2, 3, 4}
	counts := make(map[int]int)
	key := make([]byte, 16)
	keys := 1000000

	for i := 0; i < keys; i++ {
		binary.BigEndian.PutUint32(key, uint32(i))
		counts[SortByWeight(nodes, key)[0]]++
	}

	mean := float64(keys) / float64(len(nodes))
	delta := mean * 0.02 // 2%
	for node, count := range counts {
		d := mean - float64(count)
		if d > delta || (0-d) > delta {
			t.Errorf(
				"Node %d received %d keys, expected %v (+/- %v)",
				node, count, mean, delta,
			)
		}
	}
}

func BenchmarkSortByWeight10(b *testing.B) {
	_ = benchmarkSortByWeight(b, 10)
}

func BenchmarkSortByWeight100(b *testing.B) {
	_ = benchmarkSortByWeight(b, 100)
}

func BenchmarkSortByWeight1000(b *testing.B) {
	_ = benchmarkSortByWeight(b, 1000)
}

func benchmarkSortByWeight(b *testing.B, n int) int {
	key := []byte("hello, world")
	servers := make([]int, n)
	for i := 0; i < len(servers); i++ {
		servers[i] = i
	}
	b.ResetTimer()

	var x int
	for i := 0; i < b.N; i++ {
		x += SortByWeight(servers, key)[0]
	}
	return x
}
