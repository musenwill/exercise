package main

import "fmt"

func main() {

	for i := 0; i < 5; i++ {
		loop := i
		defer func() {
			fmt.Printf("loop %d\n", loop)
		}()
	}
}
