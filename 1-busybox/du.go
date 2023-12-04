package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type DuCmd struct {
	Paths []string `arg:"" optional:"" default:"."`
}

func (cmd *DuCmd) Run() error {

	for _, path := range cmd.Paths {

		abs, err := filepath.Abs(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		dirs := []string{}
		total := []int64{}

		pushDir := func(path string) {
			dirs = append(dirs, path)
			total = append(total, 0)
		}

		popDir := func() {
			dir := dirs[len(dirs)-1]
			dirs = dirs[:len(dirs)-1]

			size := total[len(total)-1]
			total = total[:len(total)-1]
			if len(total) > 0 {
				total[len(total)-1] += size
			}

			fmt.Printf("%d\t%s\n", size, filepath.Join(path, dir[len(abs):]))
		}

		err = filepath.WalkDir(abs, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return fs.SkipDir
			}

			fi, err := d.Info()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return nil
			}

			depth := strings.Count(path, "/")
			if d.IsDir() && depth >= len(dirs) {
				pushDir(path)
				return nil
			}

			if depth < len(dirs) {
				popDir()
			}

			if !d.IsDir() {
				total[len(total)-1] += fi.Size()
			}
			return nil
		})

		for len(dirs) > 0 {
			popDir()
		}
	}
	return nil
}
