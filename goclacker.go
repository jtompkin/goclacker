// Copyright 2024 Josh Tompkin
// Licensed under the MIT License

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jtompkin/goclacker/internal/stack"
)

const (
	usage string = `goclacker %s
by Josh Tompkin

usage of goclacker:
goclacker [-V] [-h] [-s] [-n] [-l] int [-c] string [-p] string [program...]
    -V, --version
        Print version information and exit.
    -h, --help
        Print usage information and exit.
    -s, --strict
        Run in strict mode: entering something that is not a number, operator,
        or defined word will return an error instead of doing nothing.
    -n, --no-display
        Do not display stack after operations. Useful if '&Nt' is in prompt.
    -l, --limit int
        stack size limit, no limit if negative (default 8)
    -c, --config string
        path to config file
    -p, --prompt string
        format string for the interactive prompt (default " &c > ")
        format specifiers:
            &l  : stack limit
            &c  : current stack size
            &Nt : top N stack values
            &s  : current stash value
    [program...]
        Any positional arguments will be interpreted and executed by the
        calculator. Interactive mode will not be entered if any positional
        arguments are supplied.
`
	defPrompt string = " &c > "
	version   string = "v1.3.1"
	fmtChar   byte   = '&'
	defLimit  int    = 8
)

func MakeStackOperator(stackLimit int, interactive bool, strict bool, noDisplay bool) *stack.StackOperator {
	actions := stack.NewOrderedMap[string, *stack.Action]()
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
	actions.Set("swap", stack.Swap)
	actions.Set("froll", stack.Froll)
	actions.Set("rroll", stack.Rroll)
	actions.Set("sum", stack.Sum)
	actions.Set("avg", stack.Average)
	actions.Set("stash", stack.Stash)
	actions.Set("pull", stack.Pull)
	actions.Set("clr", stack.Clear)
	actions.Set("words", stack.Words)
	actions.Set("help", stack.Help)
	actions.Set("cls", stack.ClearScreen)
	actions.Set("Dclip", stack.Clip)
	actions.Set("Dgrow", stack.Grow)
	actions.Set("Dfill", stack.Fill)
	so := stack.NewStackOperator(actions, stackLimit, interactive, noDisplay, strict)
	split := func(s string) []string { return strings.Split(s, " ") }
	so.DefWord(split("? help"))
	so.DefWord(split("randn rand * floor"))
	so.DefWord(split("sqrt 0.5 ^"))
	so.DefWord(split("logb log swap log / -1 ^"))
	so.DefWord(split("pi 3.141592653589793"))
	return so
}

func nonInteractive(so *stack.StackOperator, programs []string) {
	var f io.Writer
	for _, prog := range programs {
		f = os.Stdout
		err := so.ParseInput(prog)
		if err != nil {
			f = os.Stderr
			so.PrintBuf = []byte(err.Error())
		}
	}
	fmt.Fprint(f, string(so.PrintBuf))
}

func Interactive(so *stack.StackOperator) (err error) {
	fmt.Printf("goclacker %s by Josh Tompkin\n", version)
	return interactive(so)
}

func start(so *stack.StackOperator, progs []string) (err error) {
	if so.Interactive {
		err = Interactive(so)
	} else {
		nonInteractive(so, progs)
	}
	return err
}

func configure(so *stack.StackOperator, path string, promptFmt string) (err error) {
	gavePrompt := true
	if promptFmt == "\x00" {
		promptFmt = defPrompt
		gavePrompt = false
	}
	if path == "" {
		return so.MakePromptFunc(promptFmt, fmtChar)
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "parsing config file... %s\n", path)
	promptLine := scanner.Text()
	if len(promptLine) > 0 && !gavePrompt {
		promptLine = strings.TrimPrefix(promptLine, `"`)
		promptLine = strings.TrimSuffix(promptLine, `"`)
		err = so.MakePromptFunc(promptLine, fmtChar)
	} else {
		err = so.MakePromptFunc(promptFmt, fmtChar)
	}
	if err != nil {
		return err
	}
	var failed bool
	for scanner.Scan() {
		if err := so.ParseInput(strings.TrimSpace(scanner.Text())); err != nil {
			fmt.Fprint(os.Stderr, err)
			failed = true
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if !failed {
		fmt.Fprint(os.Stderr, "sucessfully parsed config file\n")
	}
	return nil
}

// main parses command line arguments and starts the program.
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
	flag.StringVar(&promptFormat, "p", "\x00", "")
	flag.StringVar(&promptFormat, "prompt", "\x00", "")
	var noDisplay bool
	flag.BoolVar(&noDisplay, "n", false, "")
	flag.BoolVar(&noDisplay, "no-display", false, "")

	flag.Usage = func() { fmt.Printf(usage, version) }
	flag.Parse()

	if printVersion {
		fmt.Printf("goclacker %s\n", version)
		return
	}

	so := MakeStackOperator(stackLimit, len(flag.Args()) == 0, strictMode, noDisplay)
	if err := configure(so, configPath, promptFormat); err != nil {
		log.Fatal(err)
	}
	if err := start(so, flag.Args()); err != io.EOF {
		log.Fatal(err)
	}
}
