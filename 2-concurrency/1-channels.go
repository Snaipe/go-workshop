package main

import (
	"fmt"
	"time"
)

func RunTicker(c chan <-time.Time, period time.Duration) {
	for {
		c <- time.Now()
		time.Sleep(period)
	}
}

func main() {
	heartbeats := make(chan time.Time)
	go RunTicker(heartbeats, 1 * time.Second)

	tick, ok := <-heartbeats

	for tick := range heartbeats {
		fmt.Println("heartbeat:", tick)
	}
}
