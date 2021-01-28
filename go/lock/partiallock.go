package lock

import (
	"errors"
	"sync"

	"go.uber.org/atomic"
)

type PartialLock struct {
	mu     sync.Mutex
	locked atomic.Bool
	postmu *PartialLock
	bind   atomic.Bool
}

func NewPartialLock() *PartialLock {
	return &PartialLock{}
}

func (p *PartialLock) WithPostLock(post *PartialLock) error {
	if err := post.Bind(); err != nil {
		return err
	}
	p.postmu = post
	return nil
}

func (p *PartialLock) Locked() bool {
	return p.locked.Load()
}

func (p *PartialLock) Bind() error {
	if p.bind.Load() {
		return errors.New("post has been bind")
	}
	p.bind.Store(true)
	return nil
}

func (p *PartialLock) Lock() error {
	if p.locked.Load() {
		return errors.New("detect relock")
	}
	if p.postmu != nil && p.postmu.Locked() {
		return errors.New("detect reverse lock")
	}

	p.mu.Lock()
	p.locked.Store(true)

	return nil
}

func (p *PartialLock) Unlock() error {
	if !p.locked.Load() {
		return errors.New("unlock not locked lock")
	}

	p.locked.Store(false)
	p.mu.Unlock()
	return nil
}
