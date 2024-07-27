// Copyright 2024 Josh Tompkin
// Licensed under the MIT License

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos

package main

import (
	"bytes"
	"io"
	"os"

	"github.com/jtompkin/goclacker/internal/stack"
	"golang.org/x/term"
)

// interactive is the unix-like implementation of interactive mode. Returns an
// io.EOF on graceful exit.
func interactive(so *stack.StackOperator, color bool) (err error) {
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), state)

	it := term.NewTerminal(os.Stdin, so.Prompt())
	ot := term.NewTerminal(os.Stdout, "")
	et := term.NewTerminal(os.Stderr, "")
	c := colors{}
	if color {
		c.out = ot.Escape.Yellow
		c.err = ot.Escape.Red
		c.reset = ot.Escape.Reset
	}
	for {
		line, err := it.ReadLine()
		if err != nil {
			return err
		}
		err = so.ParseInput(line)
		if err == io.EOF {
			return io.EOF
		}
		if bytes.Count(so.ToPrint, []byte{'\n'}) < 2 {
			ot.Write(c.out)
		}
		ot.Write(so.ToPrint)
		if err != nil {
			ot.Write(c.err)
			et.Write([]byte(err.Error()))
		}
		ot.Write(c.reset)
		it.SetPrompt(so.Prompt())
	}
}
