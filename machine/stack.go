package machine

import (
	"errors"
)

type Stack struct {
	top    *StackNode
	length int
}

type StackNode struct {
	value interface{}
	prev  *StackNode
}

// Create a new stack
func NewStack() *Stack {
	return &Stack{nil, 0}
}

// Return the number of items in the stack
func (s *Stack) Len() int {
	return s.length
}

// View the top item on the stack
func (s *Stack) Peek() interface{} {
	if s.length == 0 {
		return nil
	}
	return s.top.value
}

// Try to pop the top item off the stack and return it
func (s *Stack) TryPop() (interface{}, error) {
	if s.length == 0 {
		return nil, errors.New("Stack is empty")
	}

	node := s.top
	s.top = node.prev
	s.length--
	return node.value, nil
}

// Pop the top item off the stack and return it
func (s *Stack) Pop() interface{} {
	if v, err := s.TryPop(); err == nil {
		return v
	} else {
		panic(err)
	}
}

// Push a value onto the top of the stack
func (s *Stack) Push(value interface{}) {
	node := &StackNode{value, s.top}
	s.top = node
	s.length++
}
