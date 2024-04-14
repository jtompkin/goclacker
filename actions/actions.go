package actions

import (
	"fmt"
	"math"
	"slices"

	"github.com/jtompkin/goclacker/stack"
)

// Action implements stack.Action.
type Action struct {
	action func(*stack.StackOperator) error
	pops   int
	pushes int
	help   string
}

func (a *Action) Call(stkOp *stack.StackOperator) error {
	return a.action(stkOp)
}

func (a *Action) GetPops() int {
	return a.pops
}

func (a *Action) GetPushes() int {
	return a.pushes
}

func (a *Action) GetHelp() string {
	return a.help
}

func newAction(
	action func(*stack.StackOperator) error,
	pops int,
	pushes int,
	help string,
) *Action {
	return &Action{action: action, pops: pops, pushes: pushes, help: help}
}

func Add() *Action {
	return newAction(
		func(so *stack.StackOperator) error {
			so.Stack.Push(so.Stack.Pop() + so.Stack.Pop())
			so.Stack.Display()
			return nil
		}, 2, 1,
		"pop 'a', 'b'; push the result of 'a' + 'b'",
	)
}

func Subtract() *Action {
	return newAction(
		func(so *stack.StackOperator) error {
			x := so.Stack.Pop()
			y := so.Stack.Pop()
			so.Stack.Push(y - x)
			so.Stack.Display()
			return nil
		}, 2, 1,
		"pop 'a', 'b'; push the result of 'b' - 'a'",
	)
}

func Multiply() *Action {
	return newAction(
		func(so *stack.StackOperator) error {
			so.Stack.Push(so.Stack.Pop() * so.Stack.Pop())
			so.Stack.Display()
			return nil
		}, 2, 1,
		"pop 'a', 'b'; push the result of 'a' * 'b'",
	)
}

func Divide() *Action {
	return newAction(
		func(so *stack.StackOperator) error {
			divisor := so.Stack.Pop()
			if divisor == 0 {
				return so.Fail("cannot divide by 0", divisor)
			}
			dividend := so.Stack.Pop()
			so.Stack.Push(dividend / divisor)
			so.Stack.Display()
			return nil
		}, 2, 1,
		"pop 'a', 'b'; push the result of 'b' / 'a'",
	)
}

func Power() *Action {
	return newAction(
		func(so *stack.StackOperator) error {
			exponent := so.Stack.Pop()
			base := so.Stack.Pop()
			if base == 0 && exponent < 0 {
				return so.Fail("cannot raise 0 to negative power", base, exponent)
			}
			if base < 0 && exponent != float64(int(exponent)) {
				return so.Fail("cannot raise negative number to non-integer power", base, exponent)
			}
			so.Stack.Push(math.Pow(base, exponent))
			so.Stack.Display()
			return nil
		}, 2, 1,
		"pop 'a', 'b'; push the result of 'b' ^ 'a'",
	)
}

func Log() *Action {
	return newAction(
		func(so *stack.StackOperator) error {
			x := so.Stack.Pop()
			if x <= 0 {
				return so.Fail("cannot take logarithm of non-positive number", x)
			}
			so.Stack.Push(math.Log10(x))
			so.Stack.Display()
			return nil
		}, 1, 1,
		"pop 'a'; push the logarithm base 10 of 'a'",
	)
}

func Ln() *Action {
	return newAction(
		func(so *stack.StackOperator) error {
			x := so.Stack.Pop()
			if x <= 0 {
				return so.Fail("cannot take logarithm of non-positive number", x)
			}
			so.Stack.Push(math.Log(x))
			so.Stack.Display()
			return nil
		}, 1, 1,
		"pop 'a'; push the natural logarithm of 'a'",
	)
}

func Round() *Action {
	return newAction(
		func(so *stack.StackOperator) error {
			precision := so.Stack.Pop()
			if precision < 0 || precision != float64(int(precision)) {
				return so.Fail("precision must be non-negative integer")
			}
			ratio := math.Pow(10, precision)
			so.Stack.Push(math.Round(so.Stack.Pop()*ratio) / ratio)
			so.Stack.Display()
			return nil
		}, 2, 1,
		"pop 'a', 'b'; push the result of rounding 'b' to 'a' decimal places",
	)
}

func Stash() *Action {
	return newAction(
		func(so *stack.StackOperator) error {
			so.Stack.SetStash(so.Stack.Pop())
			so.Stack.Display()
			return nil
		}, 1, 0,
		"pop 'a'; stash 'a'",
	)
}

func Pull() *Action {
	return newAction(
		func(so *stack.StackOperator) error {
			so.Stack.Push(so.Stack.GetStash())
			so.Stack.Display()
			return nil
		}, 0, 1,
		"push the value in the stash",
	)
}

func Display() *Action {
	return newAction(
		func(so *stack.StackOperator) error {
			so.Stack.Display()
			return nil
		}, 0, 0,
		"display all values in the stack",
	)
}

func Help() *Action {
	return newAction(
		func(so *stack.StackOperator) error {
            for _, token := range *so.Tokens {
                fmt.Printf("operator: %s\t\"%s\"\n", token, so.GetActions()[token].GetHelp())
            }
			return nil
		}, 0, 0,
		"display this information screen",
	)
}

func Words() *Action {
	return newAction(
		func(so *stack.StackOperator) error {
			keys := make([]string, 0, len(so.GetWords()))
			for k := range so.GetWords() {
				keys = append(keys, k)
			}
			slices.Sort(keys)
			for _, k := range keys {
				fmt.Printf("%s: %s\n", k, so.GetWords()[k])
			}
			return nil
		}, 0, 0,
		"display all defined words",
	)
}

func Pop() *Action {
	return newAction(
		func(so *stack.StackOperator) error {
			so.Stack.Pop()
			so.Stack.Display()
			return nil
		}, 1, 0,
		"pop 'a'",
	)
}

func Clear() *Action {
	return newAction(
		func(so *stack.StackOperator) error {
			var c rune
			n := so.Stack.Len()
			if n != 1 {
				c = 's'
			}
			fmt.Printf("cleared %d value%c\n", n, c)
			so.Stack.SetValues(make([]float64, 0, so.Stack.Cap()))
			return nil
		}, 0, 0,
		"pop all values in the stack",
	)
}
