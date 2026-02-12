//go:build linux

package style

import (
	"os"
)

// StyleChecker orchestrates multiple BARR-C checkers
type StyleChecker struct {
	checkers []Checker
	parser   *Parser
}

// NewStyleChecker creates a new style checker with all BARR-C checkers enabled
func NewStyleChecker() *StyleChecker {
	return &StyleChecker{
		checkers: []Checker{
			NewNamingChecker(),
			NewTypesChecker(),
		},
		parser: NewParser(),
	}
}

// CheckFile checks a C file for BARR-C violations
func (s *StyleChecker) CheckFile(filepath string) ([]Violation, error) {
	// Read source file
	source, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// Parse the file
	tree, err := s.parser.Parse(source)
	if err != nil {
		return nil, err
	}

	// Run all checkers
	var allViolations []Violation
	for _, checker := range s.checkers {
		violations := checker.Check(tree, source)
		// Set file path on all violations
		for i := range violations {
			violations[i].File = filepath
		}
		allViolations = append(allViolations, violations...)
	}

	return allViolations, nil
}

// CheckSource checks C source code (as bytes) for BARR-C violations
func (s *StyleChecker) CheckSource(source []byte, filename string) ([]Violation, error) {
	// Parse the source
	tree, err := s.parser.Parse(source)
	if err != nil {
		return nil, err
	}

	// Run all checkers
	var allViolations []Violation
	for _, checker := range s.checkers {
		violations := checker.Check(tree, source)
		// Set file path on all violations
		for i := range violations {
			violations[i].File = filename
		}
		allViolations = append(allViolations, violations...)
	}

	return allViolations, nil
}

// AddChecker adds a custom checker to the style checker
func (s *StyleChecker) AddChecker(checker Checker) {
	s.checkers = append(s.checkers, checker)
}
