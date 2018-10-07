// The Little Big Go compiler adventure.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mewkiz/pkg/term"
)

var (
	// warn is a logger with the "warning:" prefix, which logs warning messages
	// to standard error.
	warn = log.New(os.Stderr, term.RedBold("warning:")+" ", 0)
)

func usage() {
	const use = `
Usage: lbg [OPTION]... [packages]

Little Big Go compiles and installs the packages named by the import paths.

For more about specifying packages, see 'go help packages'.
`
	fmt.Fprintln(os.Stderr, use[1:])
	flag.PrintDefaults()
}

func main() {
	// Parse command line arguments.
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	patterns := flag.Args()

	// Parse Go packages specified by patterns.
	pkgs, err := Parse(patterns...)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	// Compile Go packages into LLVM IR modules.
	modules, err := Compile(pkgs)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	for _, m := range modules {
		fmt.Println(m)
	}
}
