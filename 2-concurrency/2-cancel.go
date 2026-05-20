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

	ticker := time.NewTicker(period)
	defer ticker.Stop()

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
		case <-ticker.C:
		}
	}
}

func main() {

	ctx := context.Background()

	// Annule le context après 5 secondes
	ctx, stop := context.WithTimeout(ctx, 5*time.Second)
	defer stop()

	// Annule le context après avoir reçu SIGINT (Ctrl-C)
	ctx, stop2 := signal.NotifyContext(ctx, os.Interrupt)
	defer stop2()

	// Équivalent à context.WithTimeout
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
