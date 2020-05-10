package dispel

import (
	"sort"
)

type list []int

func (p list) Len() int {
	return len(p)
}

func (p list) Less(i, j int) bool {
	return p[i] < p[j]
}

func (p list) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p list) sum() int64 {
	var s int64 = 0

	for _, v := range p {
		s += int64(v)
	}

	return s
}

func minimumValueAfterDispel(nums []int) int64 {
	lst := list(nums)
	sort.Sort(lst)

	sum := lst.sum()

	var maxDelta int64 = 0
	for i := lst.Len() - 1; i >= 0; i-- {
		delta := int64(lst[i] * (lst.Len() - i))
		if delta > maxDelta {
			maxDelta = delta
		}
	}

	return sum - maxDelta
}
