package stack

// Pair contains a Key of type K and Value of type V.
type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

type list[K comparable, V any] struct {
	Values []*Pair[K, V]
	Len    int
}

func newList[K comparable, V any](capacity int) *list[K, V] {
	return &list[K, V]{Values: make([]*Pair[K, V], capacity)}
}

func (l *list[K, V]) increment(p *Pair[K, V]) {
	if l.Len+1 > cap(l.Values) {
		panic("list is full")
	}
	l.Values[l.Len] = p
	l.Len++
}

// OrderedMap contains a generic map of values, and an auxilary list to get
// items from in an ordered fashion.
type OrderedMap[K comparable, V any] struct {
	Pairs   map[K]*Pair[K, V]
	list    *list[K, V]
	Current int
}

// NewOrderedMap returns a pointer to an OrderedMap that has a given capacity,
// keys of type K, and values of type V.
func NewOrderedMap[K comparable, V any](capacity int) *OrderedMap[K, V] {
	om := &OrderedMap[K, V]{Pairs: make(map[K]*Pair[K, V], capacity), list: newList[K, V](capacity)}
	return om
}

// Set sets OrderedMap[key] = Pair{Key: key, Value: val} and panics if the key
// is already present in the OrderedMap.
func (om *OrderedMap[K, V]) Set(key K, val V) {
	_, present := om.Pairs[key]
	if present {
		panic("cannot redefine value in map")
	}
	p := &Pair[K, V]{Key: key, Value: val}
	om.Pairs[key] = p
	om.list.increment(p)
}

// Get returns the Pair at OrderedMap[key], and whether key is present in the
// OrderedMap.
func (om *OrderedMap[K, V]) Get(key K) (*Pair[K, V], bool) {
	p, present := om.Pairs[key]
	return p, present
}

// Next returns the next Pair in OrderedMap sequentially, from the first one set
// to the last.
//
// for p := om.Next(); p != nil; p = om.next() {...}
func (om *OrderedMap[K, V]) Next() *Pair[K, V] {
	if om.Current > om.list.Len-1 {
		return nil
	}
	p := om.list.Values[om.Current]
	om.Current++
	return om.Pairs[p.Key]
}

// Reset sets the value of OrderedMap.Current to 0.
func (om *OrderedMap[K, V]) Reset() {
    om.Current = 0
}
