package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Result struct {
	URL   string
	Body  string
	Error error
}

func main() {

	urls := []string{
		"https://google.com",
		"https://caliamis.net",
		"not-valid",
	}

	results := make(chan Result, len(urls))

	ctx := context.Background()

	ctx, stop := context.WithTimeout(ctx, 5 * time.Second)
	defer stop()

	var wg sync.WaitGroup

	for _, u := range urls {
		u := u // NOTE! necessary; see loop-variable.go

		wg.Add(1)
		go func() {
			defer wg.Done()

			req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
			if err != nil {
				results <- Result{URL: u, Error: err}
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				results <- Result{URL: u, Error: err}
				return
			}
			var out strings.Builder
			if _, err := io.Copy(&out, resp.Body); err != nil {
				results <- Result{URL: u, Error: err}
				return
			}
			results <- Result{URL: u, Body: out.String()}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.Error != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", result.Error)
		} else {
			fmt.Printf("%s: %q\n", result.URL, result.Body)
		}
	}
}
