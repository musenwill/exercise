package frequencycontrol

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// FrequencyLimitor caches result if requests too frequently to reduce system load
type FrequencyLimitor struct {
	Frequency  time.Duration // cache result if average of request interval less than it
	SampleTime time.Duration // average calculation section
	Name       string
	Logger     *zap.Logger

	mu         sync.Mutex
	cap        int                         // SampleTime / Frequency * 2
	records    []time.Time                 // cap of records = cap
	result     interface{}                 // cache result
	updateTime time.Time                   // result last update time
	fetcher    func() (interface{}, error) // call fetcher to get result
	closing    chan struct{}
	wg         *sync.WaitGroup
}

func NewFrequencyLimitor(name string, frequency, sampleTime time.Duration,
	fetcher func() (interface{}, error)) (*FrequencyLimitor, error) {
	if frequency <= 0 {
		return nil, fmt.Errorf("frequency expected to be positive, got %v", frequency)
	}
	if fetcher == nil {
		return nil, fmt.Errorf("unexpected nil fetcher")
	}

	cap := int(sampleTime / frequency * 2)

	return &FrequencyLimitor{
		Frequency:  frequency,
		SampleTime: sampleTime,
		Name:       name,
		Logger:     zap.NewNop(),
		cap:        cap,
		records:    make([]time.Time, 0, cap),
		fetcher:    fetcher,
		wg:         &sync.WaitGroup{},
		closing:    make(chan struct{}),
	}, nil
}

func (f *FrequencyLimitor) WithLogger(log *zap.Logger) {
	f.Logger = log.With(zap.String("service", f.Name))
}

func (f *FrequencyLimitor) Open() error {
	f.wg.Add(1)
	go func() {
		defer f.wg.Done()
		ticker := time.NewTicker(f.Frequency)
		defer ticker.Stop()

		for {
			select {
			case <-f.closing:
				return
			case <-ticker.C:
				result, err := f.mayFetch()
				if err != nil {
					f.Logger.Error("fetch failed", zap.Error(err))
					continue
				}
				f.result = result
			}
		}
	}()
	return nil
}

func (f *FrequencyLimitor) Close() {
	select {
	case <-f.closing:
	default:
		close(f.closing)
	}
	f.wg.Wait()
}

// Get result
func (f *FrequencyLimitor) Get() (interface{}, error) {
	select {
	case <-f.closing:
		return nil, fmt.Errorf("%s closed", f.Name)
	default:
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.records = append(f.records, time.Now())
	if len(f.records) > f.cap {
		f.records = f.records[len(f.records)-f.cap:]
	}

	if f.result != nil {
		return f.result, nil
	}

	result, err := f.fetcher()
	if err != nil {
		return nil, err
	}
	f.result = result
	f.updateTime = time.Now()
	return result, nil
}

// may fetch in lock
func (f *FrequencyLimitor) mayFetch() (interface{}, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// remove expired records
	index := 0
	for _, r := range f.records {
		if r.Before(time.Now().Add(-1 * f.SampleTime)) {
			index++
		} else {
			break
		}
	}
	f.records = f.records[index:]

	// do nothing if have no requests recently
	if len(f.records) == 0 {
		return nil, nil
	}

	// if has just fetch result then remain it
	if f.updateTime.Add(f.Frequency).After(time.Now()) {
		return f.result, nil
	}

	// calculate actual average frequency
	mean := time.Now().Sub(f.records[0]) / time.Duration(len(f.records))

	// do nothing if request frequency is low
	if mean > f.Frequency {
		return nil, nil
	}

	result, err := f.fetcher()
	if err != nil {
		return nil, err
	}
	f.result = result
	f.updateTime = time.Now()
	return result, nil
}
