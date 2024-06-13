package stack

type OrderedMap[K comparable, V any] struct {
	List    []K
	Pairs   map[K]V
	current int
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{make([]K, 0), make(map[K]V), 0}
}

func (om *OrderedMap[K, V]) Set(key K, val V) {
	if _, pres := om.Pairs[key]; pres {
		om.Delete(key)
	}
	om.List = append(om.List, key)
	om.Pairs[key] = val
}

func (om *OrderedMap[K, V]) Get(key K) (val V, pres bool) {
	val, pres = om.Pairs[key]
	return
}

func (om *OrderedMap[K, V]) Next() (key K, val V, ok bool) {
	if om.current == len(om.List) {
		om.current = 0
		return key, val, false
	}
	key = om.List[om.current]
	val = om.Pairs[key]
	om.current++
	return key, val, true
}

func (om *OrderedMap[K, V]) Delete(key K) (val V) {
	for i, v := range om.List {
		if v == key {
			lhs := om.List[:i]
			if i < len(om.List)-1 {
				om.List = append(lhs, om.List[i+1:]...)
			} else {
				om.List = lhs
			}
			val = om.Pairs[key]
			delete(om.Pairs, key)
			return val
		}
	}
	return *new(V)
}
