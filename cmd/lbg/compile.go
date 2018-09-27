package main

import (
	"fmt"
	"go/build"

	"github.com/kr/pretty"
	"github.com/llir/l/ir"
	"github.com/mewmew/lbg/internal/syntax"
	"github.com/pkg/errors"
)

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
	return nil
}

// Compiler tracks information related to the compilation of a specific Go
// package.
type compiler struct {
	// Go package.
	pkg   *build.Package
	files []*syntax.File

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
		pkg:       pkg,
		curModule: &ir.Module{}, // TODO: use ir.NewModule()
		curFunc:   nil,
		curBlock:  nil,
	}
}

func (c *compiler) compile() {
	// TODO: implement identifier resolution
	// TODO: implement type resolution
	// TODO: implement type checking
	for _, file := range c.files {
		for _, decl := range file.DeclList {
			pretty.Println("decl:", decl)
		}
	}
}
