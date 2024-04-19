package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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


func MakeStackOperator(stackLimit int, interactive bool) *stack.StackOperator {
	actionMap := *orderedmap.New[string, *stack.Action]()
	actionMap.Set("+", stack.Add())
	actionMap.Set("-", stack.Subtract())
	actionMap.Set("*", stack.Multiply())
	actionMap.Set("/", stack.Divide())
	actionMap.Set("%", stack.Modulo())
	actionMap.Set("^", stack.Power())
	actionMap.Set("!", stack.Factorial())
	actionMap.Set("log", stack.Log())
	actionMap.Set("ln", stack.Ln())
	actionMap.Set("rad", stack.Radians())
	actionMap.Set("deg", stack.Degrees())
	actionMap.Set("sin", stack.Sine())
	actionMap.Set("cos", stack.Cosine())
	actionMap.Set("tan", stack.Tangent())
	actionMap.Set("floor", stack.Floor())
	actionMap.Set("ceil", stack.Ceiling())
	actionMap.Set("round", stack.Round())
    actionMap.Set("rand", stack.Random())
	actionMap.Set(".", stack.Display())
	actionMap.Set(",", stack.Pop())
	actionMap.Set("stash", stack.Stash())
	actionMap.Set("pull", stack.Pull())
	actionMap.Set("clear", stack.Clear())
	actionMap.Set("words", stack.Words())
	actionMap.Set("help", stack.Help())
	return stack.NewStackOperator(actionMap, stackLimit, interactive)
}

func nonInteractive(so *stack.StackOperator, programs []string) {
	for _, s := range programs {
		if s, err := so.ParseInput(s); err != nil {
			fmt.Fprint(os.Stderr, err)
		} else {
			fmt.Print(s)
		}
	}
}

func interactive(so *stack.StackOperator, promptFormat string) {
	scanner := bufio.NewScanner(os.Stdin)
    if err := so.MakePromptFunc(promptFormat, fmtChar); err != nil {
        log.Fatal(err)
    }
	fmt.Print(so.Prompt())
	for scanner.Scan() {
		if s, err := so.ParseInput(scanner.Text()); err != nil {
			fmt.Fprint(os.Stderr, err)
		} else {
			fmt.Print(s)
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
		fmt.Fprint(os.Stderr, "enter 'help' to see list of operators that cannot be used as words\n")
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
