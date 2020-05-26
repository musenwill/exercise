package lru

import (
	"testing"
)

func TestLRU(t *testing.T) {
	lru := Constructor(2)
	lru.Put(1, 1)
	lru.Put(2, 2)

	if act, exp := lru.Get(1), 1; act != exp {
		t.Errorf("get %d from lru got %d expect %d", 1, act, exp)
	}
	lru.Put(3, 3)
	if act, exp := lru.Get(2), -1; act != exp {
		t.Errorf("get %d from lru got %d expect %d", 2, act, exp)
	}
}
