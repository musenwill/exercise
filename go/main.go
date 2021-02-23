package main

import "fmt"

func main() {
	buf := make([]byte, 10)
	fmt.Println(cap(buf))
	buf = buf[:5]
	fmt.Println(cap(buf))
}
