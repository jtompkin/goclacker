// Copyright 2024 Josh Tompkin
// Licensed under the MIT License

package main

import (
	"fmt"
	"testing"
)

func prompt(t *testing.T, format string, expected string) {
	so := GetStackOperator(8, false, false, false)
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
		fmt.Sprintf("%c", FmtChar):     fmt.Sprintf("%c", FmtChar),
		fmt.Sprintf(" %c", FmtChar):    fmt.Sprintf(" %c", FmtChar),
		fmt.Sprintf("%c ", FmtChar):    fmt.Sprintf("%c ", FmtChar),
		fmt.Sprintf("%c-", FmtChar):    fmt.Sprintf("%c-", FmtChar),
		fmt.Sprintf("-%c", FmtChar):    fmt.Sprintf("-%c", FmtChar),
		fmt.Sprintf(" %c > ", FmtChar): fmt.Sprintf(" %c > ", FmtChar),
		fmt.Sprintf("%c%c%c", FmtChar, FmtChar, FmtChar):                     fmt.Sprintf("%c%c%c", FmtChar, FmtChar, FmtChar),
		fmt.Sprintf("%cl%cc&3t%cs%c10t", FmtChar, FmtChar, FmtChar, FmtChar): "80N N N12N N N N N N N N N N",
	}
	for format, expected := range formats {
		prompt(t, format, expected)
	}
}

func prog(t *testing.T, program string, params progParams) {
	so := GetStackOperator(8, false, false, false)
	err := so.ParseInput(program)
	s := string(so.PrintBuf)
	if err != nil {
		if params.WantError {
			return
		}
		s = err.Error()
	}
	if params.WantError {
		t.Fatalf(`program = "%s" : wanted error, none raised`, program)
	}
	if s != params.Expected && !params.AcceptAny {
		t.Fatalf(`program = "%s" : expected = %q : got = %q`, program, params.Expected, s)
	}
}

type progParams struct {
	Expected  string
	WantError bool
	AcceptAny bool
}

func TestPrograms(t *testing.T) {
	programs := map[string]progParams{
		"":             {"", false, false},
		"      ":       {"", false, false},
		"test":         {"", false, false},
		"1 2 3 4 5 6":  {"1 2 3 4 5 6\n", false, false},
		"2 2 +":        {"4\n", false, false},
		"6 2 -":        {"4\n", false, false},
		"2 2 *":        {"4\n", false, false},
		"8 2 /":        {"4\n", false, false},
		"15 4 %":       {"3\n", false, false},
		"2 3 ^":        {"8\n", false, false},
		"4 !":          {"24\n", false, false},
		"10 log":       {"1\n", false, false},
		"10 ln":        {"2.302585092994046\n", false, false},
		"4 sqrt":       {"2\n", false, false},
		"= pi":         {"deleted word: pi\n", false, false},
		"= test 2 2 +": {"defined test : 2 2 +\n", false, false},
		"pi sqrt":      {"1.7724538509055159\n", false, false},
		"+":            {"operation error: '+' needs 2 values in stack\n", false, false},
		"-1 log":       {"operation error: cannot take logarithm of non-positive number\n", false, false},
		"-1 ln":        {"operation error: cannot take logarithm of non-positive number\n", false, false},
		"=":            {"", true, false},
		"1 0 /":        {"", true, false},
		"help":         {"", false, true},
		"words":        {"", false, true},
		"  3 4 * 4455 -    23         + 4 4332     ": {"-4420 4 4332\n", false, false},
	}
	for program, params := range programs {
		prog(t, program, params)
	}
}

func TestConfig(t *testing.T) {
	testPath := "./test/test.conf"
	DefConfigPaths = append(DefConfigPaths, testPath)
	path := CheckDefConfigPaths()
	if path != testPath {
		t.Fatalf("Default config path %s not found.", testPath)
	}
	scanner := GetConfigScanner(path)
	if scanner == nil {
		t.Fatalf("Could not open config file %s", testPath)
	}
}
