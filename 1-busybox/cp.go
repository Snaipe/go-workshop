package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type CpCmd struct {
	Paths []string `arg:""`
}

func (cmd *CpCmd) Run() error {

	if len(cmd.Paths) == 2 {
		return cp(cmd.Paths[0], cmd.Paths[1])
	}

	numops := len(cmd.Paths)-1
	errs := make(chan error, numops)

	dest := cmd.Paths[len(cmd.Paths)-1]
	for _, path := range cmd.Paths[:numops] {
		path := path
		go func() {
			newpath := filepath.Join(dest, filepath.Base(path))
			errs <- cp(path, newpath)
		}()
	}

	var hadError bool
	for i := 0; i < numops; i++ {
		err := <-errs
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			hadError = true
		}
	}
	if hadError {
		os.Exit(1)
	}

	return nil
}

func cp(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Close()
}
