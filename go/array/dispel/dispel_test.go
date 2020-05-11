package dispel

import (
	"testing"
)

type caze struct {
	lst []int
	exp int64
}

func TestDispel(t *testing.T) {
	cases := []caze{
		{[]int{2, 1, 3}, 2},
		{[]int{0}, 0},
		{[]int{1}, 0},
		{[]int{9, 8, 7, 6, 5, 4, 3, 2, 1}, 20},
		{[]int{1, 3, 2, 0, 3}, 3},
	}

	for _, v := range cases {
		act := minimumValueAfterDispel(v.lst)
		if act != v.exp {
			t.Errorf("got %v expected %v of %v", act, v.exp, v.lst)
		}
	}
}
