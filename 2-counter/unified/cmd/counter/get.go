package main

import (
	"fmt"
)

type GetCmd struct {
	HostPort string `default:"localhost:1323" help:"the host:port to connect to"`
}

func (cmd *GetCmd) Run() error {
	counter, err := getCounter("http://"+cmd.HostPort)
	if err != nil {
		return err
	}

	fmt.Println("counter:", counter)
	return nil
}
