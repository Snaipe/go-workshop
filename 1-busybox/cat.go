package main

import (
	"bufio"
	"fmt"
	"os"
)

type CatCmd struct {
	Count   bool `short:"n"`
	Escapes bool `short:"e"`

	Paths []string `arg:""`
}

func (c *CatCmd) Run() error {

	var files []*os.File
	if len(c.Paths) == 0 {
		files = append(files, os.Stdin)
	} else {
		for _, path := range c.Paths {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			files = append(files, f)
		}
	}

	lineno := 1
	for _, f := range files {
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			if c.Count {
				fmt.Printf("% 8d %s", lineno, scanner.Text())
			} else {
				os.Stdout.WriteString(scanner.Text())
			}

			if c.Escapes {
				os.Stdout.WriteString("$")
			}
			os.Stdout.WriteString("\n")
			lineno++
		}
	}

	return nil
}
