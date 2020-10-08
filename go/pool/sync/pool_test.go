package sync

import (
	"testing"
)

func TestLevelMapPool(t *testing.T) {
	for i := 0; i < 10000; i++ {
		mm := getL2Map()

		pet := getMap()
		pet["cat"] = struct{}{}
		pet["dog"] = struct{}{}
		mm["pet"] = pet

		food := getMap()
		food["egg"] = struct{}{}
		food["bread"] = struct{}{}
		mm["food"] = food

		putL2Map(mm)

		mm = getL2Map()
		if _, exist := mm["pet"]; exist {
			t.Fatal()
		}
		if _, exist := mm["food"]; exist {
			t.Fatal()
		}
	}
}

// BenchmarkLevelMapPool-8   	 5097276	       233 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLevelMapPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mm := getL2Map()

		pet := getMap()
		pet["cat"] = struct{}{}
		pet["dog"] = struct{}{}
		mm["pet"] = pet

		food := getMap()
		food["egg"] = struct{}{}
		food["bread"] = struct{}{}
		mm["food"] = food

		putL2Map(mm)
	}
}

// BenchmarkMap-8   	 3515077	       338 ns/op	     672 B/op	       4 allocs/op
func BenchmarkMap(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mm := make(map[string]interface{})

		pet := make(map[string]interface{})
		pet["cat"] = struct{}{}
		pet["dog"] = struct{}{}
		mm["pet"] = pet

		food := make(map[string]interface{})
		food["egg"] = struct{}{}
		food["bread"] = struct{}{}
		mm["food"] = food
	}
}
