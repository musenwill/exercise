package frequencycontrol_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/musenwill/exercise/performance/frequencycontrol"
	"github.com/stretchr/testify/require"
)

func TestFrequencyLimiter(t *testing.T) {
	t.Parallel()

	type tcase struct {
		frequency, sampleTime int   // millisecond
		beat                  []int // millisecond
		hit                   []int
	}
	runcase := func(c *tcase) {
		hit := make([]int, 0)
		var count int
		fl, err := frequencycontrol.NewFrequencyLimitor("frequency tester",
			time.Duration(c.frequency*10)*time.Millisecond,
			time.Duration(c.sampleTime*10)*time.Millisecond,
			func() (interface{}, error) {
				count++
				return count, nil
			})
		require.NoError(t, err)
		fl.Open()
		defer fl.Close()

		for _, b := range c.beat {
			if b > 0 {
				time.Sleep(time.Duration(b*10) * time.Millisecond)
			}
			result, err := fl.Get()
			require.NoError(t, err)
			hit = append(hit, result.(int))
		}
		if !reflect.DeepEqual(c.hit, hit) {
			t.Errorf("expect %v got %v with frequency %v and sample time %v, beats %v", c.hit, hit, c.frequency, c.sampleTime, c.beat)
		}
	}

	runcase(&tcase{10, 40, []int{0}, []int{1}})
	runcase(&tcase{10, 40, []int{5}, []int{1}})
	runcase(&tcase{10, 40, []int{15}, []int{1}})
	runcase(&tcase{10, 40, []int{25}, []int{1}})
	runcase(&tcase{10, 40, []int{35}, []int{1}})
	runcase(&tcase{10, 40, []int{45}, []int{1}})

	runcase(&tcase{10, 40, []int{0, 5}, []int{1, 1}})
	runcase(&tcase{10, 40, []int{0, 15}, []int{1, 2}})
	runcase(&tcase{10, 40, []int{0, 35}, []int{1, 2}})
	runcase(&tcase{10, 40, []int{0, 45}, []int{1, 2}})

	runcase(&tcase{10, 40, []int{0, 3, 3}, []int{1, 1, 1}})
	runcase(&tcase{10, 40, []int{0, 7, 7}, []int{1, 1, 2}})
	runcase(&tcase{10, 40, []int{0, 7, 7, 2}, []int{1, 1, 2, 2}})

	runcase(&tcase{30, 90, []int{0, 1, 1, 1, 1, 1, 1, 1, 1, 1}, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}})
	runcase(&tcase{10, 40, []int{0, 15, 20, 20}, []int{1, 2, 3, 4}})
}

func TestFrequencyLimiterLongRun(t *testing.T) {
	fl, err := frequencycontrol.NewFrequencyLimitor("frequency tester", time.Millisecond, 20*time.Millisecond,
		func() (interface{}, error) {
			return nil, nil
		})
	require.NoError(t, err)
	for i := 0; i < 1000; i++ {
		_, err := fl.Get()
		require.NoError(t, err)
		time.Sleep(30 * time.Microsecond)
	}
	for i := 0; i < 1000; i++ {
		_, err := fl.Get()
		require.NoError(t, err)
		time.Sleep(170 * time.Microsecond)
	}
}
