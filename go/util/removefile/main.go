package main

import (
	"fmt"
	"os"
)

func main() {
	if err := os.Remove("noexist"); err != nil {
		fmt.Println(err)
	}
	if err := os.Remove("noexist"); err != nil {
		fmt.Println(err)
	}
}
