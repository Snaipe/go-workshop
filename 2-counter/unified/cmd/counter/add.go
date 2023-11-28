package main

import (
	"fmt"
)

type AddCmd struct {
	HostPort string `default:"localhost:1323" help:"the host:port to connect to"`
	N        int    `arg:"" help:"the value to set the counter to"`
}

func (cmd *AddCmd) Run() error {
	counter, err := changeCounter("PATCH", "http://"+cmd.HostPort, cmd.N)
	if err != nil {
		return err
	}

	fmt.Println("counter:", counter)
	return nil
}
