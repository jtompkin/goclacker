// Copyright 2024 Josh Tompkin
// Licensed under the MIT License

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
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
	Version   string = "v1.4.0"
	FmtChar   byte   = '&'
	DefLimit  int    = 8
)

// Command line flag
var (
	PrintVersion, StrictMode, NoDisplay bool
	ConfigPath, PromptFmt               string
	StackLimit                          int
)

var DefConfigPaths = make([]string, 0)

func GetStackOperator(interactive bool) *stack.StackOperator {
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
	actions.Set("asin", stack.Arcsine)
	actions.Set("acos", stack.Arccosine)
	actions.Set("atan", stack.Arctangent)
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
	so := stack.NewStackOperator(actions, StackLimit, interactive, NoDisplay, StrictMode)
	for _, s := range []string{
		"? help",
		"randn rand * floor",
		"sqrt 0.5 ^",
		"logb log swap log / -1 ^",
	} {
		so.DefNormWord(strings.Split(s, " "))
	}
	for _, s := range []string{
		fmt.Sprintf("pi %g", math.Pi),
		fmt.Sprintf("e %g", math.E),
	} {
		so.DefValWord(strings.Split(s, " "))
	}
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
	paths := []string{".goclacker"}
	fromHome := []string{".goclacker", ".config/goclacker/goclacker.conf"}
	for _, s := range fromHome {
		paths = append(paths, fmt.Sprintf("%s%c%s", home, os.PathSeparator, s))
	}
	for _, s := range paths {
		DefConfigPaths = append(DefConfigPaths, s)
	}
	for _, path = range DefConfigPaths {
		if _, err = os.Open(path); err == nil {
			return path
		}
	}
	return ""
}

func GetConfigScanner(path string) (scanner *bufio.Scanner, msg string) {
	if path == "" {
		return nil, ""
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Sprintf("could not open config file : %v\n", err)
	}
	return bufio.NewScanner(f), fmt.Sprintf("parsing config file... %s\n", path)
}

func ReadPromptLine(scanner *bufio.Scanner) (promptFmt string, msg string) {
	if scanner == nil {
		return DefPrompt, ""
	}
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return DefPrompt, fmt.Sprintf("could not read prompt line : %v\n", err)
	}
	promptFmt = scanner.Text()
	if promptFmt == "" {
		return DefPrompt, ""
	}
	promptFmt = strings.TrimPrefix(promptFmt, `"`)
	promptFmt = strings.TrimSuffix(promptFmt, `"`)
	return promptFmt, ""
}

func ReadProgLines(scanner *bufio.Scanner, so *stack.StackOperator) (msg string) {
	if scanner == nil {
		return
	}
	for scanner.Scan() {
		if err := so.ParseInput(strings.TrimSpace(scanner.Text())); err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Sprintf("could not read config file : %v\n", err)
	}
	return "sucessfully parsed config file\n"
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

	so := GetStackOperator(len(flag.Args()) == 0)
	if ConfigPath == "\x00" {
		ConfigPath = CheckDefConfigPaths()
	}
	scanner, msg := GetConfigScanner(ConfigPath)
	fmt.Fprint(os.Stderr, msg)

	s, msg := ReadPromptLine(scanner)
	if PromptFmt == "\x00" {
		PromptFmt = s
	}
	fmt.Fprint(os.Stderr, msg)

	msg = ReadProgLines(scanner, so)
	fmt.Fprint(os.Stderr, msg)

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
	fmt.Println()
}
