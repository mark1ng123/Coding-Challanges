package file_utilities

import "fmt"

type Stack[T any] struct {
	items []T
}

func NewStack[T any]() Stack[T] {
	return Stack[T]{items: make([]T, 0)}
}

func (s *Stack[T]) isEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack[T]) Push(data T) {
	s.items = append(s.items, data)
}

func (s *Stack[T]) Pop() {
	if s.isEmpty() {
		return
	}

	s.items = s.items[:len(s.items)-1]
}

func (s *Stack[T]) Top() (T, error) {
	if s.isEmpty() {
		var zeroValue T
		return zeroValue, fmt.Errorf("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}

func (s *Stack[T]) Print() {
	for _, item := range s.items {
		fmt.Print(item, " ")
	}
	fmt.Println()
}
