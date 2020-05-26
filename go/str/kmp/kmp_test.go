package kmp

import (
	"strconv"
	"testing"
)

func TestKmp(t *testing.T) {
	cases := [][]string{
		{"", "", ""},
		{"uusd732y4", "", ""},
		{"", "d7s6yt", ""},
		{"bbc abcdab abcdabcdabde", "abcdabd", "15 22;"},
	}

	for _, c := range cases {
		if act, exp := printArray(Kmp(c[0], c[1])), c[2]; act != exp {
			t.Errorf("kmp compare between %s and %s got %s expect %s", c[0], c[1], act, exp)
		}
	}
}

func printArray(array [][]int) string {
	var r string
	for _, u := range array {
		for _, v := range u {
			r += strconv.Itoa(v) + " "
		}
		r = r[:len(r)-1]
		r += ";"
	}
	return r
}
