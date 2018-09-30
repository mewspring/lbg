package main

import (
	"fmt"
	"go/build"
	"log"

	"github.com/kr/pretty"
	"github.com/llir/l/ir"
	"github.com/llir/l/ir/types"
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
	pretty.Println("module:", c.curModule)
	return nil
}

// Compiler tracks information related to the compilation of a specific Go
// package.
type compiler struct {
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
