package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Input string `default:"-" help:"evaluate input from file"`
	Eval  string `help:"evaluate expression instead or running the REPL"`
}

func main() {
	var cli CLI
	kong.Parse(&cli)

	var in *os.File
	if cli.Input == "-" {
		in = os.Stdin
	} else {
		f, err := os.Open(cli.Input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close() // sera fermé en sortant du main()
		in = f
	}

	scanner := bufio.NewScanner(in)

	for {
		if cli.Input == "-" {
			fmt.Fprint(os.Stderr, "\x1B[94;1mλ\x1B[0m ")
		}
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()

		lex := MakeLexer(line)
		results, err := Eval(lex.Lex())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		for _, result := range results {
			fmt.Println(result)
		}
	}
	if err := scanner.Err(); err != nil { // S'il y a une erreur, l'afficher
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
