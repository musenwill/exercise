package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Second
	tick := time.NewTicker(t)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			fmt.Println(t)
			t *= 2
			if t > 10*time.Second {
				t = 10 * time.Second
			}
			tick.Reset(t)
		}
	}
}
