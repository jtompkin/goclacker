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
	Usage string = `goclacker %s
by Josh Tompkin

usage of goclacker:
goclacker [-V] [-h] [-s] [-n] [-l] int [-c] string [-p] string [program]...
    -V, --version
        Print version information and exit.
    -h, --help
        Print usage information and exit.
    -s, --strict
        Run in strict mode: entering something that is not a number, operator,
        or defined word will print an error instead of doing nothing.
    -n, --no-display
        Do not display stack after operations. Useful if '&Nt' is in prompt.
    -l, --limit int
        stack size limit, no limit if negative (default 8)
    -c, --config string
        path to config file, looks in default locations if not provided
    -p, --prompt string
        format string for the interactive prompt (default " &c > ")
        format specifiers:
            &l  : stack limit
            &c  : current stack size
            &Nt : top N stack values
            &s  : current stash value
    [program]...
        Any positional arguments will be interpreted and executed by the
        calculator. Interactive mode will not be entered if any positional
        arguments are supplied.
`
	DefPrompt string = " &c > "
	Version   string = "v1.3.2"
	FmtChar   byte   = '&'
	DefLimit  int    = 8
)

var (
	PrintVersion, StrictMode, NoDisplay bool
	ConfigPath, PromptFmt               string
	StackLimit                          int
)

func GetStackOperator(stackLimit int, interactive bool, strict bool, noDisplay bool) *stack.StackOperator {
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

// CheckDefConfigPaths checks if files exist in any of the default config file
// paths and returns the path to the first one that exists. It returns an empty
// string if none exist.
func CheckDefConfigPaths() (path string) {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	fromHome := func(s string) string { return fmt.Sprintf("%s%c%s", home, os.PathSeparator, s) }
	defConfigPaths := []string{
		".goclacker",
		fromHome(".goclacker"),
		fromHome(".config/goclacker/goclacker.conf"),
	}
	for _, path = range defConfigPaths {
		if _, err = os.Open(path); err == nil {
			return path
		}
	}
	return ""
}

func GetConfigScanner(path string) *bufio.Scanner {
	if path == "" {
		return nil
	}
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open config file : %v\n", err)
		return nil
	}
	fmt.Fprintf(os.Stderr, "parsing config file... %s\n", path)
	return bufio.NewScanner(f)
}

func ReadPromptLine(scanner *bufio.Scanner) (promptFmt string) {
	if scanner == nil {
		return DefPrompt
	}
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "could not read prompt line : %v\n", err)
		return DefPrompt
	}
	promptFmt = scanner.Text()
	promptFmt = strings.TrimPrefix(promptFmt, `"`)
	promptFmt = strings.TrimSuffix(promptFmt, `"`)
	return promptFmt
}

func ReadProgLines(scanner *bufio.Scanner, so *stack.StackOperator) {
	if scanner == nil {
		return
	}
	for scanner.Scan() {
		if err := so.ParseInput(strings.TrimSpace(scanner.Text())); err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "could not read config file : %v\n", err)
	} else {
		fmt.Fprint(os.Stderr, "sucessfully parsed config file\n")
	}
}

func RunProgram(so *stack.StackOperator, program string) {
	err := so.ParseInput(program)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	fmt.Print(string(so.PrintBuf))
}

func ExecutePrograms(so *stack.StackOperator, programs []string) (eof error) {
	for _, s := range programs {
		if err := so.ParseInput(s); err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}
	fmt.Print(string(so.PrintBuf))
	return io.EOF
}

func run() error {
	if PrintVersion {
		fmt.Printf("goclacker %s\n", Version)
		return io.EOF
	}

	so := GetStackOperator(StackLimit, len(flag.Args()) == 0, StrictMode, NoDisplay)
	if ConfigPath == "\x00" {
		ConfigPath = CheckDefConfigPaths()
	}
	scanner := GetConfigScanner(ConfigPath)
	if s := ReadPromptLine(scanner); PromptFmt == "\x00" {
		PromptFmt = s
	}
	ReadProgLines(scanner, so)

	if !so.Interactive {
		return ExecutePrograms(so, flag.Args())
	}

	if err := so.MakePromptFunc(PromptFmt, FmtChar); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "goclacker %s\n", Version)
	return interactive(so)
}

func main() {
	flag.BoolVar(&PrintVersion, "V", false, "")
	flag.BoolVar(&PrintVersion, "version", false, "")

	flag.IntVar(&StackLimit, "l", DefLimit, "")
	flag.IntVar(&StackLimit, "limit", DefLimit, "")

	flag.BoolVar(&StrictMode, "s", false, "")
	flag.BoolVar(&StrictMode, "strict", false, "")

	flag.StringVar(&ConfigPath, "c", "\x00", "")
	flag.StringVar(&ConfigPath, "config", "\x00", "")

	flag.StringVar(&PromptFmt, "p", "\x00", "")
	flag.StringVar(&PromptFmt, "prompt", "\x00", "")

	flag.BoolVar(&NoDisplay, "n", false, "")
	flag.BoolVar(&NoDisplay, "no-display", false, "")

	flag.Usage = func() { fmt.Printf(Usage, Version) }
	flag.Parse()

	if err := run(); err != io.EOF {
		log.Fatal(err)
	}
}
