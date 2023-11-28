package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/alecthomas/kong"
)

func main() {

	var cli struct {
		Serve ServeCmd `cmd:"" help:"serve the counter service"`
		Get   GetCmd   `cmd:"" help:"get counter value"`
		Set   SetCmd   `cmd:"" help:"set counter value"`
		Add   AddCmd   `cmd:"" help:"add to counter value"`
	}
	ctx := kong.Parse(&cli)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

func getCounter(url string) (int, error) {
	resp, err := http.Get(url + "/counter")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Counter int `json:"counter"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	return result.Counter, nil
}

func changeCounter(method, url string, val int) (int, error) {

	var payload struct {
		Counter int `json:"counter"`
	}
	payload.Counter = val

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return 0, err
	}

	req, err := http.NewRequest(method, url+"/counter", &buf)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, err
	}
	return payload.Counter, nil
}
