package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type stack struct {
	values []float64
}

func (stk *stack) pop() float64 {
	n := len(stk.values) - 1
	f := stk.values[n]
	stk.values = stk.values[:n]
	return f
}

func (stk *stack) push(f float64) {
	stk.values = append(stk.values, f)
}

func (stk *stack) display() {
	fmt.Println(stk.values)
}

func newStack(values []float64) *stack {
	stk := stack{}
	stk.values = values
	return &stk
}

type stackOperator struct {
	operators map[string]*operation
	stack     stack
}

func (stkOp *stackOperator) CallToken(token string) {
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

func NewStackOperator(operators map[string]*operation) *stackOperator {
	so := stackOperator{}
	so.operators = operators
	so.stack = *newStack(make([]float64, 0, 16))
	return &so
}

type operation struct {
	action func(*stackOperator)
	pops   int
	pushes int
}

func NewOperation(action func(*stackOperator), pops int, pushes int) *operation {
	op := operation{}
	op.action = action
	op.pops = pops
	op.pushes = pushes
	return &op
}

func add(stkOp *stackOperator) {
	stkOp.stack.push(stkOp.stack.pop() + stkOp.stack.pop())
}

func subtract(stkOp *stackOperator) {
	x := stkOp.stack.pop()
	y := stkOp.stack.pop()
	stkOp.stack.push(y - x)
}

func divide(stkOp *stackOperator) {
	divisor := stkOp.stack.pop()
	dividend := stkOp.stack.pop()
	stkOp.stack.push(dividend / divisor)
}

func display(stkOp *stackOperator) {
	stkOp.stack.display()
}

func interactive() {
	scanner := bufio.NewScanner(os.Stdin)
	operationMap := make(map[string]*operation)
	operationMap["+"] = NewOperation(add, 2, 1)
	operationMap["-"] = NewOperation(subtract, 2, 1)
	operationMap["/"] = NewOperation(divide, 2, 1)
	operationMap["."] = NewOperation(display, 0, 0)
	stkOp := NewStackOperator(operationMap)
	for {
		fmt.Print("  > ")
		if !scanner.Scan() {
			fmt.Println()
			return
		}
		if err := scanner.Err(); err != nil {
            fmt.Println()
			log.Fatal(err)
		}
		split := strings.Split(scanner.Text(), " ")
		for i, s := range split {
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				stkOp.CallToken(s)
			} else {
				stkOp.stack.push(f)
			}
			if i == len(split)-1 {
				stkOp.stack.display()
			}
		}
	}
}

func main() {
	interactive()
}
