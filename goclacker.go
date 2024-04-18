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
	"github.com/wk8/go-ordered-map/v2"
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
const version string = "v0.2.0"
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

func MakeStackOperator(stackLimit int, interactive bool) *stack.StackOperator {
	actionMap := *orderedmap.New[string, stack.Action]()
	actionMap.Set("+", actions.Add())
	actionMap.Set("-", actions.Subtract())
	actionMap.Set("*", actions.Multiply())
	actionMap.Set("/", actions.Divide())
    actionMap.Set("%", actions.Modulo())
    actionMap.Set("^", actions.Power())
    actionMap.Set("!", actions.Factorial())
    actionMap.Set("log", actions.Log())
    actionMap.Set("ln", actions.Ln())
    actionMap.Set("rad", actions.Radians())
    actionMap.Set("deg", actions.Degrees())
    actionMap.Set("sin", actions.Sine())
    actionMap.Set("cos", actions.Cosine())
    actionMap.Set("tan", actions.Tangent())
    actionMap.Set("floor", actions.Floor())
    actionMap.Set("ceil", actions.Ceiling())
    actionMap.Set("round", actions.Round())
    actionMap.Set(".", actions.Display())
    actionMap.Set(",", actions.Pop())
    actionMap.Set("stash", actions.Stash())
    actionMap.Set("pull", actions.Pull())
    actionMap.Set("clear", actions.Clear())
    actionMap.Set("words", actions.Words())
    actionMap.Set("help", actions.Help())
	return stack.NewStackOperator(actionMap, stackLimit, interactive)
}

func nonInteractive(so *stack.StackOperator, programs []string) {
	for _, s := range programs {
		if s, err := so.ParseInput(s); err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
		} else {
			fmt.Printf("%s", s)
		}
	}
}

func interactive(so *stack.StackOperator, promptFormat string) {
	scanner := bufio.NewScanner(os.Stdin)
	so.MakePromptFunc(promptFormat, fmtChar)
	fmt.Print(so.Prompt())
	for scanner.Scan() {
		if s, err := so.ParseInput(scanner.Text()); err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
		} else {
			fmt.Printf("%s", s)
		}
		fmt.Print(so.Prompt())
	}
	fmt.Println()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func ParseWordsFile(so *stack.StackOperator, path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var failed bool
	for scanner.Scan() {
		if _, err := so.DefWord(strings.Split(scanner.Text(), " ")); err != nil {
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

	so := MakeStackOperator(stackLimit, !(len(flag.Args()) > 0))
	if wordsPath != "" {
		ParseWordsFile(so, wordsPath)
	}
	if so.Interactive {
		interactive(so, promptFormat)
	} else {
		nonInteractive(so, flag.Args())
	}
}
