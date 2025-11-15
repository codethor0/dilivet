// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package cli

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"

	mldsa "github.com/codethor0/dilivet/code/clean"
	"github.com/codethor0/dilivet/code/clean/kats"
	"github.com/codethor0/dilivet/code/diag"
)

const defaultSigVerVectors = "code/clean/testdata/kats/ml-dsa/ML-DSA-sigVer-FIPS204-internalProjection.json"

func (a *App) runKATVerify(args []string) int {
	fs := flag.NewFlagSet("kat-verify", flag.ContinueOnError)
	fs.SetOutput(a.Err)

	vectorsPath := fs.String("vectors", defaultSigVerVectors, "path to ACVP sigVer vector JSON")
	jsonOut := fs.Bool("json", false, "emit machine-readable JSON summary")

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

	report := diag.Report{}

	for _, tg := range vectors.TestGroups {
		for _, tc := range tg.Tests {
			report.TotalTests++

			pk, err := hex.DecodeString(tc.Public)
			if err != nil {
				report.DecodeFailures++
				continue
			}
			sig, err := hex.DecodeString(tc.Signature)
			if err != nil {
				report.DecodeFailures++
				continue
			}
			msg, err := hex.DecodeString(tc.Message)
			if err != nil {
				report.DecodeFailures++
				continue
			}

			ok, verr := mldsa.Verify(pk, msg, sig)
			switch {
			case verr == nil && ok:
				report.StrictPasses++
			case verr == nil && !ok:
				// Verification returned false (signature rejected)
				report.StructuralFailures++
			default:
				// Verification error (unpacking, format, etc.)
				report.StructuralFailures++
			}
		}
	}

	if *jsonOut {
		payload := struct {
			Vectors string `json:"vectors"`
			diag.Report
			Note string `json:"note,omitempty"`
		}{
			Vectors: path,
			Report:  report,
			Note:    "Full ML-DSA verification is implemented; counts reflect complete cryptographic verification.",
		}
		enc := json.NewEncoder(a.Out)
		enc.SetIndent("", "  ")
		if err := enc.Encode(payload); err != nil {
			fmt.Fprintf(a.Err, "kat-verify: encode json: %v\n", err)
			return 1
		}
	} else {
		fmt.Fprintf(a.Out, "Vectors: %s\n", *vectorsPath)
		fmt.Fprintf(a.Out, "Total tests: %d\n", report.TotalTests)
		fmt.Fprintf(a.Out, "Strict passes: %d\n", report.StrictPasses)
		fmt.Fprintf(a.Out, "Structural warnings: %d\n", report.StructuralWarnings)
		fmt.Fprintf(a.Out, "Structural failures: %d\n", report.StructuralFailures)
		fmt.Fprintf(a.Out, "Decode failures: %d\n", report.DecodeFailures)
		fmt.Fprintf(a.Out, "Note: full ML-DSA verification is implemented; results indicate complete cryptographic verification.\n")
	}

	if report.DecodeFailures > 0 || report.StructuralFailures > 0 {
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
