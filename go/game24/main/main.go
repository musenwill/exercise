package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/musenwill/exercise/game24"
)

func main() {
	nums, err := args()
	if err != nil {
		fmt.Println(err)
		return
	}
	game24.Calc(float64(nums[0]), float64(nums[1]), float64(nums[2]), float64(nums[3]), 24)
}

func args() ([]int, error) {
	if len(os.Args) != 5 {
		return nil, fmt.Errorf("expect 4 numbers")
	}

	nums := make([]int, 0, 4)

	for i := 1; i < len(os.Args); i++ {
		num, err := strconv.Atoi(os.Args[i])
		if err != nil {
			return nil, err
		}
		nums = append(nums, num)
	}

	return nums, nil
}
