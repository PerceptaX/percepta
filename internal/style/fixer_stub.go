//go:build !linux

package style

import (
	"fmt"
)

// Fixer stub for non-Linux platforms
type Fixer struct{}

// NewFixer creates a stub fixer
func NewFixer() *Fixer {
	return &Fixer{}
}

// Fix returns an error indicating style fixing is unavailable
func (f *Fixer) Fix(violations []Violation, source []byte) ([]byte, error) {
	return nil, fmt.Errorf("style fixing requires tree-sitter (Linux only)")
}
