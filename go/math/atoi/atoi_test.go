package atoi

import (
	"testing"
)

type Case struct {
	str   string
	value int
}

func TestAtoi(t *testing.T) {
	cases := []Case{
		{"0", 0},
		{"2", 2},
		{"-1", -1},
		{"42", 42},
		{"  321", 321},
		{"- 3908", 0},
		{"  +0001", 1},
		{"   -00000000100", -100},
		{" -00919-", -919},
		{"2147483647", 2147483647},
		{"2147483648", 2147483647},
		{"-2147483648", -2147483648},
		{"-2147483649", -2147483648},
		{"214748364776436364565443", 2147483647},
		{"-21474838362555444649", -2147483648},
	}

	for _, c := range cases {
		if act, exp := atoi(c.str), c.value; act != exp {
			t.Errorf("%s atoi got %d expect %d", c.str, act, exp)
		}
	}
}
