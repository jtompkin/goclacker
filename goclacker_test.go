package main

import (
	"fmt"
	"testing"

	"github.com/jtompkin/goclacker/internal/stack"
)

func prompt(t *testing.T, format string, expected string) {
	so := MakeStackOperator(8, false, false)
	so.Stack.Stash = 12
	so.MakePromptFunc(format, '&')
	if s := so.Prompt(); s != expected {
		t.Fatalf(`format = "%s" : expected = "%s" : got  = "%s"`, format, expected, s)
	}
}

func TestPrompts(t *testing.T) {
	formats := map[string]string{
		"":                             "",
		"     ":                        "     ",
		fmt.Sprintf("%c", fmtChar):     fmt.Sprintf("%c", fmtChar),
		fmt.Sprintf(" %c", fmtChar):    fmt.Sprintf(" %c", fmtChar),
		fmt.Sprintf("%c ", fmtChar):    fmt.Sprintf("%c ", fmtChar),
		fmt.Sprintf("%c-", fmtChar):    fmt.Sprintf("%c-", fmtChar),
		fmt.Sprintf("-%c", fmtChar):    fmt.Sprintf("-%c", fmtChar),
		fmt.Sprintf(" %c > ", fmtChar): fmt.Sprintf(" %c > ", fmtChar),
		fmt.Sprintf("%c%c%c", fmtChar, fmtChar, fmtChar):                fmt.Sprintf("%c%c%c", fmtChar, fmtChar, fmtChar),
		fmt.Sprintf("%cl%cc%cs%ct", fmtChar, fmtChar, fmtChar, fmtChar): "8012NA",
	}
	for format, expected := range formats {
		prompt(t, format, expected)
	}
}

func prog(t *testing.T, program string, expected string, wantError bool, acceptAny bool) {
	so := MakeStackOperator(8, false, false)
	s, err := so.ParseInput(program)
	if err != nil {
		if wantError {
			return
		}
		s = err.Error()
	}
	if wantError {
		t.Fatalf(`program = "%s" : wanted error, none raised`, program)
	}
	if s != expected && !acceptAny {
		t.Fatalf(`program = "%s" : expected = %q : got = %q`, program, expected, s)
	}
}

type progParams struct {
	Expected  string
	WantError bool
	AcceptAny bool
}

func newProgParams(expected string, wantError bool, acceptAny bool) *progParams {
	return &progParams{Expected: expected, WantError: wantError, AcceptAny: acceptAny}
}

func TestPrograms(t *testing.T) {
	programs := map[string]*progParams{
		"":             newProgParams("", false, false),
		"      ":       newProgParams("", false, false),
		"test":         newProgParams("", false, false),
		"1 2 3 4 5 6":  newProgParams(fmt.Sprintf("1 2 3 4 5 6%s", stack.Suffix), false, false),
		"2 2 +":        newProgParams(fmt.Sprintf("4%s", stack.Suffix), false, false),
		"4 sqrt":       newProgParams(fmt.Sprintf("2%s", stack.Suffix), false, false),
		"= pi":         newProgParams(fmt.Sprintf("deleted word: pi%s", stack.Suffix), false, false),
		"= test 2 2 +": newProgParams(fmt.Sprintf(`defined word: "test" with value: "2 2 +"%s`, stack.Suffix), false, false),
		"pi sqrt":      newProgParams(fmt.Sprintf("1.7724538509055159%s", stack.Suffix), false, false),
		"+":            newProgParams(fmt.Sprintf("operation error: '+' needs 2 values in stack%s", stack.Suffix), false, false),
		"=":            newProgParams("", true, false),
		"1 0 /":        newProgParams("", true, false),
		"help":         newProgParams("", false, true),
		"words":        newProgParams("", false, true),
		"  3 4 * 4455 -    23         + 4 4332     ": newProgParams(fmt.Sprintf("-4420 4 4332%s", stack.Suffix), false, false),
	}
	for program, expected := range programs {
		prog(t, program, expected.Expected, expected.WantError, expected.AcceptAny)
	}
}
