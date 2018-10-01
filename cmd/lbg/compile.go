package main

import (
	"go/build"
	"log"
	"sort"

	"github.com/kr/pretty"
	"github.com/llir/l/ir"
	"github.com/llir/l/ir/types"
	"github.com/mewmew/lbg/internal/syntax"
	"github.com/pkg/errors"
)

// Compile compiles the given parsed Go packages into LLVM IR modules.
func Compile(pkgs map[string]*Package) (map[string]*ir.Module, error) {
	dbg.Println("compile:")
	c := NewCompiler(pkgs)

	// Compile pseudo-package builtin for predeclared identifiers.
	c.push("builtin")
	if err := c.Compile(); err != nil {
		return nil, errors.WithStack(err)
	}
	// TODO: Resolve predeclared identifiers of the universe scope.
	//c.Resolve()

	var pkgPaths []string
	for pkgPath := range c.pkgs {
		pkgPaths = append(pkgPaths, pkgPath)
	}
	sort.Strings(pkgPaths)
	for _, pkgPath := range pkgPaths {
		c.push(pkgPath)
		if err := c.Compile(); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return c.modules, nil
}

// Compiler tracks information required to compile a set of Go packages.
type Compiler struct {
	// Maps from Go package path to parsed Go package.
	pkgs map[string]*Package
	// Map from Go package path to output LLVM IR module.
	modules map[string]*ir.Module
	// Universe scope of resolved predeclared identifiers.
	universe *Scope
	// Stack of packages to compile; the package on top of the stack has no
	// unresolved dependencies.
	stack Stack
}

// NewCompiler returns a new compiler for the given parsed Go packages.
func NewCompiler(pkgs map[string]*Package) *Compiler {
	return &Compiler{
		pkgs:    pkgs,
		modules: make(map[string]*ir.Module),
	}
}

// Compile compiles the set of parsed Go packages in the stack of packages to
// compile, starting with the top element.
func (c *Compiler) Compile() error {
	dbg.Println("compile:")
	for !c.stack.Empty() {
		pkgPath := c.stack.Pop()
		dbg.Println("   pop:", pkgPath)
		pkg := c.pkgs[pkgPath]
		if err := c.compile(pkg); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// compile compiles the given parsed Go packages.
func (c *Compiler) compile(pkg *Package) error {
	c.modules[pkg.ImportPath] = &ir.Module{}
	return nil
}

// push pushes the given Go package and its transitive dependencies onto the top
// of the stack of packages to compile, those packages which are not yet
// compiled and not yet present in the stack.
func (c *Compiler) push(pkgPath string) {
	if c.stack.Contains(pkgPath) {
		return
	}
	if _, ok := c.modules[pkgPath]; ok {
		return
	}
	c.stack.Push(pkgPath)
	dbg.Println("   push:", pkgPath)
	pkg := c.pkgs[pkgPath]
	for _, importPkgPath := range pkg.Imports {
		c.push(importPkgPath)
	}
}

//func (c *Compiler) CompilePackage()

// ### [ cleanup below ] ###

/*
// compile compiles the given Go package.
func compile(pkg *build.Package) error {
	c := newCompiler(pkg)
	fmt.Printf("=== [ %s ] ========================\n", pkg.Name)
	var files []*syntax.File
	for _, goFile := range pkg.GoFiles {
		fmt.Printf("   %v\n", goFile)
		file, err := parseFile(pkg, goFile)
		if err != nil {
			return errors.WithStack(err)
		}
		files = append(files, file)
		//if err := syntax.Fdump(os.Stderr, file); err != nil {
		//	return errors.WithStack(err)
		//}
	}
	c.files = files
	//fmt.Println()
	c.compile()
	pretty.Println("module:", c.curModule)
	return nil
}
*/

// Compiler tracks information related to the compilation of a specific Go
// package.
type compiler struct {
	// Maps from Go package path to parsed Go package.
	pkgs map[string]*Package

	// Go package.
	pkg   *build.Package
	files []*syntax.File

	imports []*build.Package

	curModule *ir.Module
	curFunc   *ir.Function
	curBlock  *ir.BasicBlock
	// Maps from package import path to LLVM IR module.
	//modules map[string]*ir.Module
	// Maps from qualified identifier to LLVM IR value (function or global).
	//values map[string]value.Value
}

func newCompiler(pkg *build.Package) *compiler {
	return &compiler{
		pkg: pkg,
	}
}

func (c *compiler) compile() {
	// Initialize universal scope.
	//c.curScope = NewScope(universe)
	// Create module.
	c.curModule = &ir.Module{
		SourceFilename: c.pkg.ImportPath,
	}
	// TODO: implement identifier resolution
	//    - map[pos]value
	//    - *syntax.Name
	//    - *syntax.SelectorExpr
	// TODO: implement type resolution
	//    - map[pos]type
	// TODO: implement type checking
	for _, file := range c.files {
		for _, decl := range file.DeclList {
			switch decl := decl.(type) {
			case *syntax.VarDecl:
			// TODO: translate Go type to LLVM IR.
			//for _, name := range decl.NameList {
			//}
			//c.curModule.NewGlobalDef(name, init)
			case *syntax.FuncDecl:
				c.funcDecl(decl)
			default:
				log.Printf("support for %T not yet implemented", decl)
			}
			pretty.Println("decl:", decl)
		}
	}
}

func (c *compiler) funcDecl(decl *syntax.FuncDecl) {
	pretty.Println("func:", decl)
	pretty.Println("func type:", c.llType(decl.Type))
	// TODO: translate Go results type to LLVM IR.
	retType := types.Void
	f := c.curModule.NewFunction(decl.Name.Value, retType)
	for _, param := range decl.Type.ParamList {
		// TODO: use ir.NewParam, or even ir.Function.NewParam?
		// TODO: translate Go function parameter type to LLVM IR.
		typ := types.I32
		p := &ir.Param{
			ParamName: param.Name.Value,
			Typ:       typ,
		}
		f.Params = append(f.Params, p)
	}
}
