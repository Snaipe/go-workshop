package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// bufferPool is a pool of preallocated bytes.Buffer.
type bufferPool chan *bytes.Buffer

// makeBufferPool creates a buffer pool with the specified size.
func makeBufferPool(size int) bufferPool {
	pool := make(bufferPool, size)
	for i := 0; i < size; i++ {
		pool <- &bytes.Buffer{}
	}
	return pool
}

func (p bufferPool) Acquire() *bytes.Buffer {
	select {
	case buf := <-p:
		buf.Reset()
		return buf
	default:
		return &bytes.Buffer{}
	}
}

func (p bufferPool) Release(buf *bytes.Buffer) {
	select {
	case p <- buf:
	default:
		// NOTE: pool is full, let the GC garbage collect the buffer
	}
}

// Check200 checks whether each of the HTTP URLs return a 200 status.
//
// It returns a slice of errors with the same lenght as in. The error in
// errs[i] is nil if the request of in[i] resulted in a 200 response.
func Check200(ctx context.Context, in []string) (errs []error) {

	// Run up to 16 requests at a time
	const maxConcurrentRequests = 16

	pool := make(chan struct{}, maxConcurrentRequests)

	doRequest := func(addr string) error {
		ctx, stop := context.WithTimeout(ctx, 1*time.Minute)
		defer stop()

		req, err := http.NewRequestWithContext(ctx, "GET", addr, nil)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		_ = resp.Body.Close() // we do not care about the body or I/O errors

		if resp.StatusCode != http.StatusOK {
			// *url.Error is a standard error type, so let's use it to enrich
			// our error.
			return &url.Error{
				Op:  "Check200",
				URL: addr,
				Err: fmt.Errorf("status code %d is not 200", resp.StatusCode),
			}
		}

		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(in))

	for i := range in {
		i := i
		pool <- struct{}{} // Acquire a request slot
		go func() {
			defer func() {
				wg.Done()
				<-pool // Release the request slot back to the pool
			}()
			errs[i] = doRequest(in[i])
		}()
	}

	wg.Wait()

	return errs
}
