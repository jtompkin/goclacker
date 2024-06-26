// Copyright 2024 Josh Tompkin
// Licensed under the MIT License

package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jtompkin/goclacker/internal/stack"
	"golang.org/x/term"
)

// interactive is the windows implementation of interactive mode.
func interactive(so *stack.StackOperator) (err error) {
	it := term.NewTerminal(os.Stdin, "")
	for {
		fmt.Print(so.Prompt())
		line, err := it.ReadLine()
		if strings.TrimSpace(line) == "quit" {
			return io.EOF
		}
		if err != nil {
			return err
		}
		err = so.ParseInput(line)
		it.Write(so.PrintBuf)
		fmt.Printf("%s", so.PrintBuf)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}
}
