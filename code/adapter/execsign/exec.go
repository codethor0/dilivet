// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package execsign

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

// Bin executes an external signer or verifier binary, communicating over stdin/stdout.
type Bin struct {
	Path    string
	Timeout time.Duration
	Env     []string
	Dir     string
}

// Run executes the binary with the supplied input and returns trimmed stdout.
func (b Bin) Run(ctx context.Context, input []byte) ([]byte, error) {
	if b.Path == "" {
		return nil, errors.New("execsign: empty binary path")
	}
	if b.Timeout <= 0 {
		b.Timeout = 5 * time.Second
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, b.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, b.Path)
	cmd.Env = append(cmd.Env, b.Env...)
	if b.Dir != "" {
		cmd.Dir = b.Dir
	}

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = bytes.NewReader(input)

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("execsign: %s timed out after %s", b.Path, b.Timeout)
		}
		return nil, fmt.Errorf("execsign: %w (stderr: %s)", err, stderr.String())
	}

	return bytes.TrimSpace(stdout.Bytes()), nil
}

