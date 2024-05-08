package stack

import (
	"testing"
)

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
	for p := om.Next(); p != nil; p = om.Next() {
		if val := m[p.Key]; val != p.Value {
			t.Fatalf("om key = %q : expected value = %d : actual value = %d", p.Key, val, p.Value)
		}
	}
}

func TestOrderedMapGet(t *testing.T) {
	om, m := getMaps()
	for s, i := range m {
		p, pres := om.Get(s)
		if p.Value != i {
			t.Fatalf("om key = %q : expected value = %d : actual value = %d", p.Key, i, p.Value)
		}
		if !pres {
			t.Fatalf("om key = %q : not present", p.Key)
		}
	}
}
