package power

import (
	"errors"
)

func powerRecursion(m, n int64) (int64, error) {
	if n < 0 {
		return 0, errors.New("negative n not allowed")
	}

	if n == 0 {
		return 1, nil
	}
	if m == 0 {
		return 0, nil
	}
	if n == 1 {
		return m, nil
	}

	if n&0x01 == 0x01 {
		t, err := powerRecursion(m*m, n/2)
		if err != nil {
			return 0, err
		} else {
			return t * m, nil
		}
	} else {
		return powerRecursion(m*m, n/2)
	}
}

func powerLoop(m, n int64) (int64, error) {
	if n < 0 {
		return 0, errors.New("negative n not allowed")
	}

	if n == 0 {
		return 1, nil
	}
	if m == 0 {
		return 0, nil
	}
	if n == 1 {
		return m, nil
	}

	var s int64 = 1
	for n > 0 {
		if n&0x01 == 0x01 {
			s *= m
		}
		m *= m
		n >>= 1
	}

	return s, nil
}
