package sync

import (
	"sync"
	"sync/atomic"
)

type BufPool struct {
	maxSize int
	pool    *sync.Pool
	count   int64
	hit     int64
}

// new limited buf pool. catch buf smaller than maxBufSize.
func NewBufPool(maxBufSize int) *BufPool {
	return &BufPool{
		maxSize: maxBufSize,
		pool:    new(sync.Pool),
	}
}

// Get returns a buffer with length size from the buffer pool.
func (p *BufPool) Get(size int) []byte {
	atomic.AddInt64(&p.count, 1)

	if size > p.maxSize {
		return make([]byte, size)
	}

	x := p.pool.Get()
	if x == nil {
		return make([]byte, size)
	}
	buf := x.([]byte)

	if cap(buf) < size {
		// drop the short

		return make([]byte, size)
	}

	atomic.AddInt64(&p.hit, 1)
	buf = buf[:size]
	return buf
}

// putBuf returns a buffer to the pool.
func (p *BufPool) Put(buf []byte) {
	if cap(buf) > p.maxSize {
		return
	}

	p.pool.Put(buf)
}

func (p *BufPool) Statistics() (count, hit int64) {
	return p.count, p.hit
}
