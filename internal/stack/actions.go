// Copyright 2024 Josh Tompkin
// Licensed under the MIT License

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
		sBuf := make([]string, len(so.Stack.Values))
		for i, f := range so.Stack.Values {
			sBuf[i] = fmt.Sprint(f)
		}
		return fmt.Sprintf("[ %s ]\n", strings.Join(sBuf, " ")), nil
	}, 0, 0,
	"display all values in the stack",
}

// Help is an Action with the following description: display this information
// screen.
var Help = &Action{
	func(so *StackOperator) (string, error) {
		header := "operator"
		maxLen := len(header)
		for p := so.Actions.Next(); p != nil; p = so.Actions.Next() {
			if len(p.Key) > maxLen {
				maxLen = len(p.Key)
			}
		}
		so.Actions.Reset()
		sb := new(strings.Builder)
		pad := strings.Repeat(" ", maxLen-len(header))
		sb.WriteString(fmt.Sprintf("%s%s | %s\n", pad, header, "description"))
		for p := so.Actions.Next(); p != nil; p = so.Actions.Next() {
			pad := strings.Repeat(" ", maxLen-len(p.Key))
			sb.WriteString(fmt.Sprintf("%s%s : %s\n", pad, p.Key, p.Value.Help))
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
		header := "word"
		maxLen := len(header)
		for k := range so.Words {
			keys = append(keys, k)
			if len(k) > maxLen {
				maxLen = len(k)
			}
		}
		slices.SortFunc(keys, func(a string, b string) int {
			if len(a) > len(b) {
				return -1
			}
			if len(a) < len(b) {
				return 1
			}
			if a > b {
				return 1
			}
			if a < b {
				return -1
			}
			return 0
		})
		sb := new(strings.Builder)
		pad := strings.Repeat(" ", maxLen-len(header))
		sb.WriteString(fmt.Sprintf("%s%s | definition\n", pad, header))
		for _, k := range keys {
			pad := strings.Repeat(" ", maxLen-len(k))
			sb.WriteString(fmt.Sprintf("%s%s : %s\n", pad, k, so.Words[k]))
		}
		return sb.String(), nil
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
		return fmt.Sprintf("cleared %d value%c\n", n, c), nil
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
		newVals := make([]float64, 0, cap(so.Stack.Values))
		l := len(so.Stack.Values)
		newVals = append(newVals, so.Stack.Values[l-1])
		for _, f := range so.Stack.Values[:l-1] {
			newVals = append(newVals, f)
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
		newVals := make([]float64, 0, cap(so.Stack.Values))
		for _, f := range so.Stack.Values[1:] {
			newVals = append(newVals, f)
		}
		newVals = append(newVals, so.Stack.Values[0])
		so.Stack.Values = newVals
		return so.Stack.Display(), nil
	}, 2, 2,
	"roll the stack to the left one position",
}

// Sum is an Action with the following description: pop all values in the stack;
// push their sum.
var Sum = &Action{
	func(so *StackOperator) (toPrint string, err error) {
		var sum float64
		for len(so.Stack.Values) > 0 {
			sum += so.Stack.Pop()
		}
		so.Stack.Push(sum)
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop all values in the stack; push their sum",
}

// Average is an Action with the following description: pop all values in the
// stack; push their average.
var Average = &Action{
	func(so *StackOperator) (toPrint string, err error) {
		n := float64(len(so.Stack.Values))
		Sum.Call(so)
		so.Stack.Push(so.Stack.Pop() / n)
		return so.Stack.Display(), nil
	}, 1, 1,
	"pop all values in the stack; push their average",
}

var Clip = &Action{
	func(so *StackOperator) (toPrint string, err error) {
		c := cap(so.Stack.Values)
		so.Stack.Values = slices.Clip(so.Stack.Values)
		return fmt.Sprintf("clipped %d capacity\n", c-cap(so.Stack.Values)), nil
	}, 0, 0,
	"DEBUG: clip unused stack capacity",
}

// Grow is an Action with the following description: DEBUG: pop 'a'; push 'a';
// grow stack to accomadate 'a' more values. Will commit blasphemy and grow
// stack by 1 if cap(Stack.Values) == 0.
var Grow = &Action{
	func(so *StackOperator) (toPrint string, err error) {
		if cap(so.Stack.Values) == 0 {
			so.Stack.Values = slices.Grow(so.Stack.Values, 1)
			return fmt.Sprintf("new stack capacity is %d\n", cap(so.Stack.Values)), nil
		}
		if len(so.Stack.Values) == 0 {
			return "", nil
		}
		n := so.Stack.Pop()
		if n != float64(int(n)) {
			return "", so.Fail("cannot grow stack by non-integer value", n)
		}
		if n < 0 {
			return "", so.Fail("cannot grow stack by negative value", n)
		}
		so.Stack.Push(n)
		so.Stack.Values = slices.Grow(so.Stack.Values, int(n))
		return fmt.Sprintf("new stack capacity is %d\n", cap(so.Stack.Values)), nil
	}, 0, 0,
	"DEBUG: pop 'a'; push 'a'; grow stack to accomadate 'a' more values",
}

// Fill is an Action with the following description: DEBUG: fill stack with
// random values
var Fill = &Action{
	func(so *StackOperator) (toPrint string, err error) {
		for i := 0; i < cap(so.Stack.Values); i++ {
			if i > len(so.Stack.Values)-1 {
				so.Stack.Values = append(so.Stack.Values, float64(rand.Intn(255)))
			}
		}
		return so.Stack.Display(), nil
	}, 0, 0,
	"DEBUG: fill stack with random values",
}
