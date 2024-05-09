//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos

package main

import (
	"io"
	"os"
	"strings"

	"github.com/jtompkin/goclacker/internal/stack"
	"golang.org/x/term"
)

func interactive(so *stack.StackOperator) (err error) {
	fds := []int{int(os.Stdin.Fd()), int(os.Stderr.Fd())}
	for _, fd := range fds {
		state, err := term.MakeRaw(fd)
		if err != nil {
			return err
		}
		defer term.Restore(fd, state)
	}

	it := term.NewTerminal(os.Stdin, so.Prompt())
	et := term.NewTerminal(os.Stderr, "")
	for {
		line, err := it.ReadLine()
		if strings.TrimSpace(line) == "quit" {
			it.SetPrompt("")
			it.Write(nil)
			return io.EOF
		}
		if err != nil {
			it.SetPrompt("")
			it.Write(nil)
			return err
		}
		err = so.ParseInput(line)
		it.Write(so.PrintBuf)
		if err != nil {
			et.Write([]byte(err.Error()))
		}
		it.SetPrompt(so.Prompt())
	}
}
