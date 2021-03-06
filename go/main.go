package main

import (
	"fmt"
	"time"
)

func main() {
	var timerCh <-chan time.Time
	timerChrw := make(chan time.Time)
	close(timerChrw)
	timerCh = timerChrw

	for {
		select {
		case <-timerCh:
			fmt.Println("======")
		}
	}
}
