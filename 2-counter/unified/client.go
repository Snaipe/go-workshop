package main

import (
	"encoding/json"
	"bytes"
	"fmt"
	"os"
	"io"
	"net/http"
	"net/url"
)

type Counter struct {
	Counter int `json:"counter"`
}

func newRequest(method string, hostport string, body io.Reader) (*http.Request, error) {
	u, err := url.Parse("http://replaceme/counter")
	if err != nil {
		return nil, err
	}
	u.Host = hostport

	return http.NewRequest(method, u.String(), body)
}

type GetCmd struct{
	Address string `arg:""`
}

func (c *GetCmd) Run() error {

	req, err := newRequest("GET", c.Address, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// On décode le json du le corps de réponse dans la variable result
	var result Counter
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	fmt.Println("counter:", result.Counter)
	return nil
}

type SetCmd struct{
	Address string `arg:""`
	N       int    `arg:""`
}

func (c *SetCmd) Run() error {
	result := Counter{Counter: c.N}

	var out bytes.Buffer
	json.NewEncoder(&out).Encode(result)

	req, err := newRequest("PUT", c.Address, &out)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
	io.Copy(os.Stdout, resp.Body)
	return nil
}

type AddCmd struct{
	Address string `arg:""`
	N       int    `arg:""`
}

func (c *AddCmd) Run() error {
	inc := Counter{Counter: c.N}

	var out bytes.Buffer
	json.NewEncoder(&out).Encode(inc)

	req, err := newRequest("PATCH", c.Address, &out)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("PATCH: %v", resp.Status)
	}
	return nil
}
