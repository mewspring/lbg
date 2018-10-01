package main

import "github.com/mewmew/lbg/internal/syntax"

// Resolve resolves the identifiers of the parsed Go packages.
func (c *Compiler) Resolve() {
	// Add predeclared identifiers to universe scope.
	//
	// ref: https://golang.org/ref/spec#Predeclared_identifiers
	universe := NewScope(nil)
	builtin, ok := c.pkgs["builtin"]
	if !ok {
		panic("unable to create universe scope; cannot find package builtin")
	}
	for _, file := range builtin.files {
		for _, decl := range file.DeclList {
			// TODO: skip adding builtin.FloatType, etc to universe scope.
			switch decl := decl.(type) {
			case *syntax.ConstDecl:
				// ConstDecl only used for `true`, `false` and `iota`, for each of
				// which len(NameList) = 1.
				universe.decls[decl.NameList[0].Value] = decl
			case *syntax.FuncDecl:
				universe.decls[decl.Name.Value] = decl
			case *syntax.TypeDecl:
				universe.decls[decl.Name.Value] = decl
			case *syntax.VarDecl:
				// VarDecl only used for `nil`, for which len(NameList) = 1.
				universe.decls[decl.NameList[0].Value] = decl
			}
		}
	}
}
