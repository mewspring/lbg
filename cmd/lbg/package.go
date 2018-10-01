package main

import (
	"go/build"

	"github.com/mewmew/lbg/internal/syntax"
)

// Package is a parsed Go package.
type Package struct {
	// Go package.
	*build.Package
	// Parsed files of Go package.
	files map[string]*syntax.File
}
