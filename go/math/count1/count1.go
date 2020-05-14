package count1

func count1(n int) int {
	if n <= 0 {
		return 0
	}

	count := 0
	for m := 1; m <= n; m *= 10 {
		a, b := n/m, n%m
		count += (a + 8) / 10 * m
		if a%10 == 1 {
			count += b + 1
		}
	}

	return count
}
