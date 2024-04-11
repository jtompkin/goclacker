package stack

import "fmt"

type Stack struct {
	values []float32
}

func (stk *Stack) pop() float32 {
	n := len(stk.values) - 1
	f := stk.values[n]
	stk.values = stk.values[:n]
	return f
}

func (stk *Stack) push(f float32) {
	stk.values = append(stk.values, f)
}

func (stk *Stack) display() {
	fmt.Println(stk.values)
}

func NewStack(values []float32) *Stack {
	stk := Stack{}
	stk.values = values
	return &stk
}

type Operation struct {
	action func(*StackOperator)
	pops   int
	pushes int
}

func NewOperation(action func(*StackOperator), pops int, pushes int) *Operation {
	op := Operation{}
	op.action = action
	op.pops = pops
	op.pushes = pushes
	return &op
}

type StackOperator struct {
	operators map[string]*Operation
	stack     Stack
}

func (stkOp *StackOperator) ParseInput(input string) {

}

func (stkOp *StackOperator) ParseToken(token string) {
	op := stkOp.operators[token]
	if op == nil {
		return
	}
	stkLen := len(stkOp.stack.values)
	if stkLen < op.pops {
		return
	}
	if stkLen-op.pops+op.pushes > cap(stkOp.stack.values) {
		return
	}
	op.action(stkOp)
}

func NewStackOperator(operators map[string]*Operation, maxStack int) *StackOperator {
	stkOp := StackOperator{}
	stkOp.operators = operators
	stkOp.stack = *NewStack(make([]float32, 0, maxStack))
	return &stkOp
}
