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

const usage string = `usage of goclacker:
goclacker [-V] [-h] [-s] [-l] int [-w] string [-p] string [program...]
    -V, --version
        Print version information and exit.
    -h, --help
        Print usage information and exit.
    -s, --strict
        Run in strict mode: entering something that is not a number, operator,
        or defined word will return an error instead of doing nothing.
    -l, --limit int
        stack size limit, must be non-negative (default 8)
    -c, --config string
        path to config file
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
	version   string = "v1.1.2"
	fmtChar   byte   = '&'
	defLimit  int    = 8
)

func MakeStackOperator(stackLimit int, interactive bool, strict bool) *stack.StackOperator {
	actions := stack.NewOrderedMap[string, *stack.Action](64)
	actions.Set("+", stack.Add)
	actions.Set("-", stack.Subtract)
	actions.Set("*", stack.Multiply)
	actions.Set("/", stack.Divide)
	actions.Set("%", stack.Modulo)
	actions.Set("^", stack.Power)
	actions.Set("!", stack.Factorial)
	actions.Set("log", stack.Log)
	actions.Set("ln", stack.Ln)
	actions.Set("rad", stack.Radians)
	actions.Set("deg", stack.Degrees)
	actions.Set("sin", stack.Sine)
	actions.Set("cos", stack.Cosine)
	actions.Set("tan", stack.Tangent)
	actions.Set("floor", stack.Floor)
	actions.Set("ceil", stack.Ceiling)
	actions.Set("round", stack.Round)
	actions.Set("rand", stack.Random)
	actions.Set(".", stack.Display)
	actions.Set(",", stack.Pop)
	actions.Set("stash", stack.Stash)
	actions.Set("pull", stack.Pull)
	actions.Set("clear", stack.Clear)
	actions.Set("words", stack.Words)
	actions.Set("help", stack.Help)
	actions.Set("cls", stack.ClearScreen)
	so := stack.NewStackOperator(actions, stackLimit, interactive, strict)
	so.Words = map[string]string{
		"randn": "rand * ceil 1 -",
		"sqrt":  "0.5 ^",
		"logb":  "log stash log pull /",
		"pi":    "3.141592653589793",
	}
	return so
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

func interactive(so *stack.StackOperator) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(so.Prompt())
	for scanner.Scan() {
		if s, err := so.ParseInput(scanner.Text()); err != nil {
			fmt.Fprint(os.Stderr, err)
		} else {
			fmt.Print(s)
		}
		fmt.Print(so.Prompt())
	}
	fmt.Print("\n")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func configure(so *stack.StackOperator, path string, promptFmt string) {
	gavePrompt := true
	if promptFmt == "\000" {
		promptFmt = defPrompt
		gavePrompt = false
	}
	if path == "" {
		if err := so.MakePromptFunc(promptFmt, fmtChar); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}
		return
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	promptLine := scanner.Text()
	var failed bool
	if len(promptLine) > 0 && !gavePrompt {
		promptLine = strings.TrimPrefix(promptLine, `"`)
		promptLine = strings.TrimSuffix(promptLine, `"`)
		if err := so.MakePromptFunc(promptLine, fmtChar); err != nil {
			log.Fatal(err)
		}
		fmt.Print("sucessfully parsed prompt from file...\n")
	} else {
		if err := so.MakePromptFunc(promptFmt, fmtChar); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if _, err := so.DefWord(strings.Split(line, " ")); err != nil {
			fmt.Fprintf(os.Stderr, "definition error: %s", err.Error())
			failed = true
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	if failed {
		fmt.Fprint(os.Stderr, "enter 'help' to see list of operators that cannot be used as words...\n")
	}
	fmt.Print("sucessfully parsed config file\n")
}

func main() {
	var printVersion bool
	flag.BoolVar(&printVersion, "V", false, "")
	flag.BoolVar(&printVersion, "version", false, "")
	var stackLimit int
	flag.IntVar(&stackLimit, "l", defLimit, "")
	flag.IntVar(&stackLimit, "limit", defLimit, "")
	var strictMode bool
	flag.BoolVar(&strictMode, "s", false, "")
	flag.BoolVar(&strictMode, "strict", false, "")
	var configPath string
	flag.StringVar(&configPath, "c", "", "")
	flag.StringVar(&configPath, "config", "", "")
	var promptFormat string
	flag.StringVar(&promptFormat, "p", "\000", "")
	flag.StringVar(&promptFormat, "prompt", "\000", "")

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
	configure(so, configPath, promptFormat)
	if so.Interactive {
		interactive(so)
	} else {
		nonInteractive(so, flag.Args())
	}
}
