package stack

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Contains a slice of values and methods to operate on that slice.
type Stack struct {
	values []float64
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

func NewStack(values []float64) *Stack {
	stk := Stack{}
	stk.values = values
	return &stk
}

// Contains a function to operate on a StackOperator instance. `Pops` and
// `Pushes` represent the number of values the functions pops from and pushes to
// the stack.
type Operation struct {
	action func(*StackOperator) error
	Pops   int
	Pushes int
}

func NewOperation(action func(*StackOperator) error, pops int, pushes int) *Operation {
	op := Operation{action: action, Pops: pops, Pushes: pushes}
	return &op
}

// Contains map for converting string tokens into operations that can be called
// to operate on the stack.
type StackOperator struct {
	operators map[string]*Operation
	Stack     Stack
	MaxStack  int
}

func (stkOp *StackOperator) ParseInput(input string) {
	split := strings.Split(input, " ")
	for i, s := range split {
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

func (stkOp *StackOperator) ParseToken(token string) error {
	op := stkOp.operators[token]
	if op == nil {
		return nil
	}
	stkLen := stkOp.Stack.Len()
	if stkLen < op.Pops {
		return stkOp.Fail(fmt.Sprintf("operation needs %d values in stack", op.Pops))
	}
	if stkLen-op.Pops+op.Pushes > stkOp.MaxStack {
		return stkOp.Fail("operation would overflow stack")
	}
	if err := op.action(stkOp); err != nil {
		return stkOp.Fail(fmt.Sprintf("%s", err))
	}
	return nil
}

// Pushes all `values` to the stack and returns an error containing `message`
func (stkOp *StackOperator) Fail(message string, values ...float64) error {
	for _, value := range values {
		stkOp.Stack.Push(value)
	}
	return errors.New(message)
}

func NewStackOperator(operators map[string]*Operation, maxStack int) *StackOperator {
	stkOp := StackOperator{
		operators: operators,
		MaxStack:  maxStack,
		Stack:     *NewStack(make([]float64, 0, maxStack)),
	}
	return &stkOp
}
