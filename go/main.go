package main

import "fmt"

func main() {
	var f *foo = nil
	f.bar()
	fmt.Println(f)
}

type foo struct{}

func (f *foo) bar() {
	if f == nil {
		fmt.Println("nil.bar()")
	} else {
		fmt.Println("foo.bar()")
	}
}
