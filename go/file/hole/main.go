package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("usage: filename, hole size, zero size")
		return
	}
	filename := os.Args[1]
	holeSize, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	zeroSize, err := strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("filename: %s, hole size: %v, zero size: %v\n", filename, holeSize, zeroSize)

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	defer file.Sync()

	head := strings.Repeat("s", 512)
	tail := head
	_, err = file.Write([]byte(head))
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.Seek(holeSize, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.Write(make([]byte, zeroSize))
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.Write([]byte(tail))
	if err != nil {
		fmt.Println(err)
		return
	}
}
