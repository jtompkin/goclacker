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
	GetPops() int
	GetPushes() int
	GetHelp() string
}

// Stack contains a slice of values and methods to operate on that slice.
type Stack struct {
	values []float64
	stash  float64
}

func (stk *Stack) Pop() float64 {
	n := len(stk.values) - 1
	f := stk.values[n]
	stk.values = stk.values[:n]
	return f
}

func (stk *Stack) Push(f float64) error {
	if len(stk.values)+1 > cap(stk.values) {
		return errors.New(fmt.Sprintf("cannot push %f, stack at capacity", f))
	}
	stk.values = append(stk.values, f)
	return nil
}

func (stk *Stack) Display() {
	fmt.Println(stk.values)
}

func (stk *Stack) SetValues(values []float64) {
	stk.values = values
}

func (stk *Stack) Len() int {
	return len(stk.values)
}

func (stk *Stack) Cap() int {
	return cap(stk.values)
}

func (stk *Stack) SetStash(value float64) {
	stk.stash = value
}

func (stk *Stack) GetStash() float64 {
	return stk.stash
}

func newStack(values []float64) Stack {
	stk := Stack{values: values}
	return stk
}

// StackOperator contains map for converting string tokens into operations that
// can be called to operate on the stack.
type StackOperator struct {
	actions map[string]Action
	Tokens  *[]string
	words   map[string]string
	Stack   Stack
}

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

func (stkOp *StackOperator) DefWord(def []string) error {
	if len(def) == 0 {
		return nil
	}
	word := def[0]
	if len(def) == 1 {
		delete(stkOp.words, word)
		return nil
	}
	if strings.Contains("0123456789=.", string(word[0])) {
		return stkOp.Fail(fmt.Sprintf("could not define '%s'; cannot start word with digit, '=', or '.'", word))
	}
	if _, pres := stkOp.actions[word]; pres {
		return stkOp.Fail(fmt.Sprintf("could not define '%s'; cannot redifine operator", word))
	}
	stkOp.words[word] = strings.Join(def[1:], " ")
	return nil
}

func (so *StackOperator) ParseToken(token string) error {
	action := so.actions[token]
	if action == nil {
		def := so.words[token]
		if def == "" {
			return nil
		}
		so.ParseInput(def)
		return nil
	}
	sLen := so.Stack.Len()
	pops := action.GetPops()
	var c rune
	if pops != 1 {
		c = 's'
	}
	if sLen < pops {
		return so.Fail(fmt.Sprintf("'%s' needs %d value%c in stack", token, pops, c))
	}
	if sLen-pops+action.GetPushes() > so.Stack.Cap() {
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

func (stkOp *StackOperator) GetWords() map[string]string {
	return stkOp.words
}

func (so *StackOperator) GetActions() map[string]Action {
	return so.actions
}

func (stkOp *StackOperator) PrintHelp() error {
	return nil
}

func NewStackOperator(actions map[string]Action, orderedTokens *[]string, maxStack int) *StackOperator {
	stkOp := StackOperator{
		actions: actions,
		Tokens:  orderedTokens,
		words:   map[string]string{"sqrt": "0.5 ^", "pi": "3.141592653589793", "logb": "log stash log pull /"},
		Stack:   newStack(make([]float64, 0, maxStack)),
	}
	return &stkOp
}
