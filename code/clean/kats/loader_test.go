// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package kats

import (
	"path/filepath"
	"testing"
)

func TestLoadKeyGenVectors(t *testing.T) {
	vectors, err := LoadKeyGenVectors("")
	if err != nil {
		t.Fatalf("LoadKeyGenVectors: %v", err)
	}
	if got := vectors.Mode; got != "keyGen" {
		t.Fatalf("Mode = %q, want %q", got, "keyGen")
	}
	if len(vectors.TestGroups) == 0 {
		t.Fatalf("no test groups decoded")
	}
	for _, tg := range vectors.TestGroups {
		if tg.ParameterSet == "" {
			t.Fatalf("test group %d missing parameter set", tg.TargetGroupID)
		}
		if len(tg.Tests) == 0 {
			t.Fatalf("test group %d has no tests", tg.TargetGroupID)
		}
		for _, tc := range tg.Tests {
			if tc.Seed == "" || tc.Public == "" || tc.Secret == "" {
				t.Fatalf("test case %d missing data", tc.CaseID)
			}
		}
	}
}

func TestLoadSigGenVectors(t *testing.T) {
	vectors, err := LoadSigGenVectors("")
	if err != nil {
		t.Fatalf("LoadSigGenVectors: %v", err)
	}
	if got := vectors.Mode; got != "sigGen" {
		t.Fatalf("Mode = %q, want %q", got, "sigGen")
	}
	if len(vectors.TestGroups) == 0 {
		t.Fatalf("no test groups decoded")
	}
	for _, tg := range vectors.TestGroups {
		if tg.ParameterSet == "" {
			t.Fatalf("test group %d missing parameter set", tg.TargetGroupID)
		}
		if len(tg.Tests) == 0 {
			t.Fatalf("test group %d has no tests", tg.TargetGroupID)
		}
		for _, tc := range tg.Tests {
			if tc.Public == "" || tc.Secret == "" || tc.Signature == "" {
				t.Fatalf("test case %d missing key material or signature", tc.CaseID)
			}
		}
	}
}

func TestLoadSigVerVectors(t *testing.T) {
	vectors, err := LoadSigVerVectors("")
	if err != nil {
		t.Fatalf("LoadSigVerVectors: %v", err)
	}
	if got := vectors.Mode; got != "sigVer" {
		t.Fatalf("Mode = %q, want %q", got, "sigVer")
	}
	if len(vectors.TestGroups) == 0 {
		t.Fatalf("no test groups decoded")
	}
	for _, tg := range vectors.TestGroups {
		if tg.ParameterSet == "" {
			t.Fatalf("test group %d missing parameter set", tg.TargetGroupID)
		}
		if len(tg.Tests) == 0 {
			t.Fatalf("test group %d has no tests", tg.TargetGroupID)
		}
		for _, tc := range tg.Tests {
			if tc.Public == "" || tc.Signature == "" {
				t.Fatalf("test case %d missing public key or signature", tc.CaseID)
			}
		}
	}
}

func TestLoadCustomPath(t *testing.T) {
	path := filepath.Join(DefaultRoot, DefaultSigVerVectors)
	if _, err := LoadSigVerVectors(path); err != nil {
		t.Fatalf("LoadSigVerVectors(%q): %v", path, err)
	}
}

func TestVerifyVectorsTODO(t *testing.T) {
	t.Skip("TODO: integrate ML-DSA verification with KAT vectors")
}
