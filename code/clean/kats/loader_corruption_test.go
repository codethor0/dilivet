// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package kats

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestLoadSigVerVectors_CorruptedJSON tests that corrupted JSON files are handled gracefully.
func TestLoadSigVerVectors_CorruptedJSON(t *testing.T) {
	tDir := t.TempDir()

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "missing sig field",
			content: `{"algorithm":"ML-DSA","mode":"sigVer","testGroups":[{"tests":[{"tcId":1,"pk":"00","message":"00"}]}]}`,
			wantErr: false, // Loader doesn't validate required fields at parse time
		},
		{
			name:    "missing pk field",
			content: `{"algorithm":"ML-DSA","mode":"sigVer","testGroups":[{"tests":[{"tcId":1,"signature":"00","message":"00"}]}]}`,
			wantErr: false, // Loader doesn't validate required fields at parse time
		},
		{
			name:    "missing message field",
			content: `{"algorithm":"ML-DSA","mode":"sigVer","testGroups":[{"tests":[{"tcId":1,"pk":"00","signature":"00"}]}]}`,
			wantErr: false, // Loader doesn't validate required fields at parse time
		},
		{
			name:    "wrong type for pk (string instead of hex)",
			content: `{"algorithm":"ML-DSA","mode":"sigVer","testGroups":[{"tests":[{"tcId":1,"pk":123,"signature":"00","message":"00"}]}]}`,
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			content: `{invalid json}`,
			wantErr: true,
		},
		{
			name:    "empty JSON object",
			content: `{}`,
			wantErr: true,
		},
		{
			name:    "missing testGroups",
			content: `{"algorithm":"ML-DSA","mode":"sigVer"}`,
			wantErr: false, // Loader allows missing testGroups (empty slice)
		},
		{
			name:    "empty testGroups",
			content: `{"algorithm":"ML-DSA","mode":"sigVer","testGroups":[]}`,
			wantErr: false, // Empty but valid structure
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tDir, "corrupted.json")
			if err := os.WriteFile(path, []byte(tt.content), 0o600); err != nil {
				t.Fatalf("write test file: %v", err)
			}

			_, err := LoadSigVerVectors(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadSigVerVectors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestLoadSigVerVectors_InvalidHex tests that invalid hex strings are handled.
func TestLoadSigVerVectors_InvalidHex(t *testing.T) {
	tDir := t.TempDir()

	// Create a vector with invalid hex
	payload := map[string]interface{}{
		"algorithm": "ML-DSA",
		"mode":      "sigVer",
		"testGroups": []map[string]interface{}{
			{
				"tests": []map[string]interface{}{
					{
						"tcId":       1,
						"pk":         "zzzz", // Invalid hex
						"signature":  "00",
						"message":    "00",
						"testPassed": true,
					},
				},
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	path := filepath.Join(tDir, "invalid-hex.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	vectors, err := LoadSigVerVectors(path)
	if err != nil {
		t.Fatalf("LoadSigVerVectors should not fail on invalid hex in JSON: %v", err)
	}

	// The loader should succeed, but hex decoding will fail when used
	if len(vectors.TestGroups) == 0 || len(vectors.TestGroups[0].Tests) == 0 {
		t.Fatal("Expected at least one test case")
	}

	// The hex decode will fail when the CLI tries to use it, which is acceptable
	// The loader's job is just to parse the JSON structure
}

// TestLoadSigVerVectors_NonExistentFile tests handling of missing files.
func TestLoadSigVerVectors_NonExistentFile(t *testing.T) {
	_, err := LoadSigVerVectors("/nonexistent/path/to/file.json")
	if err == nil {
		t.Error("LoadSigVerVectors should fail on non-existent file")
	}
}

// TestLoadSigVerVectors_EmptyFile tests handling of empty files.
func TestLoadSigVerVectors_EmptyFile(t *testing.T) {
	tDir := t.TempDir()
	path := filepath.Join(tDir, "empty.json")
	if err := os.WriteFile(path, []byte{}, 0o600); err != nil {
		t.Fatalf("create empty file: %v", err)
	}

	_, err := LoadSigVerVectors(path)
	if err == nil {
		t.Error("LoadSigVerVectors should fail on empty file")
	}
}

// TestLoadSigVerVectors_MalformedStructure tests various malformed JSON structures.
func TestLoadSigVerVectors_MalformedStructure(t *testing.T) {
	tDir := t.TempDir()

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "testGroups is not an array",
			content: `{"algorithm":"ML-DSA","mode":"sigVer","testGroups":"not-an-array"}`,
			wantErr: true,
		},
		{
			name:    "tests is not an array",
			content: `{"algorithm":"ML-DSA","mode":"sigVer","testGroups":[{"tests":"not-an-array"}]}`,
			wantErr: true,
		},
		{
			name:    "test case is not an object",
			content: `{"algorithm":"ML-DSA","mode":"sigVer","testGroups":[{"tests":["not-an-object"]}]}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tDir, "malformed.json")
			if err := os.WriteFile(path, []byte(tt.content), 0o600); err != nil {
				t.Fatalf("write: %v", err)
			}

			_, err := LoadSigVerVectors(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadSigVerVectors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

