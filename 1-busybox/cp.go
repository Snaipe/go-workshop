package main

import (
	"io"
	"os"
)

type CpCmd struct {
	Source      string `arg:""`
	Destination string `arg:""`
}

func (c *CpCmd) Run() error {

	src, err := os.Open(c.Source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(c.Destination)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return dst.Close()
}
