package main

import (
	"fmt"
	"go/build"
	"os"

	"github.com/mewmew/lbg/internal/syntax"
	"github.com/pkg/errors"
)

// compile compiles the given Go package.
func compile(pkg *build.Package) error {
	fmt.Println(pkg.Name)
	for _, goFile := range pkg.GoFiles {
		fmt.Printf("   %v\n", goFile)
		file, err := parseFile(pkg, goFile)
		if err != nil {
			return errors.WithStack(err)
		}
		if err := syntax.Fdump(os.Stderr, file); err != nil {
			return errors.WithStack(err)
		}
	}
	fmt.Println()
	return nil
}
