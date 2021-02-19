package lock

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// BenchmarkNoLock-8   	1000000000	         0.549 ns/op	       0 B/op	       0 allocs/op
// BenchmarkNoLock-8   	1000000000	         0.542 ns/op	       0 B/op	       0 allocs/op
// BenchmarkNoLock-8   	1000000000	         0.545 ns/op	       0 B/op	       0 allocs/op
func BenchmarkNoLock(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		func() {
		}()
	}
}

// BenchmarkRLock-8   	74336402	        14.0 ns/op	       0 B/op	       0 allocs/op
// BenchmarkRLock-8   	84837378	        14.2 ns/op	       0 B/op	       0 allocs/op
// BenchmarkRLock-8   	77300755	        14.1 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRLock(b *testing.B) {
	var lock sync.RWMutex
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		func() {
			lock.RLock()
			defer lock.RUnlock()
		}()
	}
}

// BenchmarkLock-8   	84453649	        14.3 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLock-8   	84435949	        14.0 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLock-8   	84719203	        14.0 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLock(b *testing.B) {
	var lock sync.Mutex
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		func() {
			lock.Lock()
			defer lock.Unlock()
		}()
	}
}

func TestReRLock(t *testing.T) {
	var lock sync.RWMutex

	lock.RLock()
	go func() {
		lock.Lock()
		lock.Unlock()
	}()
	time.Sleep(time.Second)
	lock.RLock()
	lock.RUnlock()
	lock.RUnlock()
}

func TestWorker(t *testing.T) {
	resC := make(chan int)
	var n int
	for i := 0; i < 10; i++ {
		n++
		go func(i int) {
			resC <- i
		}(i)
	}

	for i := 0; i < n; i++ {
		fmt.Println(<-resC)
	}
	close(resC)
}

func TestClosedChan(t *testing.T) {
	closing := make(chan struct{})
	close(closing)
	select {
	case <-closing:
		t.Fatal()
	default:
	}
}
