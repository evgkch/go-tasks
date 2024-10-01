package main

import (
	"fmt"
)

type Stack struct {
	items []int
}

func (s *Stack) Pop() (int, error) {
	if len(s.items) == 0 {
		return 0, fmt.Errorf("stack empty")
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, nil
}

func (s *Stack) Push(v int) {
	s.items = append(s.items, v)
}
