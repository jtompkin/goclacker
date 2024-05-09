package stack

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strings"
)

type Action struct {
	// action is a function that does something to a StackOperator and returns a
	// string that is intended to be printed to the terminal screen and any
	// error that occurred during execution.
	action func(so *StackOperator) (toPrint string, err error)
	// Pops represents how many values a call to action will take from the
	// stack.
	Pops int
	// Pushes represents how many values a call to action will add to the stack.
	Pushes int
	// Help describes the purpose of the action.
	Help string
}

// Call calls the function stored in action and returns the string and error
// value returned by that function
func (a *Action) Call(so *StackOperator) (toPrint string, err error) {
	return a.action(so)
}

// Add is an Action with the following description: pop 'a', 'b'; push the
// result of 'a' + 'b'
var Add = &Action{
	func(so *StackOperator) (toPrint string, err error) {
		so.Stack.Push(so.Stack.Pop() + so.Stack.Pop())
		return so.Stack.Display(), nil
	},
	2, 1,
	"pop 'a', 'b'; push the result of 'a' + 'b'",
}

// Subtract is an Action with the following description:
var Subtract = &Action{
	func(so *StackOperator) (string, error) {
		x := so.Stack.Pop()
		y := so.Stack.Pop()
		so.Stack.Push(y - x)
		return so.Stack.Display(), nil
	}, 2, 1,
	"pop 'a', 'b'; push the result of 'b' - 'a'",
}

// Multiply is an Action with the following description: pop 'a', 'b'; push the
// result of 'a' * 'b'
var Multiply = &Action{
	func(so *StackOperator) (string, error) {
		so.Stack.Push(so.Stack.Pop() * so.Stack.Pop())
		return so.Stack.Display(), nil
	}, 2, 1,
	"pop 'a', 'b'; push the result of 'a' * 'b'",
}

// Divide is an Action with the following description: pop 'a', 'b'; push the
// result of 'b' / 'a'
var Divide = &Action{
	func(so *StackOperator) (string, error) {
		divisor := so.Stack.Pop()
		if divisor == 0 {
			return "", so.Fail("cannot divide by 0", divisor)
		}
		so.Stack.Push(so.Stack.Pop() / divisor)
		return so.Stack.Display(), nil
	}, 2, 1,
	"pop 'a', 'b'; push the result of 'b' / 'a'",
}

// Modulo is an Action with the following description: pop 'a', 'b'; push the
// remainder of 'b' / 'a'.
var Modulo = &Action{
	func(so *StackOperator) (string, error) {
		divisor := so.Stack.Pop()
		if divisor == 0 {
			return "", so.Fail("cannot divide by 0", divisor)
		}
		so.Stack.Push(math.Mod(so.Stack.Pop(), divisor))
		return so.Stack.Display(), nil
	}, 2, 1,
	"pop 'a', 'b'; push the remainder of 'b' / 'a'",
}


// Factorial is an Action with the following description: pop 'a'; push the
// factorial of 'a'.
var Factorial = &Action{
	func(so *StackOperator) (string, error) {
		x := so.Stack.Pop()
		if x != float64(int(x)) {
			return "", so.Fail("cannot take factorial of non-integer", x)
		}
		if x < 0 {
			return "", so.Fail("cannot take factorial of negative number", x)
		}
		p := 1
		for i := 2; i <= int(x); i++ {
			p *= i
		}
		so.Stack.Push(float64(p))
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the factorial of 'a'",
}

// Power is an Action with the following description: pop 'a', 'b'; push the
// result of 'b' ^ 'a'.
var Power = &Action{
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
}

// Log is an Action with the following description: pop 'a'; push the logarithm
// base 10 of 'a'.
var Log = &Action{
	func(so *StackOperator) (string, error) {
		x := so.Stack.Pop()
		if x <= 0 {
			return "", so.Fail("cannot take logarithm of non-positive number", x)
		}
		so.Stack.Push(math.Log10(x))
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the logarithm base 10 of 'a'",
}

// Ln is an Action with the following description: pop 'a'; push the natural
// logarithm of 'a'.
var Ln = &Action{
	func(so *StackOperator) (string, error) {
		x := so.Stack.Pop()
		if x <= 0 {
			return "", so.Fail("cannot take logarithm of non-positive number", x)
		}
		so.Stack.Push(math.Log(x))
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the natural logarithm of 'a'",
}

// Degrees is an Action with the following description: pop 'a'; push the result
// of converting 'a' from radians to degrees.
var Degrees = &Action{
	func(so *StackOperator) (string, error) {
		so.Stack.Push(so.Stack.Pop() * 180 / math.Pi)
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the result of converting 'a' from radians to degrees",
}

// Radians is an Action with the following description: pop 'a'; push the result
// of converting 'a' from degrees to radians.
var Radians = &Action{
	func(so *StackOperator) (string, error) {
		so.Stack.Push(so.Stack.Pop() * math.Pi / 180)
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the result of converting 'a' from degrees to radians",
}

// Sine is an Action with the following description: pop 'a'; push the sine of
// 'a' in radians.
var Sine = &Action{
	func(so *StackOperator) (string, error) {
		so.Stack.Push(math.Sin(so.Stack.Pop()))
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the sine of 'a' in radians",
}

// Cosine is an Action with the following description: pop 'a'; push the cosine
// of 'a' in radians.
var Cosine = &Action{
	func(so *StackOperator) (string, error) {
		so.Stack.Push(math.Cos(so.Stack.Pop()))
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the cosine of 'a' in radians",
}

// Tangent is an Action with the following description: pop 'a'; push the
// tangent of 'a' in radians.
var Tangent = &Action{
	func(so *StackOperator) (string, error) {
		so.Stack.Push(math.Tan(so.Stack.Pop()))
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the tangent of 'a' in radians",
}

// Floor is an Action with the following description: pop 'a'; push the greatest
// integer value less than or equal to 'a'.
var Floor = &Action{
	func(so *StackOperator) (string, error) {
		so.Stack.Push(math.Floor(so.Stack.Pop()))
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the greatest integer value less than or equal to 'a'",
}

// Ceiling is an Action with the following description: pop 'a'; push the least
// integer value greater than or equal to 'a'.
var Ceiling = &Action{
	func(so *StackOperator) (string, error) {
		so.Stack.Push(math.Ceil(so.Stack.Pop()))
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop 'a'; push the least integer value greater than or equal to 'a'",
}

// Round is an Action with the following description: pop 'a', 'b'; push the
// result of rounding 'b' to 'a' decimal places.
var Round = &Action{
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
}

// Random is an Action with the following description: push a random number
// between 0 and 1.
var Random = &Action{
	func(so *StackOperator) (string, error) {
		so.Stack.Push(rand.Float64())
		return so.Stack.Display(), nil
	}, 0, 1,
	"push a random number between 0 and 1",
}

// Stash is an Action with the following description: pop 'a'; stash 'a'.
var Stash = &Action{
	func(so *StackOperator) (string, error) {
		so.Stack.Stash = so.Stack.Pop()
		return so.Stack.Display(), nil
	}, 1, 0,
	"pop 'a'; stash 'a'",
}

// Pull is an Action with the following description: push the value in the
// stash.
var Pull = &Action{
	func(so *StackOperator) (string, error) {
		so.Stack.Push(so.Stack.Stash)
		return so.Stack.Display(), nil
	}, 0, 1,
	"push the value in the stash",
}

// Display is an Action with the following description: display all values in
// the stack.
var Display = &Action{
	func(so *StackOperator) (string, error) {
		return so.Stack.Display(), nil
	}, 0, 0,
	"display all values in the stack",
}

// Help is an Action with the following description: display this information
// screen.
var Help = &Action{
	func(so *StackOperator) (string, error) {
		sb := &strings.Builder{}
		for p := so.Actions.Next(); p != nil; p = so.Actions.Next() {
			sb.Write([]byte(fmt.Sprintf("%s%c%q\n", p.Key, '\t', p.Value.Help)))
		}
		so.Actions.Reset()
		return sb.String(), nil
	}, 0, 0,
	"display this information screen",
}

// Words is an Action with the following description: display all defined words.
var Words = &Action{
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
}

// Pop is an Action with the following description: pop 'a'.
var Pop = &Action{
	func(so *StackOperator) (string, error) {
		so.Stack.Pop()
		return so.Stack.Display(), nil
	}, 1, 0,
	"pop 'a'",
}

// Clear is an Action with the following description: pop all values in the
// stack.
var Clear = &Action{
	func(so *StackOperator) (string, error) {
		var c byte
		n := len(so.Stack.Values)
		if n != 1 {
			c = 's'
		}
		so.Stack.Values = make([]float64, 0, cap(so.Stack.Values))
		return fmt.Sprintf("cleared %d value%c%s", n, c, Suffix), nil
	}, 0, 0,
	"pop all values in the stack",
}

// ClearScreen is an Action with the following description: clear the terminal
// screen.
var ClearScreen = &Action{
	func(so *StackOperator) (string, error) {
		return "\x1b[2J\x1b[H", nil
	}, 0, 0,
	"clear the terminal screen",
}

// Swap is an Action with the following description: pop 'a', 'b'; push 'b',
// 'a'.
var Swap = &Action{
	func(so *StackOperator) (string, error) {
		x := so.Stack.Pop()
		y := so.Stack.Pop()
		so.Stack.Push(x)
		so.Stack.Push(y)
		return so.Stack.Display(), nil
	}, 2, 2,
	"pop 'a', 'b'; push 'b', 'a'",
}

// Froll is an Action with the following description: roll the stack to the
// right one position.
var Froll = &Action{
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
}

// Rroll is an Action with the following description: roll the stack to the left
// one position.
var Rroll = &Action{
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
}
