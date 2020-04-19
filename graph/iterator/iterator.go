package iterator

// New returns a new iterator.
func New(item ...interface{}) *Iterator {
	return &Iterator{data: item}
}

// Iterator is an iterator.
type Iterator struct {
	current int
	data    []interface{}
}

// Value returns the current value.
func (it *Iterator) Value() interface{} {
	return it.data[it.current]
}

// Next progresses the iterator returning true if there are still items to iterator over.
func (it *Iterator) Next() bool {
	it.current++

	if it.current >= len(it.data) {
		return false
	}

	return true
}
