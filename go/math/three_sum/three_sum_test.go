package three_sum

import (
	"testing"
)

type Case struct {
	lst    []int
	result [][]int
}

func TestThreeSum(t *testing.T) {
	cases := []Case{
		{
			[]int{-1, 0, 1, 2, -1, -4},
			[][]int{
				{-1, -1, 2},
				{-1, 0, 1},
			},
		},
	}

	for _, v := range cases {
		if act, exp := toString(ThreeSum(v.lst)), toString(v.result); act != exp {
			t.Errorf("three sum of %v got %v expect %v", v.lst, ThreeSum(v.lst), v.result)
		}
	}
}

func toString(array [][]int) string {
	str := ""

	for _, v := range array {
		str += List(v).String() + ";"
	}
	return str
}
