package actions

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strings"

	"github.com/jtompkin/goclacker/internal/stack"
)

// Action implements stack.Action.
type Action struct {
	action func(*stack.StackOperator) (string, error)
	pops   int
	pushes int
	help   string
}

// Call calls the function stored in Action.action and returns the error value
// returned by the function.
func (a *Action) Call(so *stack.StackOperator) (string, error) {
	return a.action(so)
}

func (a *Action) Pops() int {
	return a.pops
}

func (a *Action) Pushes() int {
	return a.pushes
}

func (a *Action) Help() string {
	return a.help
}

// newAction returns a pointer to Action initialized with values given to
// arguments.
func newAction(
	action func(*stack.StackOperator) (string, error),
	pops int,
	pushes int,
	help string,
) *Action {
	return &Action{action: action, pops: pops, pushes: pushes, help: help}
}

// Add pops returns a pointer to an Action that 'a', 'b'; pushes the result of
// 'a' + 'b'.
func Add() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			so.Stack.Push(so.Stack.Pop() + so.Stack.Pop())
			return so.Stack.Display(so.Interactive), nil
		}, 2, 1,
		"pop 'a', 'b'; push the result of 'a' + 'b'",
	)
}

// Subtract returns a pointer to an Action that pops 'a', 'b'; pushes the result
// of 'b' - 'a'.
func Subtract() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			x := so.Stack.Pop()
			y := so.Stack.Pop()
			so.Stack.Push(y - x)
			return so.Stack.Display(so.Interactive), nil
		}, 2, 1,
		"pop 'a', 'b'; push the result of 'b' - 'a'",
	)
}

// Multiply returns a pointer to an Action that pops 'a', 'b'; pushes the result
// of 'a' * 'b'.
func Multiply() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			so.Stack.Push(so.Stack.Pop() * so.Stack.Pop())
			return so.Stack.Display(so.Interactive), nil
		}, 2, 1,
		"pop 'a', 'b'; push the result of 'a' * 'b'",
	)
}

// Divide returns a pointer to an Action that pops 'a', 'b'; pushes the result
// of 'b' / 'a'.
func Divide() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			divisor := so.Stack.Pop()
			if divisor == 0 {
				return "", so.Fail("cannot divide by 0", divisor)
			}
			so.Stack.Push(so.Stack.Pop() / divisor)
			return so.Stack.Display(so.Interactive), nil
		}, 2, 1,
		"pop 'a', 'b'; push the result of 'b' / 'a'",
	)
}

// Modulo returns a pointer to an Action that pops 'a', 'b'; pushes the
// remainder of 'b' / 'a'
func Modulo() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			divisor := so.Stack.Pop()
			if divisor == 0 {
				return "", so.Fail("cannot divide by 0", divisor)
			}
			so.Stack.Push(math.Mod(so.Stack.Pop(), divisor))
			return so.Stack.Display(so.Interactive), nil
		}, 2, 1,
		"pop 'a', 'b'; push the remainder of 'b' / 'a'",
	)
}

func fact(x int) int {
	p := 1
	for i := 2; i <= x; i++ {
		p *= i
	}
	return p
}

// Factorial returns a pointer to an Action that pops 'a'; pushes the factorial
// of 'a'.
func Factorial() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			x := so.Stack.Pop()
			if x != float64(int(x)) {
				return "", so.Fail("cannot take factorial of non-integer", x)
			}
			if x < 0 {
				return "", so.Fail("cannot take factorial of negative number", x)
			}
			so.Stack.Push(float64(fact(int(x))))
			return so.Stack.Display(so.Interactive), nil
		}, 1, 1,
		"pop 'a'; push the factorial of 'a'",
	)
}

// Power returns a pointer to an Action that pops 'a', 'b'; pushes the result of
// 'b' ^ 'a'.
func Power() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			exponent := so.Stack.Pop()
			base := so.Stack.Pop()
			if base == 0 && exponent < 0 {
				return "", so.Fail("cannot raise 0 to negative power", base, exponent)
			}
			if base < 0 && exponent != float64(int(exponent)) {
				return "", so.Fail("cannot raise negative number to non-integer power", base, exponent)
			}
			so.Stack.Push(math.Pow(base, exponent))
			return so.Stack.Display(so.Interactive), nil
		}, 2, 1,
		"pop 'a', 'b'; push the result of 'b' ^ 'a'",
	)
}

// Log returns a pointer to an Action that pops 'a'; pushes the logarithm base
// 10 of 'a'.
func Log() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			x := so.Stack.Pop()
			if x <= 0 {
				return "", so.Fail("cannot take logarithm of non-positive number", x)
			}
			so.Stack.Push(math.Log10(x))
			return so.Stack.Display(so.Interactive), nil
		}, 1, 1,
		"pop 'a'; push the logarithm base 10 of 'a'",
	)
}

// Ln returns a pointer to an Action that pops 'a'; pushes the logarithm base 10
// of 'a'.
func Ln() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			x := so.Stack.Pop()
			if x <= 0 {
				return "", so.Fail("cannot take logarithm of non-positive number", x)
			}
			so.Stack.Push(math.Log(x))
			return so.Stack.Display(so.Interactive), nil
		}, 1, 1,
		"pop 'a'; push the natural logarithm of 'a'",
	)
}

// Degrees returns a pointer to an Action that pops 'a'; pushes the result of
// converting 'a' from radians to degrees.
func Degrees() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			so.Stack.Push(so.Stack.Pop() * 180 / math.Pi)
			return so.Stack.Display(so.Interactive), nil
		}, 1, 1,
		"pop 'a'; push the result of converting 'a' from radians to degrees",
	)
}

// Radians returns a pointer to an Action that pops 'a'; pushes the result of
// converting 'a' from degrees to radians.
func Radians() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			so.Stack.Push(so.Stack.Pop() * math.Pi / 180)
			return so.Stack.Display(so.Interactive), nil
		}, 1, 1,
		"pop 'a'; push the result of converting 'a' from degrees to radians",
	)
}

// Sine returns a pointer to an Action that pops 'a'; pushes the sine of 'a' in
// radians.
func Sine() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			so.Stack.Push(math.Sin(so.Stack.Pop()))
			return so.Stack.Display(so.Interactive), nil
		}, 1, 1,
		"pop 'a'; push the sine of 'a' in radians",
	)
}

// Cosine returns a pointer to an Action that pops 'a'; pushes the cosine of 'a'
// in radians.
func Cosine() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			so.Stack.Push(math.Cos(so.Stack.Pop()))
			return so.Stack.Display(so.Interactive), nil
		}, 1, 1,
		"pop 'a'; push the cosine of 'a' in radians",
	)
}

// Tangent returns a pointer to an Action that pops 'a'; pushes the tangent of
// 'a' in radians.
func Tangent() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			so.Stack.Push(math.Tan(so.Stack.Pop()))
			return so.Stack.Display(so.Interactive), nil
		}, 1, 1,
		"pop 'a'; push the tangent of 'a' in radians",
	)
}

// Floor returns a pointer to an Action that pops 'a'; pushes the greatest
// integer value less than 'a'.
func Floor() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			so.Stack.Push(math.Floor(so.Stack.Pop()))
			return so.Stack.Display(so.Interactive), nil
		}, 1, 1,
		"pop 'a'; push the greatest integer value less than or equal to 'a'",
	)
}

// Ceiling returns a pointer to an Action that pops 'a'; pushes the least
// integer value greater than 'a'.
func Ceiling() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			so.Stack.Push(math.Ceil(so.Stack.Pop()))
			return so.Stack.Display(so.Interactive), nil
		}, 1, 1,
		"pop 'a'; push the least integer value greater than or equal to 'a'",
	)
}

// Round returns a pointer to an Action that pops 'a', 'b'; pushes the result of
// rounding 'b' to 'a' decimal places.
func Round() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			precision := so.Stack.Pop()
			if precision < 0 || precision != float64(int(precision)) {
				return "", so.Fail("precision must be non-negative integer")
			}
			ratio := math.Pow(10, precision)
			so.Stack.Push(math.Round(so.Stack.Pop()*ratio) / ratio)
			return so.Stack.Display(so.Interactive), nil
		}, 2, 1,
		"pop 'a', 'b'; push the result of rounding 'b' to 'a' decimal places",
	)
}

// Random returns a pointer to an Action that pushes a random number between 0
// and 1
func Random() *Action {
    return newAction(
        func(so *stack.StackOperator) (string, error) {
            so.Stack.Push(rand.Float64())
            return so.Stack.Display(so.Interactive), nil
        }, 0, 1,
        "push a random number between 0 and 1",
        )
}

// Stash returns a pointer to an Action that pops 'a'; stashes 'a'.
func Stash() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			so.Stack.Stash = so.Stack.Pop()
			return so.Stack.Display(so.Interactive), nil
		}, 1, 0,
		"pop 'a'; stash 'a'",
	)
}

// Pull returns a pointer to an Action that pushes the value in the stash.
func Pull() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			so.Stack.Push(so.Stack.Stash)
			return so.Stack.Display(so.Interactive), nil
		}, 0, 1,
		"push the value in the stash",
	)
}

// Display returns a pointer to an Action that displays all values in the stack.
func Display() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			return so.Stack.Display(so.Interactive), nil
		}, 0, 0,
		"display all values in the stack",
	)
}

// Help returns a pointer to an Action that displays an information screen.
func Help() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			helps := make([]string, 0, so.Actions.Len())
			for pair := so.Actions.Oldest(); pair != nil; pair = pair.Next() {
				helps = append(helps, fmt.Sprintf(`%s%c"%s"`, pair.Key, '\t', pair.Value.Help()))
				//helps = append(helps, fmt.Sprintf(`operator: %s%c"%s"`, pair.Key, '\t', pair.Value.Help()))
			}
			return strings.Join(helps, "\n") + stack.Suffix, nil
		}, 0, 0,
		"display this information screen",
	)
}

// Words returns a pointer to an Action that displays all defined words.
func Words() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			keys := make([]string, 0, len(so.Words))
			for k := range so.Words {
				keys = append(keys, k)
			}
			slices.Sort(keys)
			defs := make([]string, len(keys))
			for i, k := range keys {
				defs[i] = fmt.Sprintf("%s : %s", k, so.Words[k])
			}
			return strings.Join(defs, "\n") + stack.Suffix, nil
		}, 0, 0,
		"display all defined words",
	)
}

// Pop returns a pointer to an Action that pops 'a'.
func Pop() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			so.Stack.Pop()
			return so.Stack.Display(so.Interactive), nil
		}, 1, 0,
		"pop 'a'",
	)
}

// Clear returns a pointer to an Action that pops all values in the stack.
func Clear() *Action {
	return newAction(
		func(so *stack.StackOperator) (string, error) {
			var c rune
			n := len(so.Stack.Values)
			if n != 1 {
				c = 's'
			}
			so.Stack.Values = make([]float64, 0, cap(so.Stack.Values))
			return fmt.Sprintf("cleared %d value%c%s", n, c, stack.Suffix), nil
		}, 0, 0,
		"pop all values in the stack",
	)
}
