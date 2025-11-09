// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

// Package cli provides shared CLI application logic for dilivet and mldsa-vet.
//
// This package eliminates code duplication between the two CLI entrypoints
// by providing a common App structure and execution logic.
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// App represents a CLI application with common functionality.
type App struct {
	Name    string // Binary name (e.g., "dilivet" or "mldsa-vet")
	Version string // Version string (injected via ldflags)
	Out     io.Writer
	Err     io.Writer
}

// Run executes the CLI application with the given arguments.
//
// It parses flags, handles common options (version, help), and executes
// the appropriate command. Returns an exit code suitable for os.Exit().
func (a *App) Run(args []string) int {
	// Default to os.Stdout/Stderr if not set
	if a.Out == nil {
		a.Out = os.Stdout
	}
	if a.Err == nil {
		a.Err = os.Stderr
	}

	// Create flag set
	fs := flag.NewFlagSet(a.Name, flag.ContinueOnError)
	fs.SetOutput(a.Err)

	version := fs.Bool("version", false, "print version and exit")
	help := fs.Bool("help", false, "show help message")

	// Parse flags
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	// Handle version flag
	if *version {
		fmt.Fprintln(a.Out, a.Version)
		return 0
	}

	// Handle help flag
	if *help {
		a.PrintHelp()
		return 0
	}

	// Handle subcommands
	if fs.NArg() > 0 {
		cmd := fs.Arg(0)
		args := fs.Args()[1:]

		switch cmd {
		case "verify":
			return a.runVerify(args)
		default:
			fmt.Fprintf(a.Err, "unknown command %q\n", cmd)
			return 1
		}
	}

	// Default behavior: print usage message
	fmt.Fprintf(a.Out, "%s - ML-DSA vetting tool. Use -help for available commands.\n", a.Name)
	return 0
}

// PrintHelp displays comprehensive usage information.
func (a *App) PrintHelp() {
	fmt.Fprintf(a.Out, `%s v%s - ML-DSA Signature Diagnostics Tool

DESCRIPTION:
    A toolkit for ML-DSA (Dilithium-like) signature diagnostics and vetting.
    Provides test harnesses, known-answer vectors, and CLI tools to validate
    post-quantum cryptographic implementations.

USAGE:
    %s [OPTIONS] <command>

COMMANDS:
    verify      Validate an ML-DSA signature against a public key

OPTIONS:
    -version    Print version and exit
    -help       Show this help message

EXAMPLES:
    %s -version
        Print the version number

    %s verify -pub pk.hex -sig sig.hex -msg msg.bin
        Verify a signature using hex-encoded key/signature files

DOCUMENTATION:
    GitHub: https://github.com/codethor0/dilivet
    Issues: https://github.com/codethor0/dilivet/issues

LICENSE:
    MIT License - see LICENSE file for details
`, a.Name, a.Version, a.Name, a.Name, a.Name)
}
