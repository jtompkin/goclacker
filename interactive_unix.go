// Copyright 2024 Josh Tompkin
// Licensed under the MIT License

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos

package main

import (
	"io"
	"os"
	"strings"

	"github.com/jtompkin/goclacker/internal/stack"
	"golang.org/x/term"
)

// interactive is the unix-like implementation of interactive mode.
func interactive(so *stack.StackOperator) (err error) {
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), state)

	it := term.NewTerminal(os.Stdin, so.Prompt())
	ot := term.NewTerminal(os.Stdout, "")
	et := term.NewTerminal(os.Stderr, "")
	for {
		line, err := it.ReadLine()
		if strings.TrimSpace(line) == "quit" {
			return io.EOF
		}
		if err != nil {
			it.SetPrompt("")
			it.Write(nil)
			return err
		}
		err = so.ParseInput(line)
		ot.Write(so.PrintBuf)
		if err != nil {
			et.Write([]byte(err.Error()))
		}
		it.SetPrompt(so.Prompt())
	}
}
