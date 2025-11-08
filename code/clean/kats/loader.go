package kats

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// DefaultRoot is the default directory containing ML-DSA KAT resources.
	DefaultRoot = "code/clean/testdata/kats/ml-dsa"

	// DefaultKeyGenVectors is the default filename for ML-DSA key generation vectors.
	DefaultKeyGenVectors = "ML-DSA-keyGen-FIPS204-internalProjection.json"

	// DefaultSigGenVectors is the default filename for ML-DSA signature generation vectors.
	DefaultSigGenVectors = "ML-DSA-sigGen-FIPS204-internalProjection.json"

	// DefaultSigVerVectors is the default filename for ML-DSA signature verification vectors.
	DefaultSigVerVectors = "ML-DSA-sigVer-FIPS204-internalProjection.json"
)

// KeyGenVectors captures the structure of ACVP ML-DSA key generation vectors.
type KeyGenVectors struct {
	Algorithm  string            `json:"algorithm"`
	Mode       string            `json:"mode"`
	Revision   string            `json:"revision"`
	IsSample   bool              `json:"isSample"`
	TestGroups []KeyGenTestGroup `json:"testGroups"`
}

// KeyGenTestGroup contains a set of key generation test cases for a parameter set.
type KeyGenTestGroup struct {
	TargetGroupID int              `json:"tgId"`
	TestType      string           `json:"testType"`
	ParameterSet  string           `json:"parameterSet"`
	Tests         []KeyGenTestCase `json:"tests"`
}

// KeyGenTestCase represents a single key generation test case.
type KeyGenTestCase struct {
	CaseID   int    `json:"tcId"`
	Deferred bool   `json:"deferred"`
	Seed     string `json:"seed"`
	Public   string `json:"pk"`
	Secret   string `json:"sk"`
}

// SigGenVectors captures the structure of ACVP ML-DSA signature generation vectors.
type SigGenVectors struct {
	Algorithm  string            `json:"algorithm"`
	Mode       string            `json:"mode"`
	Revision   string            `json:"revision"`
	IsSample   bool              `json:"isSample"`
	TestGroups []SigGenTestGroup `json:"testGroups"`
}

// SigGenTestGroup contains signature generation tests for a parameter set.
type SigGenTestGroup struct {
	TargetGroupID      int              `json:"tgId"`
	TestType           string           `json:"testType"`
	ParameterSet       string           `json:"parameterSet"`
	Deterministic      bool             `json:"deterministic"`
	SignatureInterface string           `json:"signatureInterface"`
	PreHash            string           `json:"preHash"`
	ExternalMu         bool             `json:"externalMu"`
	CornerCase         string           `json:"cornerCase"`
	Tests              []SigGenTestCase `json:"tests"`
}

// SigGenTestCase represents a single signature generation test case.
type SigGenTestCase struct {
	CaseID    int    `json:"tcId"`
	Deferred  bool   `json:"deferred"`
	Message   string `json:"message"`
	Public    string `json:"pk"`
	Secret    string `json:"sk"`
	Context   string `json:"context"`
	HashAlg   string `json:"hashAlg"`
	Signature string `json:"signature"`
}

// SigVerVectors captures the structure of ACVP ML-DSA signature verification vectors.
type SigVerVectors struct {
	Algorithm  string            `json:"algorithm"`
	Mode       string            `json:"mode"`
	Revision   string            `json:"revision"`
	IsSample   bool              `json:"isSample"`
	TestGroups []SigVerTestGroup `json:"testGroups"`
}

// SigVerTestGroup contains signature verification tests for a parameter set.
type SigVerTestGroup struct {
	TargetGroupID      int              `json:"tgId"`
	TestType           string           `json:"testType"`
	ParameterSet       string           `json:"parameterSet"`
	SignatureInterface string           `json:"signatureInterface"`
	PreHash            string           `json:"preHash"`
	ExternalMu         bool             `json:"externalMu"`
	Tests              []SigVerTestCase `json:"tests"`
}

// SigVerTestCase represents a single signature verification test vector.
type SigVerTestCase struct {
	CaseID     int    `json:"tcId"`
	TestPassed bool   `json:"testPassed"`
	Deferred   bool   `json:"deferred"`
	Public     string `json:"pk"`
	Secret     string `json:"sk"`
	Message    string `json:"message"`
	Context    string `json:"context"`
	HashAlg    string `json:"hashAlg"`
	Signature  string `json:"signature"`
	Reason     string `json:"reason"`
}

// LoadKeyGenVectors loads key generation KAT vectors from disk.
func LoadKeyGenVectors(path string) (*KeyGenVectors, error) {
	if path == "" {
		path = filepath.Join(DefaultRoot, DefaultKeyGenVectors)
	}
	resolved, err := resolvePath(path)
	if err != nil {
		return nil, err
	}
	var vectors KeyGenVectors
	if err := decodeJSON(resolved, &vectors); err != nil {
		return nil, err
	}
	if err := expectMode(vectors.Mode, "keyGen"); err != nil {
		return nil, err
	}
	return &vectors, nil
}

// LoadSigGenVectors loads signature generation KAT vectors from disk.
func LoadSigGenVectors(path string) (*SigGenVectors, error) {
	if path == "" {
		path = filepath.Join(DefaultRoot, DefaultSigGenVectors)
	}
	resolved, err := resolvePath(path)
	if err != nil {
		return nil, err
	}
	var vectors SigGenVectors
	if err := decodeJSON(resolved, &vectors); err != nil {
		return nil, err
	}
	if err := expectMode(vectors.Mode, "sigGen"); err != nil {
		return nil, err
	}
	return &vectors, nil
}

// LoadSigVerVectors loads signature verification KAT vectors from disk.
func LoadSigVerVectors(path string) (*SigVerVectors, error) {
	if path == "" {
		path = filepath.Join(DefaultRoot, DefaultSigVerVectors)
	}
	resolved, err := resolvePath(path)
	if err != nil {
		return nil, err
	}
	var vectors SigVerVectors
	if err := decodeJSON(resolved, &vectors); err != nil {
		return nil, err
	}
	if err := expectMode(vectors.Mode, "sigVer"); err != nil {
		return nil, err
	}
	return &vectors, nil
}

func decodeJSON(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("kats: open %q: %w", path, err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(v); err != nil {
		return fmt.Errorf("kats: decode %q: %w", path, err)
	}
	return nil
}

func resolvePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	root, err := moduleRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, path), nil
}

func moduleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("kats: getwd: %w", err)
	}
	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("kats: unable to locate go.mod starting from %q", dir)
		}
		dir = parent
	}
}

func expectMode(actual, expected string) error {
	if actual != expected {
		return fmt.Errorf("kats: expected mode %q, got %q", expected, actual)
	}
	return nil
}
