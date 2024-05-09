// Copyright 2024 Josh Tompkin
// Licensed under the MIT License that
// can be found in the LICENSE file

package main

import (
	"fmt"
	"testing"

	"github.com/jtompkin/goclacker/internal/stack"
)

func prompt(t *testing.T, format string, expected string) {
	so := MakeStackOperator(8, false, false, false)
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
	so := MakeStackOperator(8, false, false, false)
	err := so.ParseInput(program)
	s := string(so.PrintBuf)
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

func TestPrograms(t *testing.T) {
	s := stack.Suffix
	programs := map[string]*progParams{
		"":             {"", false, false},
		"      ":       {"", false, false},
		"test":         {"", false, false},
		"1 2 3 4 5 6":  {"1 2 3 4 5 6" + s, false, false},
		"2 2 +":        {"4" + s, false, false},
		"6 2 -":        {"4" + s, false, false},
		"2 2 *":        {"4" + s, false, false},
		"8 2 /":        {"4" + s, false, false},
		"15 4 %":       {"3" + s, false, false},
		"2 3 ^":        {"8" + s, false, false},
		"4 !":          {"24" + s, false, false},
		"10 log":       {"1" + s, false, false},
		"10 ln":        {"2.302585092994046" + s, false, false},
		"4 sqrt":       {"2" + s, false, false},
		"= pi":         {"deleted pi" + s, false, false},
		"= test 2 2 +": {"defined test : 2 2 +" + s, false, false},
		"pi sqrt":      {"1.7724538509055159" + s, false, false},
		"+":            {"operation error: '+' needs 2 values in stack" + s, false, false},
		"-1 log":       {"operation error: cannot take logarithm of non-positive number" + s, false, false},
		"-1 ln":        {"operation error: cannot take logarithm of non-positive number" + s, false, false},
		"=":            {"", true, false},
		"1 0 /":        {"", true, false},
		"help":         {"", false, true},
		"words":        {"", false, true},
		"  3 4 * 4455 -    23         + 4 4332     ": {fmt.Sprintf("-4420 4 4332%s", stack.Suffix), false, false},
	}
	for program, expected := range programs {
		prog(t, program, expected.Expected, expected.WantError, expected.AcceptAny)
	}
}
