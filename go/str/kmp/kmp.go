package kmp

func Kmp(str1, str2 string) [][]int {
	result := make([][]int, 0)

	if len(str2) == 0 {
		return result
	}

	runesA := []rune(str1)
	runesB := []rune(str2)

	state := make([]int, len(runesB))
	k := 0
	for i := 1; i < len(runesB); i += 1 {
		for k > 0 && runesB[i] != runesB[k] {
			k = state[k-1]
		}
		if runesB[i] == runesB[k] {
			k += 1
		}
		state[i] = k
	}

	k = 0
	for i := 0; i < len(runesA); i += 1 {
		for k > 0 && runesA[i] != runesB[k] {
			k = state[k-1]
		}
		if runesA[i] == runesB[k] {
			k += 1
		}
		if k >= len(state) {
			result = append(result, []int{i - k + 1, i + 1})
			k = state[k-1]
		}
	}
	return result
}
