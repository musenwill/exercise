package edit_distance

import (
	"testing"
)

type Case struct {
	StrA     string
	StrB     string
	Distance int
}

func TestDistance(t *testing.T) {
	cases := []Case{
		{"xyzab", "axyzc", 3},
		{"horse", "ros", 3},
		{"", "", 0},
		{"abcdewu", "", 7},
	}
	for _, c := range cases {
		if act, exp := editDistance(c.StrA, c.StrB), c.Distance; act != exp {
			t.Errorf("edit distance of %s and %s, got %d expect %d", c.StrA, c.StrB, act, exp)
		}
	}
}
