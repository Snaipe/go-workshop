package main

type SetCmd struct {
	HostPort string `default:"localhost:1323" help:"the host:port to connect to"`
	N        int    `arg:"" help:"the value to set the counter to"`
}

func (cmd *SetCmd) Run() error {
	_, err := changeCounter("PUT", "http://"+cmd.HostPort, cmd.N)
	return err
}

