package stack

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const Suffix string = "\n"

// Stack contains a slice of float64 and methods to operate on that slice.
type Stack struct {
	Values     []float64
	Stash      float64
	displayFmt string
}

// Pop removes the last value in Stack.Values and returns the value removed.
func (stk *Stack) Pop() float64 {
	n := len(stk.Values) - 1
	f := stk.Values[n]
	stk.Values = stk.Values[:n]
	return f
}

// Push attempts to append f to Stack.Values and returns an error if the stack
// is at capacity.
func (stk *Stack) Push(f float64) error {
	if len(stk.Values)+1 > cap(stk.Values) {
		return errors.New(fmt.Sprintf("cannot push %v, stack at capacity (%d)\n", f, cap(stk.Values)))
	}
	stk.Values = append(stk.Values, f)
	return nil
}

// Display returns a string of all values in the stack according to
// Stack.displayFmt
func (stk *Stack) Display() string {
	ss := make([]string, len(stk.Values))
	for i, f := range stk.Values {
		ss[i] = fmt.Sprint(f)
	}
	s := strings.Join(ss, " ")
	return fmt.Sprintf(stk.displayFmt, s)
}

func newStack(values []float64, displayFmt string) *Stack {
	return &Stack{Values: values, displayFmt: displayFmt}
}

// StackOperator contains a map for converting string tokens into operations
// that can be called to operate on the stack.
type StackOperator struct {
	Actions     *OrderedMap[string, *Action]
	Words       map[string]string
	Stack       *Stack
	Interactive bool
	formatters  map[byte]func(*StackOperator) string
	Prompt      func() string
	notFound    func(string) error
}

// ParseInput splits input into words and either begins defining a calculator
// word or parses each word as a token. It stops parsing input if it defines a
// word or receives an error during token parsing, and returns the string and
// error recieved. If it parses the entire input without error, it will return
// the last string received and nil.
func (so *StackOperator) ParseInput(input string) (string, error) {
	var rs string
	var err error
	input = strings.TrimSpace(input)
	split := strings.Split(input, " ")
	for i, s := range split {
		if s == "=" {
			return so.DefWord(split[i+1:])
		}
		rs, err = so.parseToken(s)
		if err != nil {
			return "", err
		}
	}
	return rs, nil
}

// DefWord adds a word to StackOperator.Words with the key being def[0] and the
// value being the rest of the slice. It deletes def[0] from StackOperator.Words
// if len(def) == 1. It returns a string and nil if the operator was successful
// and an empty string and an error if not.
func (so *StackOperator) DefWord(def []string) (string, error) {
	if len(def) == 0 {
		return "", errors.New(fmt.Sprintf(`define word: "= example 2 2 +"; remove word: "= example"%s`, Suffix))
	}
	noEmpty := make([]string, 0, len(def))
	for _, s := range def {
		if s != "" {
			noEmpty = append(noEmpty, s)
		}
	}
	if len(noEmpty) == 0 {
		return "", nil
	}
	word := noEmpty[0]
	if strings.Contains("0123456789=.", string(word[0])) {
		return "", errors.New(fmt.Sprintf("could not define %s : cannot start word with digit, '=', or '.'%s", word, Suffix))
	}
	if _, present := so.Actions.Get(word); present {
		return "", errors.New(fmt.Sprintf("could not define %s : cannot redifine operator%s", word, Suffix))
	}
	if len(noEmpty) == 1 {
		if _, present := so.Words[word]; !present {
			return "", errors.New(fmt.Sprintf("could not delete %s : not defined%s", word, Suffix))
		}
		delete(so.Words, word)
		return fmt.Sprintf("deleted %s%s", word, Suffix), nil
	}
	s := strings.Join(noEmpty[1:], " ")
	so.Words[word] = s
	return fmt.Sprintf(`defined %s : %s%s`, word, s, Suffix), nil
}

// parseToken parses token that should be one word and either pushes it to the
// stack as a number or executes it as a token. It returns the result of
// executing the token or the return value of Stack.Display and the error value
// from pushing token to the stack.
func (so *StackOperator) parseToken(token string) (string, error) {
	f, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return so.ExecuteToken(token)
	}
	err = so.Stack.Push(f)
	return so.Stack.Display(), err
}

// ExecuteToken determines if `token` is an Action token or defined word and
// executes it accordingly. Returns the string and error from doing what it
// needs to do.
func (so *StackOperator) ExecuteToken(token string) (string, error) {
	prefix := "operation error: "
	p, present := so.Actions.Get(token)
	if !present {
		def := so.Words[token]
		if def == "" {
			return "", so.notFound(token)
		}
		return so.ParseInput(def)
	}
	action := p.Value
	stkLen := len(so.Stack.Values)
	pops := action.Pops
	var c byte
	if pops != 1 {
		c = 's'
	}
	if stkLen < pops {
		return "", errors.New(fmt.Sprintf("%s'%s' needs %d value%c in stack%s", prefix, token, pops, c, Suffix))
	}
	if stkLen-pops+action.Pushes > cap(so.Stack.Values) {
		return "", errors.New(fmt.Sprintf("%s'%s' would overflow stack%s", prefix, token, Suffix))
	}
	return action.Call(so)
}

// Fail pushes all `values` to the stack and returns an error containing
// `message`. It also prints Stack.Display if the StackOperator is interactive
func (so *StackOperator) Fail(message string, values ...float64) error {
	for _, f := range values {
		so.Stack.Push(f)
	}
	if so.Interactive {
		fmt.Print(so.Stack.Display())
	}
	return errors.New(fmt.Sprintf("operation error: %s%s", message, Suffix))
}

// MakePromptFunc sets the StackOperator.prompt value that will execute any
// functions needed to build the prompt and will return a new string every time
// it is called. This function returns an error if it could not make the prompt
// function (should always return nil).
func (so *StackOperator) MakePromptFunc(format string, fmtChar byte) error {
	promptFuncs := make([]func(*StackOperator) string, 0, strings.Count(format, string(fmtChar)))
	promptFmt := []byte(format)
	for i := 0; i < len(format)-1; i++ {
		if format[i] == fmtChar {
			f := so.formatters[format[i+1]]
			if f != nil {
				promptFuncs = append(promptFuncs, f)
				promptFmt[i] = '%'
				promptFmt[i+1] = 's'
			}
			i++
		}
	}
	promptSplit := strings.SplitAfter(string(promptFmt), "%s")
	if len(promptSplit) != len(promptFuncs)+1 {
		return errors.New(fmt.Sprintf("Something done gone wrong with the prompt...\nspecifyers: %d, functions: %d", len(promptSplit)-1, len(promptFuncs)))
	}

	so.Prompt = func() string {
		sb := new(strings.Builder)
		for i, f := range promptFuncs {
			sb.WriteString(fmt.Sprintf(promptSplit[i], f(so)))
		}
		sb.WriteString(promptSplit[len(promptSplit)-1])
		return sb.String()
	}
	return nil
}

func (so *StackOperator) promptCap() string {
	return fmt.Sprint(cap(so.Stack.Values))
}

func (so *StackOperator) promptTop() string {
	stkLen := len(so.Stack.Values)
	if stkLen == 0 {
		return "NA"
	}
	return fmt.Sprint(so.Stack.Values[stkLen-1])
}

func (so *StackOperator) promptLen() string {
	return fmt.Sprint(len(so.Stack.Values))
}

func (so *StackOperator) promptStash() string {
	return fmt.Sprint(so.Stack.Stash)
}

// NewStackOperator returns a pointer to a new StackOperator, initialized to
// given arguments and a default set of defined words and formatters.
func NewStackOperator(actions *OrderedMap[string, *Action], maxStack int, interactive bool, strict bool) *StackOperator {
	var displayFmt string
	if interactive {
		displayFmt = "[ %s ]" + Suffix
	} else {
		displayFmt = "%s" + Suffix
	}
	var notFound func(string) error
	if strict {
		notFound = func(s string) error { return errors.New(fmt.Sprintf("command not found: %q\n", s)) }
	} else {
		notFound = func(s string) error { return nil }
	}
	return &StackOperator{
		Stack:       newStack(make([]float64, 0, maxStack), displayFmt),
		Actions:     actions,
		notFound:    notFound,
		Interactive: interactive,
		Words: map[string]string{
			"sqrt":  "0.5 ^",
			"pi":    "3.141592653589793",
			"logb":  "log stash log pull /",
			"randn": "rand * ceil 1 -",
		},
		formatters: map[byte]func(*StackOperator) string{
			'l': (*StackOperator).promptCap,
			't': (*StackOperator).promptTop,
			'c': (*StackOperator).promptLen,
			's': (*StackOperator).promptStash,
		},
	}
}
