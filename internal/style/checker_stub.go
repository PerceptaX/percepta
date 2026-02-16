//go:build !linux || !cgo

package style

import (
	"fmt"
)

// StyleChecker stub for non-Linux platforms (tree-sitter unavailable)
type StyleChecker struct{}

// NewStyleChecker creates a stub style checker
func NewStyleChecker() *StyleChecker {
	return &StyleChecker{}
}

// CheckFile returns an error indicating style checking is unavailable
func (s *StyleChecker) CheckFile(filepath string) ([]Violation, error) {
	return nil, fmt.Errorf("style checking requires tree-sitter (Linux only)")
}

// CheckSource returns an error indicating style checking is unavailable
func (s *StyleChecker) CheckSource(source []byte, filename string) ([]Violation, error) {
	return nil, fmt.Errorf("style checking requires tree-sitter (Linux only)")
}
