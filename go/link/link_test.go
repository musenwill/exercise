package link

import (
	"testing"
)

func TestLink(t *testing.T) {
	na := &Node{"a", nil}
	nb := &Node{"b", nil}
	nc := &Node{"c", nil}
	nd := &Node{"d", nil}
	ne := &Node{"e", nil}

	lk1 := na
	if act, exp := lk1.Len(), 1; act != exp {
		t.Errorf("link length got %d expect %d", act, exp)
	}
	if act, exp := lk1.Print(), []interface{}{"a"}; !isSame(act, exp) {
		t.Errorf("print link got %v expect %v", act, exp)
	}
	if act, exp := lk1.Reverse().Print(), []interface{}{"a"}; !isSame(act, exp) {
		t.Errorf("print reverse link got %v expect %v", act, exp)
	}

	lk2 := na
	na.Next = nb
	nb.Next = nc
	nc.Next = nd
	nd.Next = ne
	if act, exp := lk2.Len(), 5; act != exp {
		t.Errorf("link length got %d expect %d", act, exp)
	}
	if act, exp := lk2.Print(), []interface{}{"a", "b", "c", "d", "e"}; !isSame(act, exp) {
		t.Errorf("print link got %v expect %v", act, exp)
	}
	if act, exp := lk2.Reverse().Print(), []interface{}{"e", "d", "c", "b", "a"}; !isSame(act, exp) {
		t.Errorf("print reverse link got %v expect %v", act, exp)
	}
}

func isSame(lst1, lst2 []interface{}) bool {
	if len(lst1) != len(lst2) {
		return false
	}

	for i := 0; i < len(lst1); i++ {
		if lst1[i] != lst2[i] {
			return false
		}
	}

	return true
}
