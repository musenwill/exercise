package job

import (
	"fmt"
	"sync"

	"go.uber.org/atomic"
)

type Job struct {
	pool    chan int
	dealing atomic.Bool
	wg      sync.WaitGroup
}

func NewJob() *Job {
	return &Job{
		pool: make(chan int, 10),
	}
}

func (j *Job) write(i int) {
	j.pool <- i
	j.startJob()
}

func (j *Job) startJob() {
	j.wg.Add(1)
	go func() {
		defer j.wg.Done()

		if !j.dealing.CAS(false, true) {
			return
		}
		defer j.dealing.Store(false)

		for len(j.pool) > 0 {
			fmt.Println(<-j.pool)
		}
	}()
}

func (j *Job) Close() {
	j.wg.Wait()
}
