package binary

import (
	"errors"
)

type node struct {
	name        string
	left, right *node
}

func (p *node) preOrderPrint() string {
	var result string = p.name
	if p.left != nil {
		result = result + p.left.preOrderPrint()
	}
	if p.right != nil {
		result = result + p.right.preOrderPrint()
	}

	return result
}

func (p *node) inOrderPrint() string {
	var result string = p.name
	if p.left != nil {
		result = p.left.inOrderPrint() + result
	}
	if p.right != nil {
		result = result + p.right.inOrderPrint()
	}

	return result
}

func (p *node) postOrderPrint() string {
	var result string = ""
	if p.left != nil {
		result = result + p.left.postOrderPrint()
	}
	if p.right != nil {
		result = result + p.right.postOrderPrint()
	}

	return result + p.name
}

func createFromPreOrderAndInOrder(preOrder, inOrder string) (*node, error) {
	if len(preOrder) <= 0 {
		return nil, nil
	}
	if len(preOrder) != len(inOrder) {
		return nil, errors.New("invalid input")
	}

	length := len(preOrder)

	root := &node{
		preOrder[0:1], nil, nil,
	}

	i := 0
	for ; i < length; i++ {
		if inOrder[i:i+1] == root.name {
			break
		}
	}
	if i >= length {
		return nil, errors.New("invalid input")
	}

	if length > 1 {
		preOrderLeft := preOrder[1 : 1+i]
		inOrderLeft := inOrder[0:i]
		left, err := createFromPreOrderAndInOrder(preOrderLeft, inOrderLeft)
		if err != nil {
			return nil, err
		}
		root.left = left
	}

	if 1+i < length {
		preOrderRight := preOrder[1+i:]
		inOrderRight := inOrder[1+i:]
		right, err := createFromPreOrderAndInOrder(preOrderRight, inOrderRight)
		if err != nil {
			return nil, err
		}
		root.right = right
	}

	return root, nil
}
