package main

import (
	"fmt"
	"os"
	"strings"
	"slices"
)

type LsCmd struct {
	All bool `short:"a"`

	Paths []string `arg:"" optional:"" default:"."`
}

func (cmd *LsCmd) Run() error {
	for _, path := range cmd.Paths {
		dirent, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		filenames := make([]string, 0, len(dirent)+2)

		if cmd.All {
			filenames = append(filenames, ".")
			filenames = append(filenames, "..")
		}
		for _, d := range dirent {
			if !cmd.All && strings.HasPrefix(d.Name(), ".") {
				continue
			}
			filenames = append(filenames, d.Name())
		}

		slices.SortStableFunc(filenames, func(a, b string) int {
			a = strings.TrimLeft(a, ".")
			b = strings.TrimLeft(b, ".")
			switch {
			case a < b:
				return -1
			case a > b:
				return 1
			default:
				return 0
			}
		})
		for _, name := range filenames {
			fmt.Println(name)
		}
	}
	return nil
}
