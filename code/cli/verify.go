// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package cli

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	mldsa "github.com/codethor0/dilivet/code/clean"
)

const (
	formatHex = "hex"
	formatRaw = "raw"
)

func (a *App) runVerify(args []string) int {
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	fs.SetOutput(a.Err)

	pubPath := fs.String("pub", "", "path to ML-DSA public key")
	sigPath := fs.String("sig", "", "path to ML-DSA signature")
	msgPath := fs.String("msg", "", "path to message bytes")
	pubFormat := fs.String("pub-format", formatHex, "format of public key file (hex|raw)")
	sigFormat := fs.String("sig-format", formatHex, "format of signature file (hex|raw)")
	msgFormat := fs.String("msg-format", formatRaw, "format of message file (hex|raw)")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 1
	}

	if fs.NArg() != 0 {
		fmt.Fprintln(a.Err, "verify: unexpected positional arguments")
		return 1
	}

	if *pubPath == "" || *sigPath == "" || *msgPath == "" {
		fmt.Fprintln(a.Err, "verify: -pub, -sig, and -msg are required")
		return 1
	}

	pub, err := loadData(*pubPath, *pubFormat)
	if err != nil {
		fmt.Fprintf(a.Err, "verify: read public key: %v\n", err)
		return 1
	}

	sig, err := loadData(*sigPath, *sigFormat)
	if err != nil {
		fmt.Fprintf(a.Err, "verify: read signature: %v\n", err)
		return 1
	}

	msg, err := loadData(*msgPath, *msgFormat)
	if err != nil {
		fmt.Fprintf(a.Err, "verify: read message: %v\n", err)
		return 1
	}

	valid, verr := mldsa.Verify(pub, msg, sig)
	switch {
	case errors.Is(verr, mldsa.ErrNotImplemented):
		fmt.Fprintf(a.Out, "Structural checks passed (%s): full ML-DSA verification not yet implemented.\n", a.Name)
		return 0
	case verr != nil:
		fmt.Fprintf(a.Err, "verification failed: %v\n", verr)
		return 1
	default:
		if !valid {
			fmt.Fprintln(a.Err, "verification failed: signature rejected")
			return 1
		}
		fmt.Fprintln(a.Out, "Signature verified successfully.")
		return 0
	}
}

func loadData(path, format string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(format) {
	case formatHex:
		clean := stripWhitespace(string(data))
		if clean == "" {
			return nil, fmt.Errorf("empty hex input in %s", path)
		}
		buf, err := hex.DecodeString(clean)
		if err != nil {
			return nil, fmt.Errorf("hex decode %s: %w", path, err)
		}
		return buf, nil
	case formatRaw:
		return data, nil
	default:
		return nil, fmt.Errorf("unknown format %q", format)
	}
}

func stripWhitespace(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if !isWhitespace(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func isWhitespace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}
