package main

import (
	"testing"

	"github.com/jtompkin/goclacker/internal/stack"
)

func basicStackOperator() *stack.StackOperator {
    return stack.NewStackOperator(make(map[string]stack.Action), make([]string, 0), 8)
}

func prompt(t *testing.T, format string, expected string) {
    so := basicStackOperator()
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
/*
func TestProgramBasic(t *testing.T) {
    so := basicStackOperator()
    so.ParseInput("2 2 +")
    b := make([]byte, 3)
    os.Stdout.Read(b)
    if string(b) != "[4]" {
        t.Fatalf(`program = "2 2 +" : expected = "[4]" : got = %s`, b)
    }
}
*/
