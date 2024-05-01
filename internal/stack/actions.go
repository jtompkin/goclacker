package stack

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strings"
)

type Action struct {
	action func(*StackOperator) (string, error)
	Pops   int
	Pushes int
	Help   string
}

// Call calls the function stored in Action.action and returns the error value
// returned by the function.
func (a *Action) Call(so *StackOperator) (string, error) {
	return a.action(so)
}

// newAction returns a pointer to Action initialized with values given to
// arguments.
func newAction(
	action func(*StackOperator) (string, error),
	pops int,
	pushes int,
	help string,
) *Action {
	return &Action{action: action, Pops: pops, Pushes: pushes, Help: help}
}

var Add = newAction(
	func(so *StackOperator) (string, error) {
		so.Stack.Push(so.Stack.Pop() + so.Stack.Pop())
		return so.Stack.Display(), nil
	}, 2, 1,
	"pop 'a', 'b'; push the result of 'a' + 'b'",
)

var Subtract = newAction(
	func(so *StackOperator) (string, error) {
		x := so.Stack.Pop()
		y := so.Stack.Pop()
		so.Stack.Push(y - x)
		return so.Stack.Display(), nil
	}, 2, 1,
	"pop 'a', 'b'; push the result of 'b' - 'a'",
)

var Multiply = newAction(
	func(so *StackOperator) (string, error) {
		so.Stack.Push(so.Stack.Pop() * so.Stack.Pop())
		return so.Stack.Display(), nil
	}, 2, 1,
	"pop 'a', 'b'; push the result of 'a' * 'b'",
)

var Divide = newAction(
	func(so *StackOperator) (string, error) {
		divisor := so.Stack.Pop()
		if divisor == 0 {
			return "", so.Fail("cannot divide by 0", divisor)
		}
		so.Stack.Push(so.Stack.Pop() / divisor)
		return so.Stack.Display(), nil
	}, 2, 1,
	"pop 'a', 'b'; push the result of 'b' / 'a'",
)

var Modulo = newAction(
	func(so *StackOperator) (string, error) {
		divisor := so.Stack.Pop()
		if divisor == 0 {
			return "", so.Fail("cannot divide by 0", divisor)
		}
		so.Stack.Push(math.Mod(so.Stack.Pop(), divisor))
		return so.Stack.Display(), nil
	}, 2, 1,
	"pop 'a', 'b'; push the remainder of 'b' / 'a'",
)

func fact(x int) int {
	p := 1
	for i := 2; i <= x; i++ {
		p *= i
	}
	return p
}

var Factorial = newAction(
	func(so *StackOperator) (string, error) {
		x := so.Stack.Pop()
		if x != float64(int(x)) {
			return "", so.Fail("cannot take factorial of non-integer", x)
		}
		if x < 0 {
			return "", so.Fail("cannot take factorial of negative number", x)
		}
		so.Stack.Push(float64(fact(int(x))))
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the factorial of 'a'",
)

var Power = newAction(
	func(so *StackOperator) (string, error) {
		exponent := so.Stack.Pop()
		base := so.Stack.Pop()
		if base == 0 && exponent < 0 {
			return "", so.Fail("cannot raise 0 to negative power", base, exponent)
		}
		if base < 0 && exponent != float64(int(exponent)) {
			return "", so.Fail("cannot raise negative number to non-integer power", base, exponent)
		}
		so.Stack.Push(math.Pow(base, exponent))
		return so.Stack.Display(), nil
	}, 2, 1,
	"pop 'a', 'b'; push the result of 'b' ^ 'a'",
)

var Log = newAction(
	func(so *StackOperator) (string, error) {
		x := so.Stack.Pop()
		if x <= 0 {
			return "", so.Fail("cannot take logarithm of non-positive number", x)
		}
		so.Stack.Push(math.Log10(x))
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the logarithm base 10 of 'a'",
)

var Ln = newAction(
	func(so *StackOperator) (string, error) {
		x := so.Stack.Pop()
		if x <= 0 {
			return "", so.Fail("cannot take logarithm of non-positive number", x)
		}
		so.Stack.Push(math.Log(x))
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the natural logarithm of 'a'",
)

var Degrees = newAction(
	func(so *StackOperator) (string, error) {
		so.Stack.Push(so.Stack.Pop() * 180 / math.Pi)
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the result of converting 'a' from radians to degrees",
)

var Radians = newAction(
	func(so *StackOperator) (string, error) {
		so.Stack.Push(so.Stack.Pop() * math.Pi / 180)
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the result of converting 'a' from degrees to radians",
)

var Sine = newAction(
	func(so *StackOperator) (string, error) {
		so.Stack.Push(math.Sin(so.Stack.Pop()))
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the sine of 'a' in radians",
)

var Cosine = newAction(
	func(so *StackOperator) (string, error) {
		so.Stack.Push(math.Cos(so.Stack.Pop()))
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the cosine of 'a' in radians",
)

var Tangent = newAction(
	func(so *StackOperator) (string, error) {
		so.Stack.Push(math.Tan(so.Stack.Pop()))
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the tangent of 'a' in radians",
)

var Floor = newAction(
	func(so *StackOperator) (string, error) {
		so.Stack.Push(math.Floor(so.Stack.Pop()))
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the greatest integer value less than or equal to 'a'",
)

var Ceiling = newAction(
	func(so *StackOperator) (string, error) {
		so.Stack.Push(math.Ceil(so.Stack.Pop()))
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the least integer value greater than or equal to 'a'",
)

var Round = newAction(
	func(so *StackOperator) (string, error) {
		precision := so.Stack.Pop()
		if precision < 0 || precision != float64(int(precision)) {
			return "", so.Fail("precision must be non-negative integer")
		}
		ratio := math.Pow(10, precision)
		so.Stack.Push(math.Round(so.Stack.Pop()*ratio) / ratio)
		return so.Stack.Display(), nil
	}, 2, 1,
	"pop 'a', 'b'; push the result of rounding 'b' to 'a' decimal places",
)

var Random = newAction(
	func(so *StackOperator) (string, error) {
		so.Stack.Push(rand.Float64())
		return so.Stack.Display(), nil
	}, 0, 1,
	"push a random number between 0 and 1",
)

var Stash = newAction(
	func(so *StackOperator) (string, error) {
		so.Stack.Stash = so.Stack.Pop()
		return so.Stack.Display(), nil
	}, 1, 0,
	"pop 'a'; stash 'a'",
)

var Pull = newAction(
	func(so *StackOperator) (string, error) {
		so.Stack.Push(so.Stack.Stash)
		return so.Stack.Display(), nil
	}, 0, 1,
	"push the value in the stash",
)

var Display = newAction(
	func(so *StackOperator) (string, error) {
		return so.Stack.Display(), nil
	}, 0, 0,
	"display all values in the stack",
)

var Help = newAction(
	func(so *StackOperator) (string, error) {
		sb := &strings.Builder{}
		for p := so.Actions.Next(); p != nil; p = so.Actions.Next() {
			sb.Write([]byte(fmt.Sprintf("%s%c%q\n", p.Key, '\t', p.Value.Help)))
		}
		so.Actions.Reset()
		return sb.String(), nil
	}, 0, 0,
	"display this information screen",
)

var Words = newAction(
	func(so *StackOperator) (string, error) {
		keys := make([]string, 0, len(so.Words))
		for k := range so.Words {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		defs := make([]string, len(keys))
		for i, k := range keys {
			defs[i] = fmt.Sprintf("%s : %s", k, so.Words[k])
		}
		return strings.Join(defs, "\n") + Suffix, nil
	}, 0, 0,
	"display all defined words",
)

var Pop = newAction(
	func(so *StackOperator) (string, error) {
		so.Stack.Pop()
		return so.Stack.Display(), nil
	}, 1, 0,
	"pop 'a'",
)

var Clear = newAction(
	func(so *StackOperator) (string, error) {
		var c rune
		n := len(so.Stack.Values)
		if n != 1 {
			c = 's'
		}
		so.Stack.Values = make([]float64, 0, cap(so.Stack.Values))
		return fmt.Sprintf("cleared %d value%c%s", n, c, Suffix), nil
	}, 0, 0,
	"pop all values in the stack",
)

var ClearScreen = newAction(
	func(so *StackOperator) (string, error) {
		return "\033[2J\033[H", nil
	}, 0, 0,
	"clear the terminal screen",
)

var Swap = newAction(
	func(so *StackOperator) (string, error) {
		x := so.Stack.Pop()
		y := so.Stack.Pop()
		so.Stack.Push(x)
		so.Stack.Push(y)
		return so.Stack.Display(), nil
	}, 2, 2,
	"pop 'a', 'b'; push 'b', 'a'",
)

var Froll = newAction(
	func(so *StackOperator) (string, error) {
		newVals := make([]float64, len(so.Stack.Values))
		newVals[0] = so.Stack.Values[len(so.Stack.Values)-1]
		for i, f := range so.Stack.Values[:len(so.Stack.Values)-1] {
			newVals[i+1] = f
		}
		so.Stack.Values = newVals
		return so.Stack.Display(), nil
	}, 2, 2,
	"roll the stack to the right one position",
)

var Rroll = newAction(
	func(so *StackOperator) (string, error) {
		newVals := make([]float64, len(so.Stack.Values))
		for i, f := range so.Stack.Values[1:] {
			newVals[i] = f
		}
		newVals[len(newVals)-1] = so.Stack.Values[0]
		so.Stack.Values = newVals
		return so.Stack.Display(), nil
	}, 2, 2,
	"roll the stack to the left one position",
)
