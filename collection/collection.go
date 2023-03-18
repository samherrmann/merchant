package collection

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{m: make(map[K]V)}
}

type OrderedMap[K comparable, V any] struct {
	m     map[K]V
	order []K
}

func (m *OrderedMap[K, V]) Raw() map[K]V {
	return m.m
}

func (m *OrderedMap[K, V]) Set(k K, v V) {
	if _, exist := m.m[k]; !exist {
		m.order = append(m.order, k)
	}
	m.m[k] = v
}

func (m *OrderedMap[K, V]) Get(k K) (V, bool) {
	v, exists := m.m[k]
	return v, exists
}

func (m *OrderedMap[K, V]) Slice() []V {
	slice := []V{}
	for _, k := range m.order {
		slice = append(slice, m.m[k])
	}
	return slice
}

// PadSliceRight pads slice at the end.
func PadSliceRight[T any](slice []T, length int) []T {
	if len(slice) >= length {
		return slice
	}
	return append(slice, make([]T, length-len(slice))...)
}
