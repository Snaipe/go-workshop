package main

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"
	"bytes"
	"io"
	"os"
)

func main() {
	// Fait une requête HTTP GET http://localhost:1323/counter
	resp, err := http.Get("http://localhost:1323/counter")
	if err != nil {
		log.Fatal(err)
	}

	// Quand la fonction se termine, on ferme le flux de réponse
	defer resp.Body.Close()

	var result struct {
		Counter int `json:"counter"`
	}

	// On décode le json du le corps de réponse dans la variable result
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}

	fmt.Println("counter:", result.Counter)

	result.Counter = 42

	var out bytes.Buffer
	json.NewEncoder(&out).Encode(result)

	req, err := http.NewRequest("PUT", "http://localhost:1323/counter", &out)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
	io.Copy(os.Stdout, resp.Body)
}
