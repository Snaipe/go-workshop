package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"sync"
)

type NcCmd struct {
	Host string `arg:""`
	Port string `arg:"" optional:""`

	Listen   bool   `short:"l"`
	KeepOpen bool   `short:"k"`
	Command  string `short:"c"`
}

type CloseWriter interface {
	CloseWrite() error
}

func (cmd *NcCmd) Run() error {

	if cmd.Port == "" {
		cmd.Port = cmd.Host
		cmd.Host = "127.0.0.1"
	}

	redirect := func(wait *sync.WaitGroup, dst io.Writer, src io.Reader) {
		defer wait.Done()

		_, err := io.Copy(dst, src)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		cw, ok := dst.(CloseWriter)
		if ok {
			cw.CloseWrite()
		}
	}

	if cmd.Listen {
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cmd.Host, cmd.Port))
		if err != nil {
			return err
		}
		defer listener.Close()

		var waitServer sync.WaitGroup
		for {
			conn, err := listener.Accept()
			if err != nil {
				return err
			}

			waitServer.Add(1)
			go func() {
				defer waitServer.Done()
				defer conn.Close()

				if cmd.Command != "" {
					c := exec.Command("/bin/sh", "-c", cmd.Command)
					c.Stdout = conn
					c.Stderr = conn
					stdin, err := c.StdinPipe()
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
					}

					var wait sync.WaitGroup
					wait.Add(1)
					go redirect(&wait, stdin, conn)

					if err := c.Run(); err != nil {
						fmt.Fprintln(os.Stderr, err)
					}

					wait.Wait()
				} else {
					var wait sync.WaitGroup
					wait.Add(2)

					go redirect(&wait, conn, os.Stdin)
					go redirect(&wait, os.Stdout, conn)

					wait.Wait()
				}
			}()

			if !cmd.KeepOpen {
				break
			}
		}
		waitServer.Wait()
	} else {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", cmd.Host, cmd.Port))
		if err != nil {
			return err
		}
		defer conn.Close()

		var wait sync.WaitGroup
		wait.Add(2)

		go redirect(&wait, conn, os.Stdin)
		go redirect(&wait, os.Stdout, conn)

		wait.Wait()
	}
	return nil
}
