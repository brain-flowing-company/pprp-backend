package utils

type Stack[T any] struct {
	s    []T
	size int
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		s:    []T{},
		size: 0,
	}
}

func (stk *Stack[T]) Push(v T) {
	if stk.size == len(stk.s) {
		stk.s = append(stk.s, v)
	} else {
		stk.s[stk.size] = v
	}

	stk.size++
}

func (stk *Stack[T]) Pop() {
	if stk.size > 0 {
		stk.size--
	}
}

func (stk *Stack[T]) Top() T {
	return stk.s[stk.size-1]
}

func (stk *Stack[T]) Seek() []T {
	return stk.s[:stk.size]
}
