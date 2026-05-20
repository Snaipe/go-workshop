package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"sync"
)

type NcCmd struct {
	Host string `arg:"" optional:""`
	Port string `arg:"" optional:""`

	Listen   string `short:"l"`
	KeepOpen bool   `short:"k"`
	Command  string `short:"c"`
}

func client(host, port string) error {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		return err
	}
	defer conn.Close()

	errs := make(chan error, 2)
	go func() {
		_, err := io.Copy(conn, os.Stdin)
		errs <- err
	}()
	go func() {
		_, err := io.Copy(os.Stdout, conn)
		errs <- err
	}()

	err = errors.Join(<-errs, <-errs)
	if err != nil {
		return err
	}

	return nil
}

func handleConn(conn net.Conn, cmd string) {
	if cmd != "" {
		c := exec.Command("sh", "-xc", cmd)
		c.Stdout = conn
		c.Stderr = conn
		c.Stdin = conn

		if err := c.Run(); err != nil {
			fmt.Printf("exec %v: %v\n", cmd, err)
		}
	} else {
		var redirectWg sync.WaitGroup
		redirectWg.Add(2)
		go func() {
			defer redirectWg.Done()
			_, err := io.Copy(conn, os.Stdin)
			fmt.Println("copy conn<-stdin:", err)
		}()
		go func() {
			defer redirectWg.Done()
			_, err := io.Copy(os.Stdout, conn)
			fmt.Println("copy conn<-stdin:", err)
		}()

		redirectWg.Wait()
	}
}

func server(port string, keepOpen bool, cmd string) error {

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	defer l.Close()

	var wg sync.WaitGroup
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			handleConn(conn, cmd)
		}()

		if !keepOpen {
			break
		}
	}
	wg.Wait()

	return nil
}

func (c *NcCmd) Run() error {
	if c.Listen != "" {
		return server(c.Listen, c.KeepOpen, c.Command)
	} else {
		return client(c.Host, c.Port)
	}
}
