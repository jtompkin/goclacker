package actions

import (
	"fmt"
	"math"

	"github.com/jtompkin/goclacker/stack"
)

func Add(stkOp *stack.StackOperator) error {
	stkOp.Stack.Push(stkOp.Stack.Pop() + stkOp.Stack.Pop())
	stkOp.Stack.Display()
	return nil
}

func Subtract(stkOp *stack.StackOperator) error {
	x := stkOp.Stack.Pop()
	y := stkOp.Stack.Pop()
	stkOp.Stack.Push(y - x)
	stkOp.Stack.Display()
	return nil
}

func Multiply(stkOp *stack.StackOperator) error {
	stkOp.Stack.Push(stkOp.Stack.Pop() * stkOp.Stack.Pop())
	stkOp.Stack.Display()
	return nil
}

func Divide(stkOp *stack.StackOperator) error {
	divisor := stkOp.Stack.Pop()
	if divisor == 0 {
		return stkOp.Fail("cannot divde by 0", divisor)
	}
	dividend := stkOp.Stack.Pop()
	stkOp.Stack.Push(dividend / divisor)
	stkOp.Stack.Display()
	return nil
}

func Power(stkOp *stack.StackOperator) error {
	exponent := stkOp.Stack.Pop()
	base := stkOp.Stack.Pop()
	if base == 0 && exponent < 0 {
		return stkOp.Fail("cannot raise 0 to negative power", base, exponent)
	}
	if base < 0 && exponent != float64(int(exponent)) {
		return stkOp.Fail("cannot raise negative number to non-integer power", base, exponent)
	}
	stkOp.Stack.Push(math.Pow(base, exponent))
    stkOp.Stack.Display()
	return nil
}

func Display(stkOp *stack.StackOperator) error {
	stkOp.Stack.Display()
	return nil
}

func Pop(stkOp *stack.StackOperator) error {
	stkOp.Stack.Pop()
	stkOp.Stack.Display()
	return nil
}

func Clear(stkOp *stack.StackOperator) error {
	fmt.Printf("cleared %d values\n", stkOp.Stack.Len())
	stkOp.Stack.SetValues(make([]float64, 0, stkOp.MaxStack))
	return nil
}
