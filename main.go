package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/jtompkin/goclacker/actions"
	"github.com/jtompkin/goclacker/stack"
)

func interactive() {
	operationMap := map[string]*stack.Operation{
		"+":     stack.NewOperation(actions.Add, 2, 1),
		"-":     stack.NewOperation(actions.Subtract, 2, 1),
		"*":     stack.NewOperation(actions.Multiply, 2, 1),
		"/":     stack.NewOperation(actions.Divide, 2, 1),
		"^":     stack.NewOperation(actions.Power, 2, 1),
		".":     stack.NewOperation(actions.Display, 0, 0),
		",":     stack.NewOperation(actions.Pop, 1, 0),
		"clear": stack.NewOperation(actions.Clear, 0, 0),
	}
	stkOp := stack.NewStackOperator(operationMap, 8)
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
