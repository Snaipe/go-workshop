package main

import (
	"fmt"
	"bufio"
	"os"
	"strconv"
	"strings"
	"bytes"
)

type CatCmd struct {
	Lines  bool `short:"n"`
	Escape bool `short:"e"`

	Files []string `arg:"" optional:"" default:"-"`
}

func (cmd *CatCmd) Run() error {
	lineno := 1
	for _, file := range cmd.Files {
		var err error
		lineno, err = cmd.cat(file, lineno)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cmd *CatCmd) cat(src string, lineno int) (int, error) {
	var in *os.File
	if src == "-" {
		in = os.Stdin
	} else {
		f, err := os.Open(src)
		if err != nil {
			return lineno, err
		}
		defer f.Close()
		in = f
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	scanner := bufio.NewScanner(in)
	scanner.Split(ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if cmd.Escape {
			line = strconv.Quote(line)
			line = line[1:len(line)-1]
			line = strings.ReplaceAll(line, `\n`, "$\n")
		}
		if cmd.Lines {
			fmt.Fprintf(out, "% 6d\t%s", lineno, line)
		} else {
			fmt.Fprint(out, line)
		}
		lineno++
	}
	return lineno, scanner.Err()
}

func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i+1], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
