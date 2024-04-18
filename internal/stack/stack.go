package stack

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const Suffix string = "\n"

type Action interface {
	Call(*StackOperator) (string, error)
	Pops() int
	Pushes() int
	Help() string
}

// Stack contains a slice of values and methods to operate on that slice.
type Stack struct {
	Values []float64
	Stash  float64
}

// Pop removes the last value in stack.Values and returns the value removed.
func (stk *Stack) Pop() float64 {
	n := len(stk.Values) - 1
	f := stk.Values[n]
	stk.Values = stk.Values[:n]
	return f
}

// Push attempts to append f to stack.Values and returns an error if the stack
// is at capacity.
func (stk *Stack) Push(f float64) error {
	if len(stk.Values)+1 > cap(stk.Values) {
		return errors.New(fmt.Sprintf("cannot push %v, stack at capacity (%d)", f, cap(stk.Values)))
	}
	stk.Values = append(stk.Values, f)
	return nil
}

// Display prints all the values in the stack
func (stk *Stack) Display(fancy bool) string {
	var s string
	if fancy {
		s = "[ "
		for _, f := range stk.Values {
			s += fmt.Sprint(f, " ")
		}
		return s + "]" + Suffix
	}
	for _, f := range stk.Values {
		s += fmt.Sprint(f, " ")
	}
	return s + Suffix
}

func newStack(values []float64) Stack {
	stk := Stack{Values: values}
	return stk
}

func rStrip(s string, chars string) string {
	if len(s) == 0 {
		return s
	}
	for i := len(s) - 1; strings.ContainsAny(string(s[i]), chars) && i > 0; i-- {
		s = s[:i]
	}
	return s
}

// StackOperator contains map for converting string tokens into operations that
// can be called to operate on the stack.
type StackOperator struct {
	Actions     map[string]Action
	Tokens      []string
	Words       map[string]string
	Stack       Stack
	formatters  map[byte]func(*StackOperator) string
	Prompt      func() string
	Interactive bool
}

// ParseInput splits `input` into words and either starts defining a word, pushing a
// numerical value, or treating the word as a token to execute.
func (so *StackOperator) ParseInput(input string) (string, error) {
	var rs string
	input = rStrip(input, " \t")
	split := strings.Split(input, " ")
	for i, s := range split {
		if s == "=" {
			if s, err := so.DefWord(split[i+1:]); err != nil {
				return "", errors.New(fmt.Sprintf("definition error: %s\n", err))
			} else {
				return s + Suffix, nil
			}
		}
		var pErr error
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			rs, pErr = so.ParseToken(s)
		} else {
			pErr = so.Stack.Push(f)
			if i == len(split)-1 {
				return so.Stack.Display(so.Interactive), nil
			}
		}
		if pErr != nil {
			return "", errors.New(fmt.Sprintf("operation error: %s\n", pErr))
		}
	}
	return rs, nil
}

// DefWord adds a word to StackOperator.Words with the key being def[0] and the
// further values being the value. Or deletes def[0] from StackOperator.Words if
// len(def) == 1.
func (so *StackOperator) DefWord(def []string) (string, error) {
	if len(def) == 0 {
		return "", errors.New("define word: '= example 2 2 +'; remove word: '= example'")
	}
	noEmpty := make([]string, 0, len(def))
	for _, s := range def {
		if s != "" {
			noEmpty = append(noEmpty, s)
		}
	}
	word := noEmpty[0]
	if len(noEmpty) == 1 {
		delete(so.Words, word)
		return fmt.Sprintf("deleted word: %s", word), nil
	}
	if strings.Contains("0123456789=.", string(word[0])) {
		return "", errors.New(fmt.Sprintf("could not define '%s'; cannot start word with digit, '=', or '.'", word))
	}
	if _, pres := so.Actions[word]; pres {
		return "", errors.New(fmt.Sprintf("could not define '%s'; cannot redifine operator", word))
	}
	s := strings.Join(noEmpty[1:], " ")
	so.Words[word] = s
	return fmt.Sprintf(`defined word: "%s" with value: "%s"`, word, s), nil
}

// ParseToken determines if `token` is an Action token or defined word and
// executes it accordingly. Returns an error if the Action cannot be completed.
func (so *StackOperator) ParseToken(token string) (string, error) {
	action := so.Actions[token]
	if action == nil {
		def := so.Words[token]
		if def == "" { // input is neither a defined word nor an Action token
			return "", nil
		}
		return so.ParseInput(def)
	}
	sLen := len(so.Stack.Values)
	pops := action.Pops()
	var c rune
	if pops != 1 {
		c = 's'
	}
	if sLen < pops {
		return "", errors.New(fmt.Sprintf("'%s' needs %d value%c in stack", token, pops, c))
	}
	if sLen-pops+action.Pushes() > cap(so.Stack.Values) {
		return "", errors.New("operation would overflow stack")
	}
	s, err := action.Call(so)
	if err != nil {
		return "", errors.New(fmt.Sprintf("%s", err))
	}
	return s, nil
}

// Fail pushes all `values` to the stack and returns an error containing `message`
func (so *StackOperator) Fail(message string, values ...float64) error {
	for _, f := range values {
		so.Stack.Push(f)
	}
	if so.Interactive {
		fmt.Printf("%s", so.Stack.Display(true))
	}
	return errors.New(message)
}

// MakePromptFunc sets the StackOperator.prompt value that will execute any
// functions needed to build the prompt and will return a new string every time
// it is called
func (so *StackOperator) MakePromptFunc(format string, fmtChar byte) {
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
		log.Fatal("Something done gone wrong with the prompt")
	}

	so.Prompt = func() string {
		sb := new(strings.Builder)
		for i, f := range promptFuncs {
			sb.WriteString(fmt.Sprintf(promptSplit[i], f(so)))
		}
		sb.WriteString(promptSplit[len(promptSplit)-1])
		return sb.String()
	}
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
// given arguments and a default set of defined words.
func NewStackOperator(actions map[string]Action, orderedTokens []string, maxStack int, interactive bool) *StackOperator {
	stkOp := StackOperator{
		Actions:     actions,
		Tokens:      orderedTokens,
		Stack:       newStack(make([]float64, 0, maxStack)),
		Interactive: interactive,
		Words: map[string]string{
			"sqrt": "0.5 ^",
			"pi":   "3.141592653589793",
			"logb": "log stash log pull /",
		},
		formatters: map[byte]func(*StackOperator) string{
			'l': (*StackOperator).promptCap,
			't': (*StackOperator).promptTop,
			'c': (*StackOperator).promptLen,
			's': (*StackOperator).promptStash,
		},
	}
	return &stkOp
}
