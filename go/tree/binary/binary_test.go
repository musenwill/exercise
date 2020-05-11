package binary

import (
	"testing"
)

func createTree() *node {
	a := &node{"a", nil, nil}
	b := &node{"b", nil, nil}
	c := &node{"c", nil, nil}
	d := &node{"d", nil, nil}
	e := &node{"e", nil, nil}
	f := &node{"f", nil, nil}
	g := &node{"g", nil, nil}
	h := &node{"h", nil, nil}
	k := &node{"k", nil, nil}

	root := a
	a.left = b
	a.right = c
	b.left = d
	b.right = e
	e.left = f
	e.right = g
	c.left = h
	c.right = k

	return root
}

func TestPrint(t *testing.T) {
	root := createTree()
	preOrderPrint := "abdefgchk"
	inOrderPrint := "dbfegahck"
	postOrderPrint := "dfgebhkca"

	exp, act := preOrderPrint, root.preOrderPrint()
	if exp != act {
		t.Errorf("preorder print got %v exptect %v", act, exp)
	}

	exp, act = inOrderPrint, root.inOrderPrint()
	if exp != act {
		t.Errorf("inorder print got %v exptect %v", act, exp)
	}

	exp, act = postOrderPrint, root.postOrderPrint()
	if exp != act {
		t.Errorf("postorder print got %v exptect %v", act, exp)
	}
}

func TestCreateFromPreOrderAndInOrder(t *testing.T) {
	preOrderPrint := "abdefgchk"
	inOrderPrint := "dbfegahck"
	postOrderPrint := "dfgebhkca"

	root, err := createFromPreOrderAndInOrder(preOrderPrint, inOrderPrint)
	if root == nil || err != nil {
		t.Errorf("calculate postorder from preorder %v and inorder %v, got none result expect %v, %v", preOrderPrint, inOrderPrint, postOrderPrint, err)
	}

	if act, exp := root.postOrderPrint(), postOrderPrint; act != exp {
		t.Errorf("calculate postorder from preorder %v and inorder %v, got %v expect %v", preOrderPrint, inOrderPrint, act, exp)
	}
}
