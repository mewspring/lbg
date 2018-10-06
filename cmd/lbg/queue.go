package main

// Queue is a queue of package path elements.
//
// The zero value for Queue is an empty queue ready to use.
type Queue []Elem

// Elem is a package path element.
type Elem struct {
	// Package path.
	PkgPath string
	// Importer directory (used if package has a relative import or is in vendor
	// directory); empty if the package is compiled stand-alone and not imported
	// by another package.
	ImporterDir string
}

// Push pushes elem onto the end of the queue.
func (q *Queue) Push(elem Elem) {
	*q = append(*q, elem)
}

// Pop pops and returns the first element of the queue.
func (q *Queue) Pop() Elem {
	if q.Empty() {
		panic("invalid call to pop; empty queue")
	}
	elem := (*q)[0]
	*q = (*q)[1:]
	return elem
}

// Contains reports whether the queue contains elem.
func (q *Queue) Contains(elem Elem) bool {
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
