//go:build !linux

package style

// StyleFixer stub for non-Linux platforms
type StyleFixer struct{}

// NewStyleFixer creates a stub style fixer
func NewStyleFixer() *StyleFixer {
	return &StyleFixer{}
}

// ApplyFixes returns the source unchanged on non-Linux platforms
func (s *StyleFixer) ApplyFixes(violations []Violation, source []byte) ([]byte, []string) {
	// No-op: return source unchanged
	return source, nil
}
