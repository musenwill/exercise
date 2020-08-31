package typecast

import "testing"

type fruit interface {
	color() string
}

type apple struct{}

func (a *apple) color() string {
	return "red"
}

// BenchmarkNoTypeCast-8   	1000000000	         0.309 ns/op	       0 B/op	       0 allocs/op
// BenchmarkNoTypeCast-8   	1000000000	         0.285 ns/op	       0 B/op	       0 allocs/op
// BenchmarkNoTypeCast-8   	1000000000	         0.303 ns/op	       0 B/op	       0 allocs/op
func BenchmarkNoTypeCast(b *testing.B) {
	a := &apple{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.color()
	}
}

// BenchmarkTypeCast-8   	669254890	         1.94 ns/op	       0 B/op	       0 allocs/op
// BenchmarkTypeCast-8   	705645292	         1.77 ns/op	       0 B/op	       0 allocs/op
// BenchmarkTypeCast-8   	680690845	         1.84 ns/op	       0 B/op	       0 allocs/op
func BenchmarkAbstract(b *testing.B) {
	var a fruit = &apple{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.color()
	}
}
