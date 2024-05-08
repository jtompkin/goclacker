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
	"golang.org/x/term"
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
        stack size limit, no limit if negative (default 8)
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
	version   string = "v1.2.1"
	fmtChar   byte   = '&'
	defLimit  int    = 8
)

func MakeStackOperator(stackLimit int, interactive bool, strict bool) *stack.StackOperator {
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
	actions.Set("stash", stack.Stash)
	actions.Set("pull", stack.Pull)
	actions.Set("clr", stack.Clear)
	actions.Set("words", stack.Words)
	actions.Set("help", stack.Help)
	actions.Set("cls", stack.ClearScreen)
	so := stack.NewStackOperator(actions, stackLimit, interactive, strict)
	so.Words = map[string]string{
		"?":     "help",
		"randn": "rand * ceil 1 -",
		"sqrt":  "0.5 ^",
		"logb":  "log swap log / -1 ^",
		"pi":    "3.141592653589793",
	}
	return so
}

func nonInteractive(so *stack.StackOperator, programs []string) {
	for _, prog := range programs {
		f := os.Stdin
		err := so.ParseInput(prog)
		if err != nil {
			f = os.Stderr
		}
		fmt.Fprint(f, so.PrintBuf)
	}
}

func interactive(so *stack.StackOperator) (err error) {
	fds := []int{int(os.Stdin.Fd()), int(os.Stderr.Fd())}
	for _, fd := range fds {
		state, err := term.MakeRaw(fd)
		if err != nil {
			return err
		}
		defer term.Restore(fd, state)
	}

	it := term.NewTerminal(os.Stdin, so.Prompt())
	et := term.NewTerminal(os.Stderr, "")
	for {
		line, err := it.ReadLine()
		if err != nil {
            it.SetPrompt("")
            it.Write([]byte{})
			return err
		}
		err = so.ParseInput(line)
		it.Write(so.PrintBuf)
		if err != nil {
			et.Write([]byte(err.Error()))
		}
		it.SetPrompt(so.Prompt())
	}
}

func start(so *stack.StackOperator, progs []string) (err error) {
	if so.Interactive {
		err = interactive(so)
		return err
	}
	nonInteractive(so, progs)
	return err
}

func configure(so *stack.StackOperator, path string, promptFmt string) (err error) {
	gavePrompt := true
	if promptFmt == "\x00" {
		promptFmt = defPrompt
		gavePrompt = false
	}
	if path == "" {
		if err := so.MakePromptFunc(promptFmt, fmtChar); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}
		return nil
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
	promptLine := scanner.Text()
	if len(promptLine) > 0 && !gavePrompt {
		promptLine = strings.TrimPrefix(promptLine, `"`)
		promptLine = strings.TrimSuffix(promptLine, `"`)
		if err := so.MakePromptFunc(promptLine, fmtChar); err != nil {
			return err
		}
		fmt.Print("sucessfully parsed prompt from file...\n")
	} else {
		if err := so.MakePromptFunc(promptFmt, fmtChar); err != nil {
			return err
		}
	}
	var failed bool
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if _, err := so.DefWord(strings.Split(line, " ")); err != nil {
			fmt.Fprintf(os.Stderr, "definition error: %s", err.Error())
			failed = true
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if failed {
		fmt.Fprint(os.Stderr, "enter 'help' to see list of operators that cannot be used as words...\n")
	}
	fmt.Print("sucessfully parsed config file\n")
	return nil
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
	flag.StringVar(&promptFormat, "p", "\x00", "")
	flag.StringVar(&promptFormat, "prompt", "\x00", "")

	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	if printVersion {
		fmt.Printf("goclacker %s\n", version)
		os.Exit(0)
	}

	so := MakeStackOperator(stackLimit, len(flag.Args()) == 0, strictMode)
	if err := configure(so, configPath, promptFormat); err != nil {
		log.Fatal(err)
	}
	if err := start(so, flag.Args()); err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
