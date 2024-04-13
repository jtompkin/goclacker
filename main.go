package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jtompkin/goclacker/actions"
	"github.com/jtompkin/goclacker/stack"
)


func makeStackOperator(stackLimit int) *stack.StackOperator {
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
	return stack.NewStackOperator(operationMap, stackLimit)
}

func interactive(stkOp *stack.StackOperator) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("  > ")
	for scanner.Scan() {
		stkOp.ParseInput(scanner.Text())
		fmt.Print("  > ")
	}
	fmt.Println()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func parseWordsFile(stkOp *stack.StackOperator, path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := stkOp.DefWord(strings.Split(scanner.Text(), " ")); err != nil {
			fmt.Println(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	wordsPath := flag.String("w", "", "path to file containing word definitions")
    stackLimit := flag.Int("s", 8, "stack size limit")

	flag.Parse()

	stkOp := makeStackOperator(*stackLimit)
	if *wordsPath != "" {
		parseWordsFile(stkOp, *wordsPath)
	}
	if len(flag.Args()) > 0 {
        for _, program := range flag.Args() {
            stkOp.ParseInput(program)
        }
	} else {
		interactive(stkOp)
	}
}
