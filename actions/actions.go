package actions

import (
	"fmt"
	"math"
	"slices"

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

func Log(stkOp *stack.StackOperator) error {
	x := stkOp.Stack.Pop()
	if x <= 0 {
		return stkOp.Fail("cannot take logarithm of non-positive number", x)
	}
	stkOp.Stack.Push(math.Log10(x))
	stkOp.Stack.Display()
	return nil
}

func Ln(stkOp *stack.StackOperator) error {
	x := stkOp.Stack.Pop()
	if x <= 0 {
		return stkOp.Fail("cannot take logarithm of non-positive number", x)
	}
	stkOp.Stack.Push(math.Log(x))
	stkOp.Stack.Display()
	return nil
}

func Round(stkOp *stack.StackOperator) error {
	precision := stkOp.Stack.Pop()
	if precision < 0 || precision != float64(int(precision)) {
		return stkOp.Fail("precision must be non-negative integer")
	}
	ratio := math.Pow(10, precision)
	stkOp.Stack.Push(math.Round(stkOp.Stack.Pop()*ratio) / ratio)
	stkOp.Stack.Display()
	return nil
}

func Stash(stkOp *stack.StackOperator) error {
	stkOp.Stack.SetStash(stkOp.Stack.Pop())
	stkOp.Stack.Display()
	return nil
}

func Pull(stkOp *stack.StackOperator) error {
	stkOp.Stack.Push(stkOp.Stack.GetStash())
	stkOp.Stack.Display()
	return nil
}

func Display(stkOp *stack.StackOperator) error {
	stkOp.Stack.Display()
	return nil
}

func Help(stkOp *stack.StackOperator) error {
	if err := stkOp.PrintHelp(); err != nil {
		return err
	}
	return nil
}

func Words(stkOp *stack.StackOperator) error {
	words := make([]string, 0, len(stkOp.GetWords()))
	for k := range stkOp.GetWords() {
		words = append(words, k)
	}
	slices.Sort(words)
	for _, w := range words {
		fmt.Printf("%s: %s\n", w, stkOp.GetWords()[w])
	}
	return nil
}

func Pop(stkOp *stack.StackOperator) error {
	stkOp.Stack.Pop()
	stkOp.Stack.Display()
	return nil
}

func Clear(stkOp *stack.StackOperator) error {
	var c string
	n := stkOp.Stack.Len()
	if n != 1 {
		c = "s"
	}
	fmt.Printf("cleared %d value%s\n", n, c)
	stkOp.Stack.SetValues(make([]float64, 0, stkOp.Stack.Cap()))
	return nil
}
