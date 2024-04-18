package main

import (
	"testing"
)

func prompt(t *testing.T, format string, expected string) {
	so := MakeStackOperator(8, false)
	so.Stack.Stash = 12
	so.MakePromptFunc(format, '&')
	if s := so.Prompt(); s != expected {
		t.Fatalf(`format = "%s" : expected = "%s" : got  = "%s"`, format, expected, s)
	}
}

func TestPrompts(t *testing.T) {
	formats := map[string]string{
		"":         "",
		"&":        "&",
		" &":       " &",
		"& ":       "& ",
		"&a":       "&a",
		"a&":       "a&",
		"&&&":      "&&&",
		"     ":    "     ",
		"&l&c&s&t": "8012NA",
		" &c > ":   " 0 > ",
	}
	for format, expected := range formats {
		prompt(t, format, expected)
	}
}

func prog(t *testing.T, program string, expected string, wantError bool, acceptAny bool) {
	so := MakeStackOperator(8, false)
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
		"2 2 +":        newProgParams("4 \n", false, false),
		"4 sqrt":       newProgParams("2 \n", false, false),
		"= pi":         newProgParams("deleted word: pi\n", false, false),
		"= test 2 2 +": newProgParams(`defined word: "test" with value: "2 2 +"`+"\n", false, false),
		"=":            newProgParams("", true, false),
		"+":            newProgParams("", true, false),
		"1 0 /":        newProgParams("", true, false),
		"help":         newProgParams("", false, true),
		"words":        newProgParams("", false, true),
		"  3 4 * 4455 -    23         + 4 4332     ": newProgParams("-4420 4 4332 \n", false, false),
	}
	for program, expected := range programs {
		prog(t, program, expected.Expected, expected.WantError, expected.AcceptAny)
	}
}
