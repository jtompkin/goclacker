// Copyright 2024 Josh Tompkin
// Licensed under the MIT License

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/jtompkin/goclacker/internal/stack"
	"golang.org/x/term"
)

// interactive is the windows implementation of interactive mode. It returns
// io.EOF to signify a normal exit, and any other error signifies an abnormal
// exit.
func interactive(so *stack.StackOperator, color bool) (err error) {
	it := term.NewTerminal(os.Stdin, "")
	c := colors{}
	if color {
		c.out = it.Escape.Yellow
		c.err = it.Escape.Red
		c.reset = it.Escape.Reset
	}
	for {
		fmt.Print(so.Prompt())
		line, err := it.ReadLine()
		if err != nil {
			return err
		}
		err = so.ParseInput(line)
		if err == io.EOF {
			return io.EOF
		}
		if bytes.Count(so.PrintBuf, []byte{'\n'}) == 1 {
			fmt.Print(string(c.out))
		}
		fmt.Print(string(so.PrintBuf))
		if err != nil {
			fmt.Print(string(c.err))
			fmt.Fprint(os.Stderr, err)
		}
		fmt.Print(string(c.reset))
	}
}
