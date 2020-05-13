package stack

type Stack struct {
	lst []interface{}
}

func CreateStack(lst []interface{}) *Stack {
	return &Stack{lst: lst}
}

func (p *Stack) Len() int {
	return len(p.lst)
}

func (p *Stack) Empty() bool {
	return p.Len() == 0
}

func (p *Stack) Peek() interface{} {
	if p.Empty() {
		return nil
	}
	return p.lst[p.Len()-1]
}

func (p *Stack) Pop() interface{} {
	if p.Empty() {
		return nil
	}

	result := p.lst[p.Len()-1]
	p.lst = p.lst[:p.Len()-1]

	return result
}

func (p *Stack) Push(e interface{}) {
	p.lst = append(p.lst, e)
}

func (p *Stack) Reverse() {
	head := 0
	tail := p.Len() - 1

	for head < tail {
		p.lst[head], p.lst[tail] = p.lst[tail], p.lst[head]
		head++
		tail--
	}
}

func (p *Stack) Iter(f func(index int, e interface{})) {
	lst := make([]interface{}, 0)
	copy(lst, p.lst)

	for i, v := range lst {
		f(i, v)
	}
}
