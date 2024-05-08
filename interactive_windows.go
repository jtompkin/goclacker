package main

import (
	"fmt"
	"os"

	"github.com/jtompkin/goclacker/internal/stack"
	"golang.org/x/term"
)

func interactive(so *stack.StackOperator) (err error) {
	it := term.NewTerminal(os.Stdin, "")
	for {
		fmt.Print(so.Prompt())
		line, err := it.ReadLine()
		if err != nil {
			it.SetPrompt("")
			it.Write(nil)
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
