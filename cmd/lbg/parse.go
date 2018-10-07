package main

import (
	"log"
	"os"

	"github.com/mewkiz/pkg/term"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

var (
	// dbg1 is a logger with the "parse:" prefix, which logs debug messages to
	// standard error.
	dbg1 = log.New(os.Stderr, term.MagentaBold("parse:")+" ", 0)
)

// Parse parses the Go packages specified by the given patterns.
func Parse(patterns ...string) ([]*packages.Package, error) {
	// Expand patterns into Go package paths.
	dbg1.Println("patterns:")
	for _, pattern := range patterns {
		dbg1.Printf("   %v", pattern)
	}
	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	dbg1.Println("packages (pre-order):")
	pre := func(pkg *packages.Package) bool {
		dbg1.Println("   id:", pkg.ID)
		return true
	}
	packages.Visit(pkgs, pre, nil)
	return pkgs, nil
}
