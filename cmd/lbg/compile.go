package main

import (
	"fmt"
	"go/ast"
	"log"
	"os"
	"sort"

	"github.com/llir/l/ir"
	"github.com/mewkiz/pkg/term"
	"golang.org/x/tools/go/packages"
)

var (
	// dbg2 is a logger with the "compile:" prefix, which logs debug messages to
	// standard error.
	dbg2 = log.New(os.Stderr, term.CyanBold("compile:")+" ", 0)
)

// Compile compiles the given parsed Go packages into LLVM IR modules.
func Compile(pkgs []*packages.Package) ([]*ir.Module, error) {
	dbg2.Println("packages (post-order):")
	c := NewCompiler(pkgs)

	// Compile pseudo-package builtin for predeclared identifiers.
	dbg2.Println("   id:", "builtin")
	c.modules["builtin"] = &ir.Module{
		SourceFilename: "builtin",
	}
	// TODO: handle "builtin"

	post := func(pkg *packages.Package) {
		dbg2.Println("   id:", pkg.ID)
		c.compile(pkg)
	}
	packages.Visit(pkgs, nil, post)

	// Return LLVM IR modules.
	var ids []string
	for id := range c.modules {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	var ms []*ir.Module
	for _, id := range ids {
		m := c.modules[id]
		ms = append(ms, m)
	}
	return ms, nil
}

// Compiler tracks information required to compile a set of Go packages.
type Compiler struct {
	// Maps from Go package ID to parsed Go package.
	pkgs map[string]*packages.Package
	// Map from Go package ID to output LLVM IR module.
	modules map[string]*ir.Module

	// reset after compiling each package.
	curPkg    *packages.Package
	curModule *ir.Module

	// reset after compiling each function.
	curFunc  *ir.Function
	curBlock *ir.BasicBlock
}

// NewCompiler returns a new compiler for the given parsed Go packages.
func NewCompiler(pkgs []*packages.Package) *Compiler {
	ps := make(map[string]*packages.Package)
	for _, pkg := range pkgs {
		ps[pkg.ID] = pkg
	}
	return &Compiler{
		pkgs:    ps,
		modules: make(map[string]*ir.Module),
	}
}

// === [ compile ] =============================================================

// compile compiles the given Go package into an LLVM IR module.
func (c *Compiler) compile(pkg *packages.Package) {
	// Create LLVM IR module and add scaffolding LLVM IR values for global
	// variable, function and type declarations.
	c.indexPackage(pkg)

	// Compile Go source files.
	for _, file := range pkg.Syntax {
		c.compileFile(file)
	}

	// Reset compiler state for the current Go package.
	c.resetPackage()
}

// compileFile compiles the given Go source file into an LLVM IR module.
func (c *Compiler) compileFile(file *ast.File) {
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			c.compileFuncDecl(decl)
		case *ast.GenDecl:
			c.compileGenDecl(decl)
		default:
			panic(fmt.Errorf("support for declaration %T not yet implemented", decl))
		}
	}
}

// --- [ function ] ------------------------------------------------------------

// compileFuncDecl compiles the given function declaration.
func (c *Compiler) compileFuncDecl(decl *ast.FuncDecl) {
	//pretty.Println("FuncDecl:", decl)
}

// compileGenDecl compiles the given generic declaration.
func (c *Compiler) compileGenDecl(decl *ast.GenDecl) {
	//pretty.Println("GenDecl:", decl)
}

// ### [ Helper functions ] ####################################################

// === [ index ] ===============================================================

// indexPackage creates an LLVM IR module for the given Go package, and adds
// scaffolding LLVM IR values for top-level global variable, function and type
// declarations.
func (c *Compiler) indexPackage(pkg *packages.Package) {
	// Create LLVM IR module.
	c.curPkg = pkg
	m := &ir.Module{
		SourceFilename: pkg.ID,
	}
	c.curModule = m
	c.modules[pkg.ID] = m

	// Create scaffolding LLVM IR values for global variable, function and type
	// declarations.
	for _, file := range pkg.Syntax {
		c.indexFile(file)
	}
}

// indexFile adds scaffolding LLVM IR values for top-level global variable,
// function and type declarations of the given Go file.
func (c *Compiler) indexFile(file *ast.File) {
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			c.indexFuncDecl(decl)
		case *ast.GenDecl:
			c.indexGenDecl(decl)
		default:
			panic(fmt.Errorf("support for declaration %T not yet implemented", decl))
		}
	}
}

// --- [ function ] ------------------------------------------------------------

// indexFuncDecl adds scaffolding LLVM IR values for the given function
// declaration.
func (c *Compiler) indexFuncDecl(decl *ast.FuncDecl) {
	//pretty.Println("FuncDecl:", decl)
}

// indexGenDecl adds scaffolding LLVM IR values for the given generic
// declaration.
func (c *Compiler) indexGenDecl(decl *ast.GenDecl) {
	//pretty.Println("GenDecl:", decl)
}

// === [ reset ] ===============================================================

// resetPackage resets the compiler state for the current Go package.
func (c *Compiler) resetPackage() {
	c.curPkg = nil
	c.curModule = nil
}
