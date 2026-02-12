package style

import (
	sitter "github.com/smacker/go-tree-sitter"
)

// Rule represents a BARR-C coding standard rule
type Rule struct {
	ID       string
	Name     string
	Severity string // "error", "warning"
	Category string // "naming", "types", "safety", etc.
}

// Violation represents a detected violation of a BARR-C rule
type Violation struct {
	Rule       Rule
	File       string
	Line       int
	Column     int
	Message    string
	Suggestion string // auto-fix suggestion if available
}

// Checker is the interface that all style checkers must implement
type Checker interface {
	Check(tree *sitter.Tree, source []byte) []Violation
}
