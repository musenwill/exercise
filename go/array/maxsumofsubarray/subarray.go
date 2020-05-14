package maxsumofsubarray

func FindGreatestSumOfSubArray(array []int) int {
	maxSum, sum := 0, 0

	for _, v := range array {
		sum += v
		if sum > maxSum {
			maxSum = sum
		}
		if sum < 0 {
			sum = 0
		}
	}

	return maxSum
}
