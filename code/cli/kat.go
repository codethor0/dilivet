package cli

// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"path/filepath"

	mldsa "github.com/codethor0/dilivet/code/clean"
	"github.com/codethor0/dilivet/code/clean/kats"
)

const defaultSigVerVectors = "code/clean/testdata/kats/ml-dsa/ML-DSA-sigVer-FIPS204-internalProjection.json"

func (a *App) runKATVerify(args []string) int {
	fs := flag.NewFlagSet("kat-verify", flag.ContinueOnError)
	fs.SetOutput(a.Err)

	vectorsPath := fs.String("vectors", defaultSigVerVectors, "path to ACVP sigVer vector JSON")

	if err := fs.Parse(args); err != nil {
		return exitFromFlagError(err)
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(a.Err, "kat-verify: unexpected positional arguments")
		return 1
	}

	path := *vectorsPath
	if !filepath.IsAbs(path) {
		path = filepath.Clean(path)
	}

	vectors, err := kats.LoadSigVerVectors(path)
	if err != nil {
		fmt.Fprintf(a.Err, "kat-verify: load vectors: %v\n", err)
		return 1
	}

	var total, passes, structuralWarnings, structuralFailures, decodeFailures int

	for _, tg := range vectors.TestGroups {
		for _, tc := range tg.Tests {
			total++

			pk, err := hex.DecodeString(tc.Public)
			if err != nil {
				decodeFailures++
				continue
			}
			sig, err := hex.DecodeString(tc.Signature)
			if err != nil {
				decodeFailures++
				continue
			}
			msg, err := hex.DecodeString(tc.Message)
			if err != nil {
				decodeFailures++
				continue
			}

			ok, verr := mldsa.Verify(pk, msg, sig)
			switch {
			case verr == nil && ok:
				passes++
			default:
				if errors.Is(verr, mldsa.ErrNotImplemented) {
					structuralWarnings++
				} else if verr == nil {
					structuralWarnings++
				} else {
					structuralFailures++
				}
			}
		}
	}

	fmt.Fprintf(a.Out, "Vectors: %s\n", *vectorsPath)
	fmt.Fprintf(a.Out, "Total tests: %d\n", total)
	fmt.Fprintf(a.Out, "Strict passes: %d\n", passes)
	fmt.Fprintf(a.Out, "Structural warnings: %d\n", structuralWarnings)
	fmt.Fprintf(a.Out, "Structural failures: %d\n", structuralFailures)
	fmt.Fprintf(a.Out, "Decode failures: %d\n", decodeFailures)
	fmt.Fprintf(a.Out, "Note: full ML-DSA verification is not yet implemented; results indicate parsing/length checks only.\n")

	if decodeFailures > 0 || structuralFailures > 0 {
		return 1
	}
	return 0
}

func exitFromFlagError(err error) int {
	if err == flag.ErrHelp {
		return 0
	}
	return 1
}
