package lru

import "fmt"

type Node struct {
	Key   int
	Value int
	Pre   *Node
	Next  *Node
}

type LRUCache struct {
	Dict     map[int]*Node
	Root     *Node
	Tail     *Node
	Size     int
	Capacity int
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		Dict:     make(map[int]*Node),
		Root:     nil,
		Tail:     nil,
		Size:     0,
		Capacity: capacity,
	}
}

func (this *LRUCache) Get(key int) int {
	v, ok := this.Dict[key]
	if !ok {
		return -1
	}
	this.moveHead(v)
	return v.Value
}

func (this *LRUCache) Put(key int, value int) {
	v, ok := this.Dict[key]
	if !ok {
		if this.Size >= this.Capacity {
			n := this.delLast()
			if n != nil {
				delete(this.Dict, n.Key)
				this.Size -= 1
			}
		}
		node := &Node{
			Key:   key,
			Value: value,
			Pre:   nil,
			Next:  nil,
		}
		this.addHead(node)
		this.Dict[key] = node
		this.Size += 1
	} else {
		v.Value = value
		this.moveHead(v)
	}
}

func (this *LRUCache) addHead(node *Node) {
	if this.Root != nil {
		this.Root.Pre = node
	}
	node.Next = this.Root
	this.Root = node
	if this.Tail == nil {
		this.Tail = node
	}
}

func (this *LRUCache) moveHead(node *Node) {
	pre := node.Pre
	next := node.Next
	if pre == nil {
		return
	}
	pre.Next = next
	if next == nil {
		this.Tail = pre
	} else {
		next.Pre = pre
	}

	node.Next = this.Root
	if this.Root != nil {
		this.Root.Pre = node
	}
	this.Root = node

}

func (this *LRUCache) delLast() *Node {
	if this.Tail == nil {
		return nil
	}
	result := this.Tail
	pre := result.Pre
	if pre == nil {
		this.Root = nil
		this.Tail = nil
		return result
	}
	pre.Next = nil
	this.Tail = pre
	return result
}

func (this *LRUCache) print() {
	fmt.Printf("%v\n", this.Dict)
	for n := this.Root; n != nil; n = n.Next {
		fmt.Printf("(%d, %d)->", n.Key, n.Value)
	}
	fmt.Println()
}

/**
 * Your LRUCache object will be instantiated and called as such:
 * obj := Constructor(capacity);
 * param_1 := obj.Get(key);
 * obj.Put(key,value);
 */
