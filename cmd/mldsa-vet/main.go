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
		Name:    "mldsa-vet",
		Version: version,
		Out:     stdout,
		Err:     stderr,
	}
	return app.Run(args)
}
