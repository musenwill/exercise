package power

import (
	"testing"
)

func TestPowerRecursion(t *testing.T) {
	// testPower(t, powerRecursion)
	testPower(t, powerLoop)
}

func testPower(t *testing.T, f func(m, n int64) (int64, error)) {
	cases := [][3]int64{
		{0, 1, 0},
		{0, 100, 0},
		{0, 101, 0},
		{0, 0, 1},

		{1, 0, 1},
		{-100, 0, 1},
		{100, 0, 1},

		{1, 100, 1},
		{1, 101, 1},
		{-1, 100, 1},
		{-1, 101, -1},

		{2, 1, 2},
		{2, 2, 4},
		{2, 5, 32},
		{2, 10, 1024},
		{2, 31, 2147483648},
		{-2, 32, 4294967296},
		{-2, 33, -8589934592},
	}

	for _, v := range cases {
		m, n, exp := v[0], v[1], v[2]
		act, err := f(m, n)
		if exp != act || err != nil {
			t.Errorf("%v^%v got %v expected %v, %v", m, n, act, exp, err)
		}
	}
}
