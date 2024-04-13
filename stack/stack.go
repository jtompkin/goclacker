package stack

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

func NewStack(values []float64) *Stack {
	stk := Stack{}
	stk.values = values
	return &stk
}

// Operation contains a function to operate on a StackOperator instance. `Pops`
// and `Pushes` represent the number of values the functions pops from and
// pushes to the stack.
type Operation struct {
	action func(*StackOperator) error
	Pops   int
	Pushes int
	Help   string
}

func NewOperation(action func(*StackOperator) error, pops int, pushes int, help string) *Operation {
	op := Operation{action: action, Pops: pops, Pushes: pushes, Help: help}
	return &op
}

// StackOperator contains map for converting string tokens into operations that
// can be called to operate on the stack.
type StackOperator struct {
	operators map[string]*Operation
	words     map[string]string
	Stack     Stack
}

func (stkOp *StackOperator) ParseInput(input string) {
	split := strings.Split(input, " ")
	for i, s := range split {
		if s == "=" {
			dErr := stkOp.DefWord(split[i+1:])
			if dErr != nil {
				fmt.Fprintf(os.Stderr, "definition error: %s\n", dErr)
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
	if _, pres := stkOp.operators[word]; pres {
		return stkOp.Fail(fmt.Sprintf("could not define '%s'; cannot redifine operator", word))
	}
	stkOp.words[word] = strings.Join(def[1:], " ")
	return nil
}

func (stkOp *StackOperator) ParseToken(token string) error {
	op := stkOp.operators[token]
	if op == nil {
		def := stkOp.words[token]
		if def == "" {
			return nil
		}
		stkOp.ParseInput(def)
		return nil
	}
	stkLen := stkOp.Stack.Len()
	if stkLen < op.Pops {
		return stkOp.Fail(fmt.Sprintf("operation needs %d values in stack", op.Pops))
	}
	if stkLen-op.Pops+op.Pushes > stkOp.Stack.Cap() {
		return stkOp.Fail("operation would overflow stack")
	}
	if err := op.action(stkOp); err != nil {
		return stkOp.Fail(fmt.Sprintf("%s", err))
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

func (stkOp *StackOperator) PrintHelp() error {
	f, err := os.Open("data/help.tab")
	if err != nil {
		return stkOp.Fail("could not read help file")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), "\t")
		if len(split) == 2 {
			fmt.Printf("operator: %s\t\"%s\"\n", split[0], split[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return stkOp.Fail("could not read help file")
	}
	return nil
}

func NewStackOperator(operators map[string]*Operation, maxStack int) *StackOperator {
	stkOp := StackOperator{
		operators: operators,
		words:     map[string]string{"sqrt": "0.5 ^", "pi": "3.141592653589793", "logb": "log stash log pull /"},
		Stack:     *NewStack(make([]float64, 0, maxStack)),
	}
	return &stkOp
}
