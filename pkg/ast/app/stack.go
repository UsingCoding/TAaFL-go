package app

type (
	StackElementPointer uint
)

type Stack []Node

func (s *Stack) Add(node Node) StackElementPointer {
	*s = append(*s, node)
	// There is never will be nil point deference since we append to stack earlier
	return *s.currentTopPointer() - 1
}

func (s Stack) currentTopPointer() *StackElementPointer {
	stackLen := uint(len(s))
	if stackLen == 0 {
		return nil
	}
	return (*StackElementPointer)(&stackLen)
}
