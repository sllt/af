package datastructure

import (
	"testing"

	"github.com/sllt/af/internal"
)

func TestLinkedStack_Push(t *testing.T) {
	t.Parallel()

	assert := internal.NewAssert(t, "TestLinkedStack_Push")

	stack := NewLinkedStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	values := stack.Data()
	size := stack.Size()

	assert.Equal([]int{3, 2, 1}, values)
	assert.Equal(3, size)
}

func TestLinkedStack_Pop(t *testing.T) {
	t.Parallel()

	assert := internal.NewAssert(t, "TestLinkedStack_Pop")

	stack := NewLinkedStack[int]()
	_, err := stack.Pop()
	assert.IsNotNil(err)

	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	topItem, err := stack.Pop()
	assert.IsNil(err)
	assert.Equal(3, *topItem)

	stack.Print()
	assert.Equal([]int{2, 1}, stack.Data())
}

func TestLinkedStack_Peak(t *testing.T) {
	t.Parallel()

	assert := internal.NewAssert(t, "TestLinkedStack_Peak")

	stack := NewLinkedStack[int]()
	_, err := stack.Peak()
	assert.IsNotNil(err)

	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	topItem, err := stack.Peak()
	assert.IsNil(err)
	assert.Equal(3, *topItem)

	assert.Equal([]int{3, 2, 1}, stack.Data())
}

func TestLinkedStack_Empty(t *testing.T) {
	t.Parallel()

	assert := internal.NewAssert(t, "TestLinkedStack_Empty")

	stack := NewLinkedStack[int]()
	assert.Equal(true, stack.IsEmpty())
	assert.Equal(0, stack.Size())

	stack.Push(1)
	assert.Equal(false, stack.IsEmpty())
	assert.Equal(1, stack.Size())

	stack.Clear()
	assert.Equal(true, stack.IsEmpty())
	assert.Equal(0, stack.Size())
}
