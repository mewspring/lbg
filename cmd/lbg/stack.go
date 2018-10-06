package main

// Stack is a stack of string elements.
//
// The zero value for Stack is an empty stack ready to use.
type Stack []string

// Push pushes elem onto the stack the top of the stack.
func (s *Stack) Push(elem string) {
	*s = append(*s, elem)
}

// Pop pops and returns the top element of the stack.
func (s *Stack) Pop() string {
	if s.Empty() {
		panic("invalid call to pop; empty stack")
	}
	elem := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return elem
}

// Contains reports whether the stack contains elem.
func (s *Stack) Contains(elem string) bool {
	for _, e := range *s {
		if e == elem {
			return true
		}
	}
	return false
}

// Empty reports whether the stack is empty.
func (s *Stack) Empty() bool {
	return len(*s) == 0
}
