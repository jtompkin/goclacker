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
        print usge information
    -s, --stack-limit int
        stack size limit (default 8)
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

func makePrompt(so *stack.StackOperator, format string) string {
	flags := map[byte]func(*stack.StackOperator) (float64, error){
		'l': func(so *stack.StackOperator) (float64, error) { return float64(cap(so.Stack.Values)), nil },
		't': func(so *stack.StackOperator) (float64, error) { return so.Stack.Top() },
		'c': func(so *stack.StackOperator) (float64, error) { return float64(len(so.Stack.Values)), nil },
		's': func(so *stack.StackOperator) (float64, error) { return so.Stack.Stash, nil },
	}
	var prompt string
	var i int
	for i < len(format) {
		if c := format[i]; c == '&' {
			if i == len(format)-1 {
				return prompt
			}
			if fmtFunc := flags[format[i+1]]; fmtFunc != nil {
				if f, err := fmtFunc(so); err != nil {
					prompt += fmt.Sprint(err)
				} else {
					prompt += fmt.Sprint(f)
				}
			}
			i += 2
		} else {
			prompt += fmt.Sprintf("%c", c)
			i++
		}
	}
	return prompt
}

func interactive(so *stack.StackOperator, promptFormat string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(makePrompt(so, promptFormat))
	for scanner.Scan() {
		so.ParseInput(scanner.Text())
		fmt.Print(makePrompt(so, promptFormat))
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
	var promptFormat string
	flag.StringVar(&promptFormat, "p", " &c > ", "format string for the interactive prompt")
	flag.StringVar(&promptFormat, "prompt", " &c > ", "format string for the interactive prompt")
	var printVersion bool
	flag.BoolVar(&printVersion, "V", false, "print version information")
	flag.BoolVar(&printVersion, "version", false, "print version information")

	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	if printVersion {
		fmt.Printf("goclacker %s\n", version)
		return
	}

	stkOp := makeStackOperator(stackLimit)
	if wordsPath != "" {
		parseWordsFile(stkOp, wordsPath)
	}
	if len(flag.Args()) > 0 {
		for _, program := range flag.Args() {
			stkOp.ParseInput(program)
		}
		return
	}
	interactive(stkOp, promptFormat)
}
