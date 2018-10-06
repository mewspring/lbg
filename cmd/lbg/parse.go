package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mewkiz/pkg/term"
	"github.com/mewmew/lbg/internal/syntax"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/buildutil"
)

var (
	// dbg1 is a logger with the "parse:" prefix, which logs debug messages to
	// standard error.
	dbg1 = log.New(os.Stderr, term.MagentaBold("parse:")+" ", 0)
)

// Parse parses the Go packages specified by the given patterns.
func Parse(patterns []string) (map[string]*Package, error) {
	dbg1.Println("patterns:")
	for _, pattern := range patterns {
		dbg1.Printf("   %v", pattern)
	}
	p := NewParser(&build.Default)
	pkgPaths, err := expandPatterns(p.ctxt, patterns)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, pkgPath := range pkgPaths {
		p.push(Elem{PkgPath: pkgPath})
	}
	// Parse pseudo-package builtin for predeclared identifiers.
	p.push(Elem{PkgPath: "builtin"})
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
		elem := p.queue.Pop()
		pkg, err := parsePkg(p.ctxt, elem.PkgPath, elem.ImporterDir)
		if err != nil {
			if _, ok := errors.Cause(err).(*build.NoGoError); ok {
				// Skip directories without Go files.
				//log.Printf("skipping directory %q with no Go files", e.Dir)
				continue
			}
			return errors.WithStack(err)
		}
		p.pkgs[elem.PkgPath] = pkg
		// TODO: check if there exists an exports data file for the Go package, to
		// avoid re-parsing.
		for _, importPkgPath := range pkg.Imports {
			elem := Elem{
				PkgPath:     importPkgPath,
				ImporterDir: pkg.Dir,
			}
			p.push(elem)
		}
	}
	return nil
}

// push pushes the given Go package path onto the queue of packages to parse, if
// the package is not yet parsed and not yet present in the queue.
func (p *Parser) push(elem Elem) {
	if p.queue.Contains(elem) {
		return
	}
	if _, ok := p.pkgs[elem.PkgPath]; ok {
		return
	}
	p.queue.Push(elem)
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
	dbg1.Println("patterns (fixed):")
	for _, p := range ps {
		dbg1.Printf("   %v", p)
	}
	m := buildutil.ExpandPatterns(ctxt, ps)
	var pkgPaths []string
	for pkgPath := range m {
		pkgPaths = append(pkgPaths, pkgPath)
	}
	sort.Strings(pkgPaths)
	dbg1.Println("package paths:")
	for _, pkgPath := range pkgPaths {
		dbg1.Printf("   %v", pkgPath)
	}
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

// parsePkg parses the given Go package. The importer directory is used if
// package has a relative import or is in vendor directory. An empty import
// directory is used if the package is compiled stand-alone and not imported by
// another package.
func parsePkg(ctxt *build.Context, pkgPath string, importerDir string) (*Package, error) {
	dbg1.Println("parsing package:", pkgPath)
	if pkgPath == "C" {
		// TODO: figure out how to support cgo.
		return &Package{Package: &build.Package{ImportPath: "C"}}, nil
	}
	goPkg, err := ctxt.Import(pkgPath, importerDir, build.ImportComment)
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

// parseFile parses the given Go file.
func parseFile(pkg *build.Package, goFile string) (*syntax.File, error) {
	dbg1.Println("   parsing file:", goFile)
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
