package main

import (
	"os"
)

type LsCmd struct {
	All  bool   `short:"a"`
	Path string `arg:"" optional:""`
}

func (c *LsCmd) Run() error {
	if c.Path == "" {
		var err error
		c.Path, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	dirents, err := os.ReadDir(c.Path)
	if err != nil {
		return err
	}

	for _, ent := range dirents {
		name := ent.Name()

		if !c.All && len(name) > 0 && name[0] == '.' {
			continue
		}

		os.Stdout.WriteString(name)
		os.Stdout.WriteString("\n")
	}

	return nil
}
