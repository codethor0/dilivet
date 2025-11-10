// © 2025 Thor Thor
// Contact: codethor@gmail.com
// LinkedIn: https://www.linkedin.com/in/thor-thor0
// SPDX-License-Identifier: MIT

package diag

// Report aggregates diagnostic counters during signing/verification.
type Report struct {
	Rejections int
}

// NewReport returns an empty diagnostic report.
func NewReport() (*Report, error) {
	return &Report{}, nil
}
