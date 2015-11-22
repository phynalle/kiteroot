package kiteroot

// Stack implements the stack container
type Stack []*Element

// Push adds the given element to the stack.
func (s *Stack) Push(e *Element) {
	*s = append(*s, e)
}

// Top returns the last element in the stack if it is not empty.
// Otherwise, nil is returned.
func (s *Stack) Top() (e *Element) {
	if s.Len() == 0 {
		return nil
	}
	return (*s)[s.Len()-1]
}

// Pop removes top element from stack and returns it. if stack is empty,
// nil is returned.
func (s *Stack) Pop() (e *Element) {
	if s.Len() > 0 {
		n := s.Len() - 1
		e, *s = (*s)[n], (*s)[:n]
		return
	}
	return nil
}

// Len returns the number of element currently in the stack.
func (s *Stack) Len() int {
	return len(*s)
}
