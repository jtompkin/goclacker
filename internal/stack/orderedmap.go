// Copyright 2024 Josh Tompkin
// Licensed under the MIT License that
// can be found in the LICENSE file

package stack

// Pair contains a Key of type K and Value of type V.
type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

// OrderedMap contains a generic map of values, and an auxilary list to get
// items from in an ordered fashion.
type OrderedMap[K comparable, V any] struct {
	Pairs   map[K]*Pair[K, V]
	list    []*Pair[K, V]
	Current int
}

// NewOrderedMap returns a pointer to an OrderedMap with initialized values that
// is ready to use.
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{Pairs: make(map[K]*Pair[K, V], 16), list: make([]*Pair[K, V], 0, 16)}
}

// Set sets OrderedMap[key] = Pair{Key: key, Value: val} and panics if the key
// is already present in the OrderedMap.
func (om *OrderedMap[K, V]) Set(key K, val V) {
	_, present := om.Pairs[key]
	if present {
		panic("cannot redefine value in OrderedMap")
	}
	p := &Pair[K, V]{Key: key, Value: val}
	om.Pairs[key] = p
	om.list = append(om.list, p)
}

// Get returns the Pair at OrderedMap[key], and whether key is present in the
// OrderedMap.
func (om *OrderedMap[K, V]) Get(key K) (pair *Pair[K, V], present bool) {
	pair, present = om.Pairs[key]
	return pair, present
}

// Next returns the next Pair in OrderedMap sequentially, from the first one set
// to the last.
//
// for p := om.Next(); p != nil; p = om.next() {...}
func (om *OrderedMap[K, V]) Next() (pair *Pair[K, V]) {
	if om.Current > len(om.list)-1 {
		return nil
	}
	p := om.list[om.Current]
	om.Current++
	return om.Pairs[p.Key]
}

// Reset sets the value of OrderedMap.Current to 0.
func (om *OrderedMap[K, V]) Reset() {
	om.Current = 0
}
