// Copyright 2024 Josh Tompkin
// Licensed under the MIT License

package stack

import "testing"

func getMaps() (*OrderedMap[string, int], map[string]int) {
	om := NewOrderedMap[string, int]()
	om.Set("arthropods", 0)
	om.Set("are", 1)
	om.Set("cool", 2)
	m := map[string]int{
		"are":        1,
		"arthropods": 0,
		"cool":       2,
	}
	return om, m
}

func TestOrderedMapIteration(t *testing.T) {
	om, m := getMaps()
	for k, v, ok := om.Next(); ok; k, v, ok = om.Next() {
		if mapVal := m[k]; mapVal != v {
			t.Fatalf("OrderedMap key = %q : expected value = %d : actual value = %d", k, mapVal, v)
		}
	}
}

func TestOrderedMapGet(t *testing.T) {
	om, m := getMaps()
	for s, i := range m {
		v, pres := om.Get(s)
		if !pres {
			t.Fatalf("OrderedMap key = %q : not present", s)
		}
		if v != i {
			t.Fatalf("OrderedMap key = %q : expected value = %d : actual value = %d", s, i, v)
		}
	}
}
