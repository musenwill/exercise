package lock

import (
	"errors"
	"sync"

	"go.uber.org/atomic"
)

type PartialLock struct {
	pre        sync.Mutex
	post       sync.Mutex
	prelocked  atomic.Bool
	postlocked atomic.Bool
}

func NewPartialLock() *PartialLock {
	return &PartialLock{}
}

func (p *PartialLock) Prelock() error {
	if p.prelocked.Load() {
		return errors.New("detect relock")
	}
	if p.postlocked.Load() {
		return errors.New("detect reverse lock")
	}

	p.pre.Lock()
	p.prelocked.Store(true)

	return nil
}

func (p *PartialLock) Preunlock() error {
	if !p.prelocked.Load() {
		return errors.New("unlock not locked lock")
	}

	p.prelocked.Store(false)
	p.pre.Unlock()
	return nil
}

func (p *PartialLock) Postlock() error {
	if p.postlocked.Load() {
		return errors.New("detect relock")
	}

	p.post.Lock()
	p.postlocked.Store(true)

	return nil
}

func (p *PartialLock) Postunlock() error {
	if !p.postlocked.Load() {
		return errors.New("unlock not locked lock")
	}

	p.postlocked.Store(false)
	p.post.Unlock()
	return nil
}
