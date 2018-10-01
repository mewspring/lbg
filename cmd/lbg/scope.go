package main

import (
	"github.com/mewmew/lbg/internal/syntax"
)

// A Scope tracks the set of declarations in a lexical scope and links to its
// containing scope.
type Scope struct {
	// Containing scope.
	outer *Scope
	// Maps from Go identifier to the corresponding Go definition.
	decls map[string]syntax.Decl
}

// NewScope returns a new scope nested in the given outer scope.
func NewScope(outer *Scope) *Scope {
	return &Scope{
		outer: outer,
		decls: make(map[string]syntax.Decl),
	}
}

/*
	// Maps from Go type name to the corresponding LLVM IR type.
	types map[string]types.Type
	// Maps from Go operand name to the corresponding LLVM IR value.
	values map[string]value.Value
*/
