package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jtompkin/goclacker/internal/stack"
)

const usage string = `Usage of goclacker:
goclacker [-V] [-h] [-s] [-l] int [-w] string [-p] string [program...]
    -V, --version
        Print version information and exit.
    -h, --help
        Print usage information and exit.
    -s, --strict
        Run in strict mode: entering something that is not a number, operator,
        or defined word will return an error instead of doing nothing.
    -l, --stack-limit int
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
    [program...]
        Any positional arguments will be interpreted and executed by the
        calculator. Interactive mode will not be entered if any positional
        arguments are supplied.
`

const (
	defPrompt string = " &c > "
	version   string = "v1.0.1"
	fmtChar   byte   = '&'
	defLimit  int    = 8
)

func MakeStackOperator(stackLimit int, interactive bool, strict bool) *stack.StackOperator {
	actions := stack.NewOrderedMap[string, *stack.Action](26)
	actions.Set("+", stack.Add())
	actions.Set("-", stack.Subtract())
	actions.Set("*", stack.Multiply())
	actions.Set("/", stack.Divide())
	actions.Set("%", stack.Modulo())
	actions.Set("^", stack.Power())
	actions.Set("!", stack.Factorial())
	actions.Set("log", stack.Log())
	actions.Set("ln", stack.Ln())
	actions.Set("rad", stack.Radians())
	actions.Set("deg", stack.Degrees())
	actions.Set("sin", stack.Sine())
	actions.Set("cos", stack.Cosine())
	actions.Set("tan", stack.Tangent())
	actions.Set("floor", stack.Floor())
	actions.Set("ceil", stack.Ceiling())
	actions.Set("round", stack.Round())
	actions.Set("rand", stack.Random())
	actions.Set(".", stack.Display())
	actions.Set(",", stack.Pop())
	actions.Set("stash", stack.Stash())
	actions.Set("pull", stack.Pull())
	actions.Set("clear", stack.Clear())
	actions.Set("words", stack.Words())
	actions.Set("help", stack.Help())
	actions.Set("cls", stack.Cls())
	return stack.NewStackOperator(actions, stackLimit, interactive, strict)
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
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if failed {
		fmt.Fprint(os.Stderr, "enter 'help' to see list of operators that cannot be used as words\n")
	} else {
        fmt.Printf("sucessfully parsed words file: %q", path)
    }

}

func main() {
	var printVersion bool
	flag.BoolVar(&printVersion, "V", false, "")
	flag.BoolVar(&printVersion, "version", false, "")
	var stackLimit int
	flag.IntVar(&stackLimit, "l", defLimit, "")
	flag.IntVar(&stackLimit, "stack-limit", defLimit, "")
	var strictMode bool
	flag.BoolVar(&strictMode, "s", false, "")
	flag.BoolVar(&strictMode, "strict", false, "")
	var wordsPath string
	flag.StringVar(&wordsPath, "words-file", "", "")
	flag.StringVar(&wordsPath, "w", "", "")
	var promptFormat string
	flag.StringVar(&promptFormat, "p", defPrompt, "")
	flag.StringVar(&promptFormat, "prompt", defPrompt, "")

	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	if stackLimit < 0 {
		fmt.Print("argument error: -s, --stack-limit must be non-negative\n\n")
		fmt.Print(usage)
		os.Exit(1)
	}
	if printVersion {
		fmt.Printf("goclacker %s\n", version)
		os.Exit(0)
	}

	so := MakeStackOperator(stackLimit, !(len(flag.Args()) > 0), strictMode)
	if wordsPath != "" {
		ParseWordsFile(so, wordsPath)
	}
	if so.Interactive {
		interactive(so, promptFormat)
	} else {
		nonInteractive(so, flag.Args())
	}
}
