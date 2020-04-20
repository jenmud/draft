package iterator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {
	item1 := 1
	item2 := 2
	item3 := 3

	iter := Iterator{data: []interface{}{item1, item2, item3}}

	actual1 := iter.Value()
	assert.Equal(t, item1, actual1)
	assert.Equal(t, true, iter.Next())

	actual2 := iter.Value()
	assert.Equal(t, item2, actual2)
	assert.Equal(t, true, iter.Next())

	actual3 := iter.Value()
	assert.Equal(t, item3, actual3)
	assert.Equal(t, false, iter.Next())
}

func TestSize(t *testing.T) {
	iter := Iterator{data: []interface{}{1, 2, 3}}
	assert.Equal(t, 3, iter.Size())
}
