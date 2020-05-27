package frog_pass_river

func CanCross(stones []int) bool {
	length := len(stones)
	if length < 2 {
		return true
	}
	steps := make([]map[int]bool, len(stones)) // 保存每一块石子都需要哪些步长才能跨越至它
	steps[0] = map[int]bool{0: true}           // 初始为 0

	for i, v := range steps { // 广度优先遍历
		if len(v) == 0 { // 如果某一块石头没有步长可以到达它，说明问题无解，直接返回
			return false
		}
		distance := stones[i]
		for u := range v { // 遍历该石头保存的所有步长
			s := []int{u - 1, u, u + 1}
			for _, k := range s {
				if k > 0 { // 步子大于零才有意义
					index := binarySearch(distance+k, i+1, length-1, stones)
					if index > 0 {
						if steps[index] == nil {
							steps[index] = make(map[int]bool)
						}
						steps[index][k] = true
						if index >= length-1 {
							return true
						}
					}
				}
			}
		}

	}
	return false
}

func binarySearch(value, left, right int, stones []int) int {
	for left <= right {
		middle := (left + right) / 2
		if stones[middle] == value {
			return middle
		} else if stones[middle] > value {
			right = middle - 1
		} else {
			left = middle + 1
		}
	}
	return -1
}
