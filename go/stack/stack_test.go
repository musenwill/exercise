package stack

import "testing"

func TestStack(t *testing.T) {
	stack := CreateStack(make([]interface{}, 0))
	var none interface{} = nil

	if act, exp := stack.Len(), 0; act != exp {
		t.Errorf("stack length got %v expect %v", act, exp)
	}
	if act, exp := stack.Empty(), true; act != exp {
		t.Errorf("stack empty got %v expect %v", act, exp)
	}
	if act, exp := stack.Peek(), none; act != exp {
		t.Errorf("stack peek got %v expect %v", act, exp)
	}
	if act, exp := stack.Pop(), none; act != exp {
		t.Errorf("stack pop got %v expect %v", act, exp)
	}

	stack.Push("a")
	stack.Push("b")
	stack.Push("c")
	stack.Push("d")

	if act, exp := stack.Len(), 4; act != exp {
		t.Errorf("stack length got %v expect %v", act, exp)
	}
	if act, exp := stack.Empty(), false; act != exp {
		t.Errorf("stack empty got %v expect %v", act, exp)
	}
	if act, exp := stack.Peek().(string), "d"; act != exp {
		t.Errorf("stack peek got %v expect %v", act, exp)
	}
	if act, exp := stack.Pop().(string), "d"; act != exp {
		t.Errorf("stack pop got %v expect %v", act, exp)
	}
}
