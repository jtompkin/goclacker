package actions

import "github.com/jtompkin/goclacker/stack"

func add(stkOp *stack.StackOperator) {
    stkOp.stack.push(stkOp.stack.pop() + stkOp.stack.pop())
    stkOp.operators
}
