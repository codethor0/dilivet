package main

import (
	"fmt"
	"os"

	"github.com/codethor0/ml-dsa-debug-whitepaper/internal/version"
)

// Fires before main(); prints version if requested and exits early.
func init() {
	for _, a := range os.Args[1:] {
		if a == "-version" || a == "--version" || a == "version" {
			fmt.Println(version.String())
			os.Exit(0)
		}
	}
}
