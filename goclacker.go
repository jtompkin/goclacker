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

const Usage string = `goclacker <version>
by Josh Tompkin

usage of goclacker:
goclacker [-V] [-h] [-s] [-d] [-r] [-l] int [-c] string [-p] string [program]...
    -V, --version
        Print version information and exit.
    -h, --help
        Print usage information and exit.
    -s, --strict
        Run in strict mode: entering anything that is not a number, operator,
        or defined word will print an error instead of doing nothing.
    -d, --no-display
        Do not display stack after operations: useful if '&Nt' is in prompt.
    -r, --no-color
        Do not color output in interactive mode.
    -l, --limit int
        Provide the stack size limit. There is no limit if a negative number is
        provided. (default 8)
    -c, --config string
        Provide the path to the config file to use. Goclacker looks in the
        default locations if not provided; provide an empty string to not use
        default config files.
    -p, --prompt string
        Provide the format string for the interactive prompt. (default " &c > ")
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

const (
	DefPrompt string = " &c > "
	Version   string = "v1.4.3"
	FmtChar   byte   = '&'
	DefLimit  int    = 8
)

// Command line flags
var (
	PrintVersion, StrictMode, Display, Color bool
	ConfigPath, PromptFmt                    string
	StackLimit                               int
)

var DefConfigPaths = make([]string, 0)

type colors struct {
	out   []byte
	err   []byte
	reset []byte
}

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
	actions.Set("quit", stack.Quit)
	actions.Set("Dclip", stack.Clip)
	actions.Set("Dgrow", stack.Grow)
	actions.Set("Dfill", stack.Fill)
	so := stack.NewStackOperator(actions, StackLimit, interactive, Display, StrictMode)
	so.Words["?"] = "help"
	so.Words["randn"] = "rand * floor"
	so.Words["sqrt"] = "0.5 ^"
	so.Words["logb"] = "log swap log / -1 ^"
	so.ValWords["pi"] = math.Pi
	so.ValWords["e"] = math.E
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
	fmt.Print(string(so.ToPrint))
	return io.EOF
}

func run() (err error) {
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

	if err = so.MakePromptFunc(PromptFmt, FmtChar); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "goclacker %s\n", Version)
	err = interactive(so, Color)
	fmt.Println()
	return err
}

func main() {
	flag.BoolVar(&PrintVersion, "V", false, "")
	flag.BoolVar(&PrintVersion, "version", false, "")

	flag.BoolVar(&StrictMode, "s", false, "")
	flag.BoolVar(&StrictMode, "strict", false, "")

	flag.BoolVar(&Display, "d", false, "")
	flag.BoolVar(&Display, "no-display", false, "")

	flag.BoolVar(&Color, "r", false, "")
	flag.BoolVar(&Color, "no-color", false, "")

	flag.IntVar(&StackLimit, "l", DefLimit, "")
	flag.IntVar(&StackLimit, "limit", DefLimit, "")

	flag.StringVar(&ConfigPath, "c", "\x00", "")
	flag.StringVar(&ConfigPath, "config", "\x00", "")

	flag.StringVar(&PromptFmt, "p", "\x00", "")
	flag.StringVar(&PromptFmt, "prompt", "\x00", "")

	flag.Usage = func() { fmt.Print(strings.Replace(Usage, "<version>", Version, 1)) }
	flag.Parse()

	Display = !Display
	Color = !Color

	if err := run(); err != io.EOF {
		log.Fatal(err)
	}
}
