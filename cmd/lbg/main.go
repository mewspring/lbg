// The Little Big Go compiler adventure.
package main

import (
	"flag"
	"log"
)

func main() {
	flag.Parse()
	patterns := flag.Args()
	pkgs, err := parse(patterns)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	for _, pkg := range pkgs {
		if err := compile(pkg); err != nil {
			log.Fatalf("%+v", err)
		}
	}
}
