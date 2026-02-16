//go:build linux && cgo

package style

import (
	"fmt"
	"regexp"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

// NamingChecker checks BARR-C naming conventions
type NamingChecker struct {
	parser *Parser
}

// NewNamingChecker creates a new naming convention checker
func NewNamingChecker() *NamingChecker {
	return &NamingChecker{
		parser: NewParser(),
	}
}

// BARR-C naming patterns:
// - Functions: Module_Function() format (PascalCase_PascalCase)
// - Variables: snake_case
// - Constants: UPPER_SNAKE
var (
	// Function naming: Module_Function (e.g., LED_Init, UART_SendByte)
	// At least one underscore, starts with uppercase, no consecutive underscores
	funcNamePattern = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*_[A-Z][a-zA-Z0-9]*(_[A-Z][a-zA-Z0-9]*)*$`)

	// Variable naming: snake_case (e.g., status_flag, uart_buffer)
	// Starts with lowercase, allows underscores
	varNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

	// Constant naming: UPPER_SNAKE (e.g., MAX_BUFFER_SIZE, LED_PIN)
	// All uppercase with underscores
	constNamePattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
)

// Check implements the Checker interface
func (n *NamingChecker) Check(tree *sitter.Tree, source []byte) []Violation {
	var violations []Violation

	// Check function names
	functions := n.parser.FindNodesByType(tree, NodeFunction)
	for _, funcNode := range functions {
		funcName := n.parser.GetFunctionName(funcNode, source)
		if funcName == "" {
			continue
		}

		// Skip main() - it's a special case required by C standard
		if funcName == "main" {
			continue
		}

		if !funcNamePattern.MatchString(funcName) {
			violations = append(violations, Violation{
				Rule:       FunctionNamingRule,
				File:       "",
				Line:       int(funcNode.StartPoint().Row) + 1,
				Column:     int(funcNode.StartPoint().Column) + 1,
				Message:    fmt.Sprintf("Function '%s' should use Module_Function format (got: %s)", funcName, funcName),
				Suggestion: n.suggestFunctionName(funcName),
			})
		}
	}

	// Check variable names
	declarations := n.parser.FindNodesByType(tree, NodeDeclaration)
	for _, declNode := range declarations {
		varName := n.parser.GetVariableName(declNode, source)
		if varName == "" {
			continue
		}

		// Determine if this is a macro-like constant based on context
		// Only global const variables at file scope should use UPPER_SNAKE
		// Local const variables should use snake_case
		isGlobalConst := n.isGlobalConstDeclaration(declNode, source)

		if isGlobalConst {
			// Check constant naming (UPPER_SNAKE) - only for global consts
			if !constNamePattern.MatchString(varName) {
				violations = append(violations, Violation{
					Rule:       ConstantNamingRule,
					File:       "",
					Line:       int(declNode.StartPoint().Row) + 1,
					Column:     int(declNode.StartPoint().Column) + 1,
					Message:    fmt.Sprintf("Global constant '%s' should use UPPER_SNAKE format (got: %s)", varName, varName),
					Suggestion: n.suggestConstantName(varName),
				})
			}
		} else {
			// Check variable naming (snake_case) - for all non-global-const variables
			if !varNamePattern.MatchString(varName) {
				violations = append(violations, Violation{
					Rule:       VariableNamingRule,
					File:       "",
					Line:       int(declNode.StartPoint().Row) + 1,
					Column:     int(declNode.StartPoint().Column) + 1,
					Message:    fmt.Sprintf("Variable '%s' should use snake_case format (got: %s)", varName, varName),
					Suggestion: n.suggestVariableName(varName),
				})
			}
		}
	}

	return violations
}

// isGlobalConstDeclaration checks if a declaration is a global const (should use UPPER_SNAKE)
// Local const variables should use snake_case, only global consts use UPPER_SNAKE
func (n *NamingChecker) isGlobalConstDeclaration(declNode *sitter.Node, source []byte) bool {
	// Check if has const qualifier
	hasConst := false
	for i := 0; i < int(declNode.ChildCount()); i++ {
		child := declNode.Child(i)
		if child.Type() == "type_qualifier" {
			content := child.Content(source)
			if content == "const" {
				hasConst = true
				break
			}
		}
	}

	if !hasConst {
		return false
	}

	// Check if this is at global scope (parent is translation_unit)
	parent := declNode.Parent()
	if parent != nil && parent.Type() == "translation_unit" {
		return true
	}

	return false
}

// suggestFunctionName converts a name to Module_Function format
func (n *NamingChecker) suggestFunctionName(name string) string {
	// If already has underscore, capitalize appropriately
	if strings.Contains(name, "_") {
		parts := strings.Split(name, "_")
		for i := range parts {
			if len(parts[i]) > 0 {
				parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
			}
		}
		return strings.Join(parts, "_")
	}

	// No underscore - suggest adding a module prefix
	capitalized := strings.ToUpper(name[:1]) + name[1:]
	return fmt.Sprintf("Module_%s", capitalized)
}

// suggestVariableName converts a name to snake_case
func (n *NamingChecker) suggestVariableName(name string) string {
	// Convert camelCase or PascalCase to snake_case
	var result strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// suggestConstantName converts a name to UPPER_SNAKE
func (n *NamingChecker) suggestConstantName(name string) string {
	// Convert to snake_case first, then uppercase
	snakeCase := n.suggestVariableName(name)
	return strings.ToUpper(snakeCase)
}
