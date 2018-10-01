// The Little Big Go compiler adventure.
package main

import (
	"flag"
	"log"
)

func main() {
	// Parse command line arguments.
	flag.Parse()
	patterns := flag.Args()

	// Parse Go packages specified by patterns.
	pkgs, err := Parse(patterns)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	// Compile Go packages.
	c := NewCompiler(pkgs)
	if err := c.Compile(); err != nil {
		log.Fatalf("%+v", err)
	}
}
