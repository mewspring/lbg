package main

import (
	"fmt"
	"go/build"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mewmew/lbg/internal/syntax"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/buildutil"
)

// parse parses the set of Go packages specified by the given patterns.
func parse(patterns []string) ([]*build.Package, error) {
	// Note, relative patterns (e.g. "." and "../") are not yet well supported by
	// buildutil.ExpandPatterns. There is a TODO in that package to extend
	// support, but for now, we use our own implementation.
	ctxt := &build.Default
	ps, err := fixPatterns(ctxt, patterns)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	m := buildutil.ExpandPatterns(ctxt, ps)
	var pkgPaths []string
	for pkgPath := range m {
		pkgPaths = append(pkgPaths, pkgPath)
	}
	sort.Strings(pkgPaths)
	var pkgs []*build.Package
	for _, pkgPath := range pkgPaths {
		pkg, err := parsePkg(ctxt, pkgPath)
		if err != nil {
			if _, ok := errors.Cause(err).(*build.NoGoError); ok {
				// Skip directories without Go files.
				continue
			}
			return nil, errors.WithStack(err)
		}
		pkgs = append(pkgs, pkg)
	}
	return pkgs, nil
}

// fixPatterns translates relative import patterns to qualified import patterns.
func fixPatterns(ctxt *build.Context, patterns []string) ([]string, error) {
	var ps []string
	for _, pattern := range patterns {
		switch {
		// Relative to current directory.
		case pattern == "." || strings.HasPrefix(pattern, "./"):
			absPath, err := filepath.Abs(".")
			if err != nil {
				return nil, errors.WithStack(err)
			}
			pkgPath, err := findPkgPath(ctxt, absPath)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			p := fmt.Sprintf("%s%s", pkgPath, pattern[len("."):])
			ps = append(ps, p)
		// Relative to parent directory.
		case pattern == ".." || strings.HasPrefix(pattern, "../"):
			absPath, err := filepath.Abs("..")
			if err != nil {
				return nil, errors.WithStack(err)
			}
			pkgPath, err := findPkgPath(ctxt, absPath)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			p := fmt.Sprintf("%s%s", pkgPath, pattern[len(".."):])
			ps = append(ps, p)
		// Use pattern as is.
		default:
			ps = append(ps, pattern)
		}
	}
	return ps, nil
}

// findPkgPath returns the qualified package path (as relative to a Go src
// directory) for the given absolute path.
func findPkgPath(ctxt *build.Context, absPath string) (string, error) {
	for _, srcDir := range ctxt.SrcDirs() {
		if filepath.HasPrefix(absPath, srcDir) {
			pkgPath, err := filepath.Rel(srcDir, absPath)
			if err != nil {
				return "", errors.WithStack(err)
			}
			return pkgPath, nil
		}
	}
	return "", errors.Errorf("unable to locate %q in Go src directories `%s`", absPath, ctxt.SrcDirs())
}

// parsePkg parses the given Go package.
func parsePkg(ctxt *build.Context, pkgPath string) (*build.Package, error) {
	for _, srcDir := range ctxt.SrcDirs() {
		pkg, err := ctxt.Import(pkgPath, srcDir, build.ImportComment)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return pkg, nil
	}
	return nil, errors.Errorf("unable to locate package %q in Go src directories `%s`", pkgPath, ctxt.SrcDirs())
}

// parseFile parses the given Go file.
func parseFile(pkg *build.Package, goFile string) (*syntax.File, error) {
	absGoPath := filepath.Join(pkg.Dir, goFile)
	mode := syntax.CheckBranches
	errh := func(err error) {
		log.Printf("compile error: %v", err)
	}
	file, err := syntax.ParseFile(absGoPath, errh, nil, mode)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return file, nil
}
