package main

import (
	"fmt"
	"time"
)

func RunTicker(done chan struct{}, c chan <-time.Time, period time.Duration) {
	for {
		select {
		case c <- time.Now():
		case <-done:
			return
		}
		time.Sleep(period)
	}
}

func main() {
	done := make(chan struct{})
	heartbeats := make(chan time.Time)
	go RunTicker(done, heartbeats, 1 * time.Second)

	go func() {
		time.Sleep(5*time.Second)
		close(done)
	}()

	for tick := range heartbeats {
		fmt.Println("heartbeat:", tick)
	}

	// Équivalent 
	//for {
	//	tick, ok := <-heartbeats
	//	if !ok {
	//		break
	//	}
	//	...
	//}
}
