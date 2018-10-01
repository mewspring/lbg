package main

import (
	"fmt"
	"go/build"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kr/pretty"
	"github.com/mewmew/lbg/internal/syntax"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/buildutil"
)

// Parse parses the Go packages specified by the given patterns.
func Parse(patterns []string) (map[string]*Package, error) {
	p := NewParser(&build.Default)
	pkgPaths, err := expandPatterns(p.ctxt, patterns)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, pkgPath := range pkgPaths {
		p.push(pkgPath)
	}
	if err := p.Parse(); err != nil {
		return nil, errors.WithStack(err)
	}
	return p.pkgs, nil
}

// Parser parses Go packages.
type Parser struct {
	// Tracks build context information (e.g. Go source directories).
	ctxt *build.Context
	// Maps from Go package path to parsed Go package.
	pkgs map[string]*Package
	// Queue of packages to parse.
	queue Queue
}

// NewParser returns a new parser for parsing Go packages.
func NewParser(ctxt *build.Context) *Parser {
	return &Parser{
		ctxt: ctxt,
		pkgs: make(map[string]*Package),
	}
}

// Parse parses the queued Go packages and their transitive imports.
func (p *Parser) Parse() error {
	for !p.queue.Empty() {
		pkgPath := p.queue.Pop()
		pkg, err := parsePkg(p.ctxt, pkgPath)
		if err != nil {
			if _, ok := errors.Cause(err).(*build.NoGoError); ok {
				// Skip directories without Go files.
				//log.Printf("skipping directory %q with no Go files", e.Dir)
				continue
			}
			return errors.WithStack(err)
		}
		p.pkgs[pkgPath] = pkg
		for _, importPkgPath := range pkg.Imports {
			p.push(importPkgPath)
		}
	}
	return nil
}

// push pushes the given Go package path onto the queue of packages to parse, if
// the package is not yet parsed and not yet present in the queue.
func (p *Parser) push(pkgPath string) {
	if p.queue.Contains(pkgPath) {
		return
	}
	if _, ok := p.pkgs[pkgPath]; ok {
		return
	}
	p.queue.Push(pkgPath)
}

// ### [ Helper functions ] ####################################################

// expandPatterns returns the Go package paths specified by the given patterns.
func expandPatterns(ctxt *build.Context, patterns []string) ([]string, error) {
	// Note, relative patterns (e.g. "." and "../") are not yet well supported by
	// buildutil.ExpandPatterns. There is a TODO in that package to extend
	// support, but for now, we use our own implementation.
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
	return pkgPaths, nil
}

// TODO: remove fixPatterns once buildutil.ExpandPatterns handles relative
// patterns.

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
			p := pattern
			parents := 0
			for ; strings.HasPrefix(p, "../"); p = p[len("../"):] {
				parents++
			}
			if p == ".." {
				parents++
				p = ""
			}
			absPath, err := filepath.Abs(strings.Repeat("../", parents))
			if err != nil {
				return nil, errors.WithStack(err)
			}
			pkgPath, err := findPkgPath(ctxt, absPath)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			ps = append(ps, strings.Join([]string{pkgPath, p}, "/"))
		// Use pattern as is.
		default:
			ps = append(ps, pattern)
		}
	}
	pretty.Println("patterns:", ps)
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
func parsePkg(ctxt *build.Context, pkgPath string) (*Package, error) {
	for _, srcDir := range ctxt.SrcDirs() {
		goPkg, err := ctxt.Import(pkgPath, srcDir, build.ImportComment)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		pkg := &Package{
			Package: goPkg,
			files:   make(map[string]*syntax.File),
		}
		for _, goFile := range pkg.GoFiles {
			file, err := parseFile(goPkg, goFile)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			pkg.files[goFile] = file
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
