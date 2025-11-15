// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package diag

// Report aggregates diagnostic counters during signing/verification.
type Report struct {
	TotalTests         int `json:"total_tests"`
	StrictPasses       int `json:"strict_passes"`
	StructuralWarnings int `json:"structural_warnings"`
	StructuralFailures int `json:"structural_failures"`
	DecodeFailures     int `json:"decode_failures"`
}

// NewReport returns an empty diagnostic report.
func NewReport() (*Report, error) {
	return &Report{}, nil
}
