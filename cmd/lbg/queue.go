package main

// Queue is a queue of strings.
//
// The zero value for Queue is an empty queue ready to use.
type Queue []string

// Push pushes s onto the queue.
func (q *Queue) Push(s string) {
	*q = append(*q, s)
}

// Pop pops and returns an element from the queue.
func (q *Queue) Pop() string {
	if q.Empty() {
		panic("invalid call to pop; empty queue")
	}
	s := (*q)[0]
	*q = (*q)[1:]
	return s
}

// Contains reports whether the queue contains s.
func (q *Queue) Contains(s string) bool {
	for _, r := range *q {
		if r == s {
			return true
		}
	}
	return false
}

// Empty reports whether the queue is empty.
func (q *Queue) Empty() bool {
	return len(*q) == 0
}
