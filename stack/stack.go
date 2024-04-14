package stack

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Action interface {
	Call(*StackOperator) error
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
		return errors.New(fmt.Sprintf("cannot push %v, stack at capacity", f))
	}
	stk.Values = append(stk.Values, f)
	return nil
}

// Top returns the last value in stack.Values and returns an error if
// stack.Values is empty.
func (stk *Stack) Top() (float64, error) {
	if len(stk.Values) == 0 {
		return 0, errors.New("NA")
	}
	return stk.Values[len(stk.Values)-1], nil
}

// Display prints all the values in the stack
func (stk *Stack) Display() {
	fmt.Println(stk.Values)
}

func newStack(values []float64) Stack {
	stk := Stack{Values: values}
	return stk
}

// StackOperator contains map for converting string tokens into operations that
// can be called to operate on the stack.
type StackOperator struct {
	Actions map[string]Action
	Tokens  []string
	Words   map[string]string
	Stack   Stack
}

// ParseInput splits `input` into words and either starts defining a word, pushing a
// numerical value, or treating the word as a token to execute.
func (stkOp *StackOperator) ParseInput(input string) {
	split := strings.Split(input, " ")
	for i, s := range split {
		if s == "=" {
			if err := stkOp.DefWord(split[i+1:]); err != nil {
				fmt.Fprintf(os.Stderr, "definition error: %s\n", err)
			}
			return
		}
		var pErr error
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			pErr = stkOp.ParseToken(s)
		} else {
			pErr = stkOp.Stack.Push(f)
			if i == len(split)-1 {
				stkOp.Stack.Display()
			}
		}
		if pErr != nil {
			fmt.Fprintf(os.Stderr, "operation error: %s\n", pErr)
		}
	}
}

// DefWord adds a word to StackOperator.Words with the key being def[0] and the
// further values being the value. Or deletes def[0] from StackOperator.Words if
// len(def) == 1.
func (so *StackOperator) DefWord(def []string) error {
	if len(def) == 0 {
		return so.Fail("define word: '= example 2 2 +'; remove word: '= example'")
	}
	word := def[0]
	if len(def) == 1 {
		delete(so.Words, word)
		return nil
	}
	if strings.Contains("0123456789=.", string(word[0])) {
		return so.Fail(fmt.Sprintf("could not define '%s'; cannot start word with digit, '=', or '.'", word))
	}
	if _, pres := so.Actions[word]; pres {
		return so.Fail(fmt.Sprintf("could not define '%s'; cannot redifine operator", word))
	}
	so.Words[word] = strings.Join(def[1:], " ")
	return nil
}

// ParseToken determines if `token` is an Action token or defined word and
// executes it accordingly. Returns an error if the Action cannot be completed.
func (so *StackOperator) ParseToken(token string) error {
	action := so.Actions[token]
	if action == nil {
		def := so.Words[token]
		if def == "" { // input is neither a defined word nor an Action token
			return nil
		}
		so.ParseInput(def)
		return nil
	}
	sLen := len(so.Stack.Values)
	pops := action.Pops()
	var c rune
	if pops != 1 {
		c = 's'
	}
	if sLen < pops {
		return so.Fail(fmt.Sprintf("'%s' needs %d value%c in stack", token, pops, c))
	}
	if sLen-pops+action.Pushes() > cap(so.Stack.Values) {
		return so.Fail("operation would overflow stack")
	}
	if err := action.Call(so); err != nil {
		return so.Fail(fmt.Sprintf("%s", err))
	}
	return nil
}

// Fail pushes all `values` to the stack and returns an error containing `message`
func (stkOp *StackOperator) Fail(message string, values ...float64) error {
	for _, value := range values {
		stkOp.Stack.Push(value)
	}
	return errors.New(message)
}

// NewStackOperator returns a pointer to a new StackOperator, initialized to
// given arguments and a default set of defined words.
func NewStackOperator(actions map[string]Action, orderedTokens []string, maxStack int) *StackOperator {
	stkOp := StackOperator{
		Actions: actions,
		Tokens:  orderedTokens,
		Stack:   newStack(make([]float64, 0, maxStack)),
		Words: map[string]string{
			"sqrt": "0.5 ^",
			"pi":   "3.141592653589793",
			"logb": "log stash log pull /",
		},
	}
	return &stkOp
}
