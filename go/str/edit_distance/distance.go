package edit_distance

func editDistance(a, b string) int {

	runeX, runeY := []rune(a), []rune(b)
	lenX, lenY := len(runeX), len(runeY)
	if lenY < lenX {
		runeX, runeY = runeY, runeX
		lenX, lenY = lenY, lenX
	}

	table := make([]int, lenY)
	for i := 0; i < len(table); i += 1 {
		table[i] = i + 1
	}

	for i := 0; i < lenX; i += 1 {
		upLeft := i
		left := i + 1
		for j := 0; j < lenY; j += 1 {
			up := table[j]

			var cur int
			if runeX[i] == runeY[j] {
				cur = upLeft
			} else {
				cur = min(left, upLeft, up) + 1
			}

			left = cur
			upLeft = up
			table[j] = cur
		}
	}

	if lenY > 0 {
		return table[len(table)-1]
	} else {
		return 0
	}
}

func min(a, b, c int) int {
	r := a
	if b < r {
		r = b
	}
	if c < r {
		r = c
	}
	return r
}
