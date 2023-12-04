package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"net"
)

// Tout type implémentant la méthode Write(p []byte) (n int, err error)
// implémente automatiquement cette interface
type Writer interface {
	Write(p []byte) (n int, err error)
}

// package fmt
//
// func Fprintf(w io.Writer, format string, a ...any) (n int, err error)

func main() {

	sujet := "tout le monde"
	if len(os.Args) > 2 {
		sujet = os.Args[1]
	}

	// Écrire dans un *os.File
	fmt.Fprintf(os.Stdout, "Bonjour %s!\n", sujet)

	// Écrire dans un net.Conn
	sock, _ := net.Dial("tcp", "localhost:1234")
	fmt.Fprintf(sock, "Bonjour %s!\n", sujet)

	// Écrire dans un buffer
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Bonjour %s!\n", sujet)

	// Écrire dans un string builder
	var str strings.Builder
	fmt.Fprintf(&str, "Bonjour %s!\n", sujet)
}

func WriteAndClose(writeCloser io.WriteCloser, buf []byte) {
	_, err := writeCloser.Write(buf)
	if err != nil {
		return err
	}
	return writeCloser.Close()
}
