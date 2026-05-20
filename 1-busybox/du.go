package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type DuCmd struct {
	Path string `arg:"" optional:""`
}

func (c *DuCmd) Run() error {

	type SizeStack struct {
		Path string
		Size int64
	}

	var sizes []SizeStack

	tally := func(depth int) {
		for len(sizes) > depth {
			last := sizes[len(sizes)-1]
			fmt.Printf("%s %d\n", last.Path, last.Size)

			sizes[len(sizes)-2].Size += last.Size
			sizes = sizes[:len(sizes)-1]
		}
	}

	filepath.WalkDir(c.Path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		depth := strings.Count(path, string(filepath.Separator)) + 1

		if len(sizes) < depth {
			sizes = append(sizes, SizeStack{Path: path})
		}

		tally(depth)

		size := &sizes[depth-1]

		if d.Type().IsRegular() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			size.Size += info.Size()
		}

		return err
	})

	tally(1)

	fmt.Printf("%s %d\n", sizes[0].Path, sizes[0].Size)
	return nil
}
