// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package main

import (
	"io"
	"os"

	"github.com/codethor0/dilivet/code/cli"
)

var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	app := &cli.App{
		Name:    "dilivet",
		Version: version,
		Out:     stdout,
		Err:     stderr,
	}
	return app.Run(args)
}
