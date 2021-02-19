package main

import (
	"fmt"
	"time"
)

func main() {
	closing := make(chan bool)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	go func() {
		time.Sleep(time.Second)
		close(closing)
	}()

	for {
		select {
		case <-closing:
			return
		case <-ticker.C:
			fmt.Println("=========")
		}
	}
}
