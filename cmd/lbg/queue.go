package main

// Queue is a queue of strings.
//
// The zero value for Queue is an empty queue ready to use.
type Queue []string

// Push pushes elem onto the end of the queue.
func (q *Queue) Push(elem string) {
	*q = append(*q, elem)
}

// Pop pops and returns the first element of the queue.
func (q *Queue) Pop() string {
	if q.Empty() {
		panic("invalid call to pop; empty queue")
	}
	elem := (*q)[0]
	*q = (*q)[1:]
	return elem
}

// Contains reports whether the queue contains elem.
func (q *Queue) Contains(elem string) bool {
	for _, e := range *q {
		if e == elem {
			return true
		}
	}
	return false
}

// Empty reports whether the queue is empty.
func (q *Queue) Empty() bool {
	return len(*q) == 0
}
