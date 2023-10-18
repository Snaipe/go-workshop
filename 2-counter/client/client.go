package main

import (
	"bytes"
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"os"
)

func main() {
	var payload struct {
		Counter int `json:"counter"`
	}

	var err error
	payload.Counter, err = strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "http://localhost:1323/counter", &buf)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var result struct {
		Counter int `json:"counter"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("counter:", result.Counter)
}
