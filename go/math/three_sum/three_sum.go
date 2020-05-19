package three_sum

import (
	"sort"
)

func ThreeSum(nums []int) [][]int {
	lst := List(nums)
	if lst.Len() < 3 {
		return nil
	}
	sort.Sort(lst)
	result := make(map[string][]int)

	l1, l2 := 0, 1
	r := lst.Len() - 1
	for l2 < r {
		t := lst[l1] + lst[l2] + lst[r]
		if t == 0 {
			answer := []int{lst[l1], lst[l2], lst[r]}
			result[List(answer).String()] = answer
			r -= 1
		} else if t < 0 {
			l1 += 1
			l2 += 1
		} else {
			r -= 1
		}
	}

	l := 0
	r1 := lst.Len() - 1
	r2 := r1 - 1
	for l < r2 {
		t := lst[l] + lst[r1] + lst[r2]
		if t == 0 {
			answer := []int{lst[l], lst[r2], lst[r1]}
			result[List(answer).String()] = answer
			l += 1
		} else if t < 0 {
			l += 1
		} else {
			r1 -= 1
			r2 -= 1
		}
	}

	ret := make([][]int, 0)
	for _, v := range result {
		ret = append(ret, v)
	}

	return ret
}

type List []int

func (p List) Len() int {
	return len(p)
}

func (p List) Less(i, j int) bool {
	return p[i] < p[j]
}

func (p List) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p List) String() string {
	str := ""
	for _, v := range p {
		str += string(v) + " "
	}
	return str
}
