package link

type Node struct {
	Item interface{}
	Next *Node
}

func (p *Node) Len() int {
	count := 0
	n := p
	for n != nil {
		count++
		n = n.Next
	}

	return count
}

func (p *Node) Reverse() *Node {
	var pre, cur, next *Node
	cur = p

	for cur != nil {
		next = cur.Next
		cur.Next = pre
		pre = cur
		cur = next
	}

	return pre
}

func (p *Node) Print() []interface{} {
	lst := make([]interface{}, 0)
	n := p
	for n != nil {
		lst = append(lst, n.Item)
		n = n.Next
	}
	return lst
}
