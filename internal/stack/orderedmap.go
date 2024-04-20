package stack


type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

type list[K comparable, V any] struct {
	Values  []*Pair[K, V]
	current int
}

func newList[K comparable, V any](capacity int) *list[K, V] {
	return &list[K, V]{Values: make([]*Pair[K, V], capacity)}
}

func (l *list[K, V]) increment(p *Pair[K, V]) {
	if l.current+1 > cap(l.Values) {
		panic("list is full")
	}
	l.Values[l.current] = p
	l.current++
}

type OrderedMap[K comparable, V any] struct {
	Pairs   map[K]*Pair[K, V]
	list    *list[K, V]
	Current int
}

func NewOrderedMap[K comparable, V any](capacity int) *OrderedMap[K, V] {
	om := &OrderedMap[K, V]{Pairs: make(map[K]*Pair[K, V], capacity), list: newList[K, V](capacity)}
	return om
}

func (om *OrderedMap[K, V]) Set(key K, val V) {
	_, present := om.Pairs[key]
	if present {
		panic("cannot redefine value in map")
	}
	p := &Pair[K, V]{Key: key, Value: val}
	om.Pairs[key] = p
	om.list.increment(p)
}

func (om *OrderedMap[K, V]) Get(key K) (*Pair[K, V], bool) {
    p, present := om.Pairs[key]
    return p, present
}

func (om *OrderedMap[K, V]) Next() *Pair[K, V] {
    if om.Current > om.list.current-1 {
        return nil
    }
	p := om.list.Values[om.Current]
	om.Current++
	return om.Pairs[p.Key]
}
