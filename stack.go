package kiteroot

type Stack []*Element

func (s *Stack) Push(e *Element) {
	*s = append(*s, e)
}

func (s *Stack) Top() (e *Element) {
	if s.Len() == 0 {
		return nil
	}
	return (*s)[s.Len()-1]
}

func (s *Stack) Pop() (e *Element) {
	if s.Len() > 0 {
		n := s.Len() - 1
		e, *s = (*s)[n], (*s)[:n]
		return
	}
	return nil
}

func (s *Stack) Len() int {
	return len(*s)
}
