package count1

import (
	"testing"
)

func TestCount1(t *testing.T) {
	cases := [][]int{
		{0, 0},
		{1, 1},
		{2, 1},
		{10, 2},
		{13, 6},
		{20, 12},
		{99, 20},
		{101, 23},
		{211, 144},
		{234, 154},
	}

	for _, v := range cases {
		if act, exp := count1(v[0]), v[1]; act != exp {
			t.Errorf("count of 1 between 1 and %d got %d exp %d", v[0], act, exp)
		}
	}
}
