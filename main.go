package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/jtompkin/goclacker/actions"
	"github.com/jtompkin/goclacker/stack"
)

const StackLimit int = 8

func interactive() {
	operationMap := map[string]*stack.Operation{
		"+":     stack.NewOperation(actions.Add, 2, 1),
		"-":     stack.NewOperation(actions.Subtract, 2, 1),
		"*":     stack.NewOperation(actions.Multiply, 2, 1),
		"/":     stack.NewOperation(actions.Divide, 2, 1),
		"^":     stack.NewOperation(actions.Power, 2, 1),
		"log":   stack.NewOperation(actions.Log, 1, 1),
		"ln":    stack.NewOperation(actions.Ln, 1, 1),
		".":     stack.NewOperation(actions.Display, 0, 0),
		",":     stack.NewOperation(actions.Pop, 1, 0),
		"stash": stack.NewOperation(actions.Stash, 1, 0),
		"pull":  stack.NewOperation(actions.Pull, 0, 1),
		"round": stack.NewOperation(actions.Round, 2, 1),
		"words": stack.NewOperation(actions.Words, 0, 0),
		"clear": stack.NewOperation(actions.Clear, 0, 0),
	}
	stkOp := stack.NewStackOperator(operationMap, StackLimit)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("  > ")
		// Catch EOF
		if !scanner.Scan() {
			fmt.Println()
			return
		}
		if err := scanner.Err(); err != nil {
			fmt.Println()
			log.Fatal(err)
		}
		stkOp.ParseInput(scanner.Text())
	}
}

func main() {
	interactive()
}
