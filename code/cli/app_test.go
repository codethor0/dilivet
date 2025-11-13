// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package cli

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApp_Version(t *testing.T) {
	var out bytes.Buffer
	app := &App{
		Name:    "test-cli",
		Version: "v1.2.3",
		Out:     &out,
	}

	exitCode := app.Run([]string{"-version"})

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	output := out.String()
	if !strings.Contains(output, "v1.2.3") {
		t.Errorf("Expected version in output, got: %s", output)
	}
}

func TestApp_Help(t *testing.T) {
	var out bytes.Buffer
	app := &App{
		Name:    "test-cli",
		Version: "v1.0.0",
		Out:     &out,
	}

	exitCode := app.Run([]string{"-help"})

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	output := out.String()
	if !strings.Contains(output, "USAGE:") {
		t.Errorf("Expected help message, got: %s", output)
	}
	if !strings.Contains(output, "test-cli") {
		t.Errorf("Expected app name in help, got: %s", output)
	}
}

func TestApp_DefaultBehavior(t *testing.T) {
	var out bytes.Buffer
	app := &App{
		Name:    "dilivet",
		Version: "dev",
		Out:     &out,
	}

	exitCode := app.Run([]string{})

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	output := out.String()
	if !strings.Contains(output, "ML-DSA vetting tool") {
		t.Errorf("Expected default message, got: %s", output)
	}
}

func TestApp_InvalidFlag(t *testing.T) {
	var out, errOut bytes.Buffer
	app := &App{
		Name:    "test-cli",
		Version: "dev",
		Out:     &out,
		Err:     &errOut,
	}

	exitCode := app.Run([]string{"-invalid"})

	if exitCode == 0 {
		t.Error("Expected non-zero exit code for invalid flag")
	}
}

func TestApp_MultipleFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantExit int
		wantOut  string
	}{
		{
			name:     "version takes precedence",
			args:     []string{"-version"},
			wantExit: 0,
			wantOut:  "v1.0.0",
		},
		{
			name:     "help flag",
			args:     []string{"-help"},
			wantExit: 0,
			wantOut:  "USAGE:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			app := &App{
				Name:    "test",
				Version: "v1.0.0",
				Out:     &out,
			}

			exitCode := app.Run(tt.args)

			if exitCode != tt.wantExit {
				t.Errorf("Expected exit code %d, got %d", tt.wantExit, exitCode)
			}

			if !strings.Contains(out.String(), tt.wantOut) {
				t.Errorf("Expected %q in output, got: %s", tt.wantOut, out.String())
			}
		})
	}
}

func TestApp_VerifyCommand(t *testing.T) {
	tDir := t.TempDir()

	pub := bytes.Repeat([]byte{0xAA}, 1312)
	sig := bytes.Repeat([]byte{0xBB}, 2420)
	msg := []byte("hello world")

	pubPath := filepath.Join(tDir, "pk.hex")
	sigPath := filepath.Join(tDir, "sig.hex")
	msgPath := filepath.Join(tDir, "msg.bin")

	if err := os.WriteFile(pubPath, []byte(hex.EncodeToString(pub)), 0o600); err != nil {
		t.Fatalf("write pub: %v", err)
	}
	if err := os.WriteFile(sigPath, []byte(hex.EncodeToString(sig)), 0o600); err != nil {
		t.Fatalf("write sig: %v", err)
	}
	if err := os.WriteFile(msgPath, msg, 0o600); err != nil {
		t.Fatalf("write msg: %v", err)
	}

	var out, errOut bytes.Buffer
	app := &App{
		Name:    "dilivet",
		Version: "dev",
		Out:     &out,
		Err:     &errOut,
	}

	exitCode := app.Run([]string{
		"verify",
		"-pub", pubPath,
		"-sig", sigPath,
		"-msg", msgPath,
	})

	// With dummy data, verification should fail
	if exitCode == 0 {
		t.Fatalf("verify should fail with dummy data, exit = %d", exitCode)
	}
	if !strings.Contains(errOut.String(), "verification failed") {
		t.Fatalf("unexpected stderr: %q", errOut.String())
	}
}

func TestApp_VerifyCommandMissingArgs(t *testing.T) {
	var errOut bytes.Buffer
	app := &App{
		Name:    "dilivet",
		Version: "dev",
		Err:     &errOut,
	}

	exitCode := app.Run([]string{"verify"})

	if exitCode == 0 {
		t.Fatal("expected non-zero exit when required flags missing")
	}
	if !strings.Contains(errOut.String(), "-pub, -sig, and -msg are required") {
		t.Fatalf("unexpected stderr: %q", errOut.String())
	}
}

func TestApp_KATVerifyCommand(t *testing.T) {
	tDir := t.TempDir()

	pk := strings.Repeat("AA", 1312)
	sig := strings.Repeat("BB", 2420)
	msg := "00"

	payload := map[string]interface{}{
		"algorithm": "ML-DSA",
		"mode":      "sigVer",
		"revision":  "FIPS204",
		"isSample":  false,
		"testGroups": []map[string]interface{}{
			{
				"tgId":               1,
				"testType":           "AFT",
				"parameterSet":       "ML-DSA-44",
				"signatureInterface": "external",
				"preHash":            "pure",
				"externalMu":         false,
				"tests": []map[string]interface{}{
					{
						"tcId":       1,
						"testPassed": true,
						"deferred":   false,
						"pk":         pk,
						"sk":         "",
						"message":    msg,
						"context":    "",
						"hashAlg":    "none",
						"signature":  sig,
						"reason":     "",
					},
				},
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal vectors: %v", err)
	}

	vectorPath := filepath.Join(tDir, "vectors.json")
	if err := os.WriteFile(vectorPath, data, 0o600); err != nil {
		t.Fatalf("write vectors: %v", err)
	}

	var out, errOut bytes.Buffer
	app := &App{
		Name:    "dilivet",
		Version: "dev",
		Out:     &out,
		Err:     &errOut,
	}

	exitCode := app.Run([]string{"kat-verify", "-vectors", vectorPath})
	// kat-verify produces a report; exit code 1 indicates test failures (expected with dummy data)
	// Exit code 0 would mean all tests passed
	if exitCode != 1 {
		t.Fatalf("kat-verify should exit 1 with failing tests, got exit = %d, stderr=%q, stdout=%q", exitCode, errOut.String(), out.String())
	}

	// Should produce a report with test results
	if !strings.Contains(out.String(), "Total tests:") && !strings.Contains(out.String(), "TotalTests") {
		t.Fatalf("unexpected stdout: %q", out.String())
	}
}
