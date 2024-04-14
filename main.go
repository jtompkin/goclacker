package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
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

func checkTokens(tokens []string, actions map[string]stack.Action) error {
	notFoundOrdered := make([]string, 0, len(actions))
	for k := range actions {
		if slices.Index(tokens, k) == -1 {
			notFoundOrdered = append(notFoundOrdered, k)
		}
	}
	notFoundAction := make([]string, 0, len(tokens))
	for _, s := range tokens {
		if actions[s] == nil {
			notFoundAction = append(notFoundAction, s)
		}
	}
	var err error
	if len(notFoundOrdered) > 0 {
		err = errors.Join(err, errors.New(fmt.Sprintf("%s not found in orderedTokens", notFoundOrdered)))
	}
	if len(notFoundAction) > 0 {
		err = errors.Join(err, errors.New(fmt.Sprintf("%s not found in actionMap", notFoundAction)))
	}
	return err
}

func makeStackOperator(stackLimit int) *stack.StackOperator {
	orderedTokens := []string{
		"+", "-", "*", "/", "^", "log", "ln", ".", ",", "rad", "deg", "sin",
		"cos", "tan", "stash", "pull", "round", "clear", "words", "help",
	}
	actionMap := map[string]stack.Action{
		"+":     actions.Add(),
		"-":     actions.Subtract(),
		"*":     actions.Multiply(),
		"/":     actions.Divide(),
		"^":     actions.Power(),
		"log":   actions.Log(),
		"ln":    actions.Ln(),
		"rad":   actions.Radians(),
		"deg":   actions.Degrees(),
		"sin":   actions.Sine(),
		"cos":   actions.Cosine(),
		"tan":   actions.Tangent(),
		".":     actions.Display(),
		",":     actions.Pop(),
		"stash": actions.Stash(),
		"pull":  actions.Pull(),
		"round": actions.Round(),
		"clear": actions.Clear(),
		"words": actions.Words(),
		"help":  actions.Help(),
	}
	if err := checkTokens(orderedTokens, actionMap); err != nil {
		log.Fatal(err)
	}
	return stack.NewStackOperator(actionMap, &orderedTokens, stackLimit)
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
