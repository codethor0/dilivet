package main

import (
	"os"

	"github.com/codethor0/dilivet/code/cli"
)

var version = "dev"

func main() {
	app := &cli.App{
		Name:    "dilivet",
		Version: version,
	}
	os.Exit(app.Run(os.Args[1:]))
}
