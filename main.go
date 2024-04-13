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

const usage = `Usage of goclacker:
    -s, --stack-limit int
        stack size limit
    -w, --words-file string
        path to file containing word definitions
    -n, --no-count
        do not print stack size in interactive prompt
`

func makeStackOperator(stackLimit int) *stack.StackOperator {
    help := map[string]string{
        "+": "pop 'a', 'b'; push the result of 'a' + 'b'",
        "-":"pop 'a', 'b'; push the result of 'b' - 'a'",
        "*":"pop 'a', 'b'; push the result of 'a' * 'b'",
        "/":"pop 'a', 'b'; push the result of 'b' / 'a'",
        "^":"pop 'a', 'b'; push the result of 'b' ^ 'a'",
        "log":,
        "ln":,
        ".":,
        "stash":,
    }
	operationMap := map[string]*stack.Operation{
		"+":     stack.NewOperation(actions.Add, 2, 1, help["+"]),
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
		"clear": stack.NewOperation(actions.Clear, 0, 0),
		"words": stack.NewOperation(actions.Words, 0, 0),
		"help":  stack.NewOperation(actions.Help, 0, 0),
	}
	return stack.NewStackOperator(operationMap, stackLimit)
}

func withCount(stkOp *stack.StackOperator) {
	fmt.Printf(" %d > ", stkOp.Stack.Len())
}

func noCount(_ *stack.StackOperator) {
	fmt.Print("  > ")
}

func interactive(stkOp *stack.StackOperator, promptFunc func(*stack.StackOperator)) {
	scanner := bufio.NewScanner(os.Stdin)
	promptFunc(stkOp)
	for scanner.Scan() {
		stkOp.ParseInput(scanner.Text())
		promptFunc(stkOp)
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
	failed := false
	for scanner.Scan() {
		if err := stkOp.DefWord(strings.Split(scanner.Text(), " ")); err != nil {
			fmt.Printf("definition error: %s\n", err)
			failed = true
		}
	}

	if failed {
		fmt.Fprintln(os.Stderr, "enter 'help' to see list of operators that cannot be used as words")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var wordsPath string
	flag.StringVar(&wordsPath, "words-file", "", "path to file containing word definitions")
	flag.StringVar(&wordsPath, "w", "", "path to file containing word definitions")
	var stackLimit int
	flag.IntVar(&stackLimit, "s", 8, "stack size limit")
	flag.IntVar(&stackLimit, "stack-limit", 8, "stack size limit")
	var printCount bool
	flag.BoolVar(&printCount, "n", false, "do not print stack size in interactive prompt")
	flag.BoolVar(&printCount, "no-count", false, "do not print stack size in interactive prompt")

	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	stkOp := makeStackOperator(stackLimit)
	if wordsPath != "" {
		parseWordsFile(stkOp, wordsPath)
	}
	if len(flag.Args()) > 0 {
		for _, program := range flag.Args() {
			stkOp.ParseInput(program)
		}
	} else {
		promptFunc := withCount
		if printCount {
			promptFunc = noCount
		}
		interactive(stkOp, promptFunc)
	}
}
