package main

import (
	"fmt"
	"go/build"
)

// compile compiles the given Go package.
func compile(pkg *build.Package) error {
	fmt.Println(pkg.Name)
	for _, goFile := range pkg.GoFiles {
		fmt.Printf("   %v\n", goFile)
	}
	fmt.Println()
	return nil
}
