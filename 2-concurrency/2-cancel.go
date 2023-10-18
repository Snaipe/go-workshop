package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func RunTicker(ctx context.Context, c chan<- time.Time, period time.Duration) {
	defer close(c)

	timer := time.NewTimer(period)
	defer timer.Stop()

	defer func() {
		fmt.Println("exiting ticker:", ctx.Err())
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case c <- time.Now():
		}

		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}
		timer.Reset(period)
	}
}

func main() {
	ctx := context.Background()

	// Cancel after 5 seconds
	ctx, stop := context.WithTimeout(ctx, 5*time.Second)
	defer stop()

	// Cancel if program gets SIGINT
	ctx, stop2 := signal.NotifyContext(ctx, os.Interrupt)
	defer stop2()

	// Equivalent to context.WithTimeout
	//
	//ctx, cancel := context.WithCancel(ctx)
	//go func() {
	//	time.Sleep(5 * time.Second)
	//	cancel()
	//}()

	heartbeats := make(chan time.Time)
	go RunTicker(ctx, heartbeats, 1*time.Second)

	for tick := range heartbeats {
		fmt.Println("heartbeat:", tick)
	}
}
