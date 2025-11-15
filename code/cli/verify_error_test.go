// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVerify_MissingFile(t *testing.T) {
	var errOut bytes.Buffer
	app := &App{
		Name: "dilivet",
		Err:  &errOut,
	}

	exitCode := app.Run([]string{
		"verify",
		"-pub", "nope.hex",
		"-sig", "nope.sig",
		"-msg", "nope.msg",
	})

	if exitCode == 0 {
		t.Error("Expected non-zero exit code for missing file")
	}
	if !strings.Contains(errOut.String(), "read") {
		t.Errorf("Expected error message about reading file, got: %q", errOut.String())
	}
}

func TestVerify_EmptyFile(t *testing.T) {
	tDir := t.TempDir()

	emptyPath := filepath.Join(tDir, "empty.hex")
	if err := os.WriteFile(emptyPath, []byte{}, 0o600); err != nil {
		t.Fatalf("create empty file: %v", err)
	}

	var errOut bytes.Buffer
	app := &App{
		Name: "dilivet",
		Err:  &errOut,
	}

	exitCode := app.Run([]string{
		"verify",
		"-pub", emptyPath,
		"-sig", emptyPath,
		"-msg", emptyPath,
	})

	if exitCode == 0 {
		t.Error("Expected non-zero exit code for empty hex file")
	}
	// Should get error about empty hex input or hex decode failure
	if !strings.Contains(errOut.String(), "empty") && !strings.Contains(errOut.String(), "hex decode") {
		t.Errorf("Expected error about empty/hex decode, got: %q", errOut.String())
	}
}

func TestVerify_NonHexContent(t *testing.T) {
	tDir := t.TempDir()

	badPath := filepath.Join(tDir, "bad.hex")
	if err := os.WriteFile(badPath, []byte("zzzz"), 0o600); err != nil {
		t.Fatalf("create bad hex file: %v", err)
	}

	var errOut bytes.Buffer
	app := &App{
		Name: "dilivet",
		Err:  &errOut,
	}

	exitCode := app.Run([]string{
		"verify",
		"-pub", badPath,
		"-sig", badPath,
		"-msg", badPath,
	})

	if exitCode == 0 {
		t.Error("Expected non-zero exit code for invalid hex")
	}
	if !strings.Contains(errOut.String(), "hex decode") {
		t.Errorf("Expected hex decode error, got: %q", errOut.String())
	}
}

func TestVerify_InvalidFormat(t *testing.T) {
	tDir := t.TempDir()

	testPath := filepath.Join(tDir, "test.bin")
	if err := os.WriteFile(testPath, []byte("test"), 0o600); err != nil {
		t.Fatalf("create test file: %v", err)
	}

	var errOut bytes.Buffer
	app := &App{
		Name: "dilivet",
		Err:  &errOut,
	}

	exitCode := app.Run([]string{
		"verify",
		"-pub", testPath,
		"-sig", testPath,
		"-msg", testPath,
		"-pub-format", "invalid",
	})

	if exitCode == 0 {
		t.Error("Expected non-zero exit code for invalid format")
	}
	if !strings.Contains(errOut.String(), "unknown format") {
		t.Errorf("Expected format error, got: %q", errOut.String())
	}
}

func TestVerify_HexWithWhitespace(t *testing.T) {
	tDir := t.TempDir()

	// Hex with spaces, newlines, tabs should be handled
	hexWithWhitespace := "00 01 02\n03\t04\r05"
	pubPath := filepath.Join(tDir, "pk.hex")
	if err := os.WriteFile(pubPath, []byte(hexWithWhitespace), 0o600); err != nil {
		t.Fatalf("create hex file: %v", err)
	}

	sigPath := filepath.Join(tDir, "sig.hex")
	if err := os.WriteFile(sigPath, []byte("0000000000000000"), 0o600); err != nil {
		t.Fatalf("create sig file: %v", err)
	}

	msgPath := filepath.Join(tDir, "msg.bin")
	if err := os.WriteFile(msgPath, []byte("test"), 0o600); err != nil {
		t.Fatalf("create msg file: %v", err)
	}

	var errOut bytes.Buffer
	app := &App{
		Name: "dilivet",
		Err:  &errOut,
	}

	// Should not panic, should handle whitespace gracefully
	exitCode := app.Run([]string{
		"verify",
		"-pub", pubPath,
		"-sig", sigPath,
		"-msg", msgPath,
	})

	// May fail verification, but should not panic or have decode errors
	if exitCode == 0 {
		t.Log("Verification succeeded (unexpected but not an error)")
	}
	// Check that we didn't get a hex decode error
	if strings.Contains(errOut.String(), "hex decode") {
		t.Errorf("Whitespace should be stripped, got error: %q", errOut.String())
	}
}

func TestVerify_CRLFLineEndings(t *testing.T) {
	tDir := t.TempDir()

	// Hex file with CRLF line endings
	hexCRLF := "0001020304050607\r\n08090a0b0c0d0e0f"
	pubPath := filepath.Join(tDir, "pk.hex")
	if err := os.WriteFile(pubPath, []byte(hexCRLF), 0o600); err != nil {
		t.Fatalf("create hex file: %v", err)
	}

	sigPath := filepath.Join(tDir, "sig.hex")
	if err := os.WriteFile(sigPath, []byte("0000000000000000"), 0o600); err != nil {
		t.Fatalf("create sig file: %v", err)
	}

	msgPath := filepath.Join(tDir, "msg.bin")
	if err := os.WriteFile(msgPath, []byte("test"), 0o600); err != nil {
		t.Fatalf("create msg file: %v", err)
	}

	var errOut bytes.Buffer
	app := &App{
		Name: "dilivet",
		Err:  &errOut,
	}

	exitCode := app.Run([]string{
		"verify",
		"-pub", pubPath,
		"-sig", sigPath,
		"-msg", msgPath,
	})

	// Should handle CRLF without hex decode errors
	if strings.Contains(errOut.String(), "hex decode") {
		t.Errorf("CRLF should be handled, got error: %q", errOut.String())
	}
	_ = exitCode // May fail verification, that's OK
}

func TestVerify_UTF8BOM(t *testing.T) {
	tDir := t.TempDir()

	// Hex file with UTF-8 BOM
	hexWithBOM := "\xEF\xBB\xBF0001020304050607"
	pubPath := filepath.Join(tDir, "pk.hex")
	if err := os.WriteFile(pubPath, []byte(hexWithBOM), 0o600); err != nil {
		t.Fatalf("create hex file: %v", err)
	}

	sigPath := filepath.Join(tDir, "sig.hex")
	if err := os.WriteFile(sigPath, []byte("0000000000000000"), 0o600); err != nil {
		t.Fatalf("create sig file: %v", err)
	}

	msgPath := filepath.Join(tDir, "msg.bin")
	if err := os.WriteFile(msgPath, []byte("test"), 0o600); err != nil {
		t.Fatalf("create msg file: %v", err)
	}

	var errOut bytes.Buffer
	app := &App{
		Name: "dilivet",
		Err:  &errOut,
	}

	exitCode := app.Run([]string{
		"verify",
		"-pub", pubPath,
		"-sig", sigPath,
		"-msg", msgPath,
	})

	// BOM might cause hex decode issues, but should be handled gracefully
	_ = exitCode // May fail, that's OK - just checking for panics
}

func TestVerify_LargeInput(t *testing.T) {
	tDir := t.TempDir()

	// Create a 10MB message file
	largeMsg := make([]byte, 10*1024*1024)
	for i := range largeMsg {
		largeMsg[i] = byte(i % 256)
	}

	msgPath := filepath.Join(tDir, "large.msg")
	if err := os.WriteFile(msgPath, largeMsg, 0o600); err != nil {
		t.Fatalf("create large file: %v", err)
	}

	// Create minimal valid-looking hex files
	pubPath := filepath.Join(tDir, "pk.hex")
	pubHex := strings.Repeat("00", 1312) // ML-DSA-44 public key size
	if err := os.WriteFile(pubPath, []byte(pubHex), 0o600); err != nil {
		t.Fatalf("create pub file: %v", err)
	}

	sigPath := filepath.Join(tDir, "sig.hex")
	sigHex := strings.Repeat("00", 2420) // ML-DSA-44 signature size
	if err := os.WriteFile(sigPath, []byte(sigHex), 0o600); err != nil {
		t.Fatalf("create sig file: %v", err)
	}

	var errOut bytes.Buffer
	app := &App{
		Name: "dilivet",
		Err:  &errOut,
	}

	// Should not panic or hang
	exitCode := app.Run([]string{
		"verify",
		"-pub", pubPath,
		"-sig", sigPath,
		"-msg", msgPath,
	})

	// May fail verification, but should complete without panic
	_ = exitCode
	if strings.Contains(errOut.String(), "panic") {
		t.Errorf("Should not panic on large input, got: %q", errOut.String())
	}
}

func TestVerify_RelativePath(t *testing.T) {
	tDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Create files in temp dir
	pubPath := "pk.hex"
	if err := os.WriteFile(pubPath, []byte(strings.Repeat("00", 1312)), 0o600); err != nil {
		t.Fatalf("create pub: %v", err)
	}

	sigPath := "sig.hex"
	if err := os.WriteFile(sigPath, []byte(strings.Repeat("00", 2420)), 0o600); err != nil {
		t.Fatalf("create sig: %v", err)
	}

	msgPath := "msg.bin"
	if err := os.WriteFile(msgPath, []byte("test"), 0o600); err != nil {
		t.Fatalf("create msg: %v", err)
	}

	var errOut bytes.Buffer
	app := &App{
		Name: "dilivet",
		Err:  &errOut,
	}

	// Should handle relative paths
	exitCode := app.Run([]string{
		"verify",
		"-pub", pubPath,
		"-sig", sigPath,
		"-msg", msgPath,
	})

	_ = exitCode // May fail verification, that's OK
	if strings.Contains(errOut.String(), "read") && !strings.Contains(errOut.String(), "verification") {
		t.Errorf("Relative path should work, got: %q", errOut.String())
	}
}

func TestVerify_PathTraversal(t *testing.T) {
	var errOut bytes.Buffer
	app := &App{
		Name: "dilivet",
		Err:  &errOut,
	}

	// Try path traversal (should fail safely)
	exitCode := app.Run([]string{
		"verify",
		"-pub", "../../../etc/passwd",
		"-sig", "../../../etc/passwd",
		"-msg", "../../../etc/passwd",
	})

	// Should fail with file read error, not panic
	if exitCode == 0 {
		t.Error("Expected non-zero exit for path traversal attempt")
	}
	if !strings.Contains(errOut.String(), "read") {
		t.Errorf("Expected read error, got: %q", errOut.String())
	}
}

