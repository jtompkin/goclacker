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

	"github.com/jtompkin/goclacker/internal/actions"
	"github.com/jtompkin/goclacker/internal/stack"
)

const usage string = `Usage of goclacker:
    -V, --version
        print version information
    -h, --help
        print usage information
    -s, --stack-limit int
        stack size limit (default 8); must be non-negative
    -w, --words-file string
        path to file containing word definitions
    -p, --prompt string
        format string for the interactive prompt (default " &c > ")
        format specifiers:
            &l : stack limit
            &c : current stack size
            &t : top stack value
            &s : current stash value
`
const version string = "v0.1.1"
const fmtChar byte = '&'

const defStackLimit int = 8
const defPrompt string = " &c > "

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
		"+", "-", "*", "/", "^", "log", "ln", "rad", "deg", "sin", "cos", "tan",
		"stash", "pull", "round", ".", ",", "clear", "words", "help",
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
	return stack.NewStackOperator(actionMap, orderedTokens, stackLimit)
}

func interactive(so *stack.StackOperator, promptFormat string) {
	scanner := bufio.NewScanner(os.Stdin)
	so.MakePromptFunc(promptFormat, fmtChar)
	fmt.Print(so.Prompt())
	for scanner.Scan() {
		so.ParseInput(scanner.Text())
		fmt.Print(so.Prompt())
	}
	fmt.Println()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func parseWordsFile(so *stack.StackOperator, path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var failed bool
	for scanner.Scan() {
		if err := so.DefWord(strings.Split(scanner.Text(), " ")); err != nil {
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
	flag.IntVar(&stackLimit, "s", defStackLimit, "stack size limit")
	flag.IntVar(&stackLimit, "stack-limit", defStackLimit, "stack size limit")
	var promptFormat string
	flag.StringVar(&promptFormat, "p", defPrompt, "format string for the interactive prompt")
	flag.StringVar(&promptFormat, "prompt", defPrompt, "format string for the interactive prompt")
	var printVersion bool
	flag.BoolVar(&printVersion, "V", false, "print version information")
	flag.BoolVar(&printVersion, "version", false, "print version information")

	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	if stackLimit < 0 {
		fmt.Print("-s, --stack-limit must be non-negative\n\n")
		fmt.Print(usage)
		os.Exit(1)
	}

	if printVersion {
		fmt.Printf("goclacker %s\n", version)
		os.Exit(0)
	}

	so := makeStackOperator(stackLimit)
	if wordsPath != "" {
		parseWordsFile(so, wordsPath)
	}
	if len(flag.Args()) > 0 {
		for _, program := range flag.Args() {
			so.ParseInput(program)
		}
		os.Exit(0)
	}
	interactive(so, promptFormat)
}
