package style

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

// TypesChecker checks BARR-C type safety rules
type TypesChecker struct {
	parser *Parser
}

// NewTypesChecker creates a new type safety checker
func NewTypesChecker() *TypesChecker {
	return &TypesChecker{
		parser: NewParser(),
	}
}

// Type mappings: primitive types to stdint.h equivalents
var typeReplacements = map[string]string{
	"unsigned char":  "uint8_t",
	"signed char":    "int8_t",
	"unsigned short": "uint16_t",
	"signed short":   "int16_t",
	"short":          "int16_t",
	"unsigned int":   "uint32_t",
	"signed int":     "int32_t",
	"unsigned long":  "uint32_t or uint64_t (depends on platform)",
	"signed long":    "int32_t or int64_t (depends on platform)",
	"long":           "int32_t or int64_t (depends on platform)",
}

// Check implements the Checker interface
func (t *TypesChecker) Check(tree *sitter.Tree, source []byte) []Violation {
	var violations []Violation

	// Check all declarations for non-stdint types
	declarations := t.parser.FindNodesByType(tree, NodeDeclaration)
	for _, declNode := range declarations {
		typeSpec := t.parser.GetTypeSpecifier(declNode, source)
		if typeSpec == "" {
			continue
		}

		// Check if this type should be replaced with stdint.h type
		if replacement, shouldReplace := typeReplacements[typeSpec]; shouldReplace {
			varName := t.parser.GetVariableName(declNode, source)
			violations = append(violations, Violation{
				Rule:   StdintTypesRule,
				File:   "",
				Line:   int(declNode.StartPoint().Row) + 1,
				Column: int(declNode.StartPoint().Column) + 1,
				Message: fmt.Sprintf(
					"Type '%s' should be replaced with stdint.h type (use %s instead)",
					typeSpec,
					replacement,
				),
				Suggestion: t.suggestStdintReplacement(typeSpec, varName),
			})
		}

		// Check for const correctness on pointers in declarations
		t.checkConstPointers(declNode, source, &violations)
	}

	// Check function parameters for const correctness
	functions := t.parser.FindNodesByType(tree, NodeFunction)
	for _, funcNode := range functions {
		t.checkFunctionParameterConst(funcNode, source, &violations)
	}

	return violations
}

// checkConstPointers checks if pointer declarations have proper const qualifiers
func (t *TypesChecker) checkConstPointers(declNode *sitter.Node, source []byte, violations *[]Violation) {
	// Look for pointer declarators
	var checkNode func(*sitter.Node)
	checkNode = func(node *sitter.Node) {
		if node.Type() == "pointer_declarator" {
			// Check if const qualifier is present
			if !t.parser.IsConstPointer(node) {
				// This is a warning, not an error - const is recommended but not always required
				varName := t.extractPointerName(node, source)
				*violations = append(*violations, Violation{
					Rule:   ConstPointersRule,
					File:   "",
					Line:   int(node.StartPoint().Row) + 1,
					Column: int(node.StartPoint().Column) + 1,
					Message: fmt.Sprintf(
						"Pointer '%s' should consider const qualifier for safety (const uint8_t* instead of uint8_t*)",
						varName,
					),
					Suggestion: fmt.Sprintf("Add const qualifier: const ... *%s", varName),
				})
			}
		}

		// Recurse into children
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child != nil {
				checkNode(child)
			}
		}
	}

	checkNode(declNode)
}

// extractPointerName extracts the identifier from a pointer declarator
func (t *TypesChecker) extractPointerName(ptrNode *sitter.Node, source []byte) string {
	for i := 0; i < int(ptrNode.ChildCount()); i++ {
		child := ptrNode.Child(i)
		if child.Type() == "identifier" {
			return child.Content(source)
		}
		// Recurse for nested pointers
		if child.Type() == "pointer_declarator" {
			return t.extractPointerName(child, source)
		}
	}
	return "unknown"
}

// checkFunctionParameterConst checks function parameters for const correctness on pointers
func (t *TypesChecker) checkFunctionParameterConst(funcNode *sitter.Node, source []byte, violations *[]Violation) {
	// Find parameter_list within function
	for i := 0; i < int(funcNode.ChildCount()); i++ {
		child := funcNode.Child(i)
		if child.Type() == "function_declarator" {
			// Look for parameter_list
			for j := 0; j < int(child.ChildCount()); j++ {
				paramList := child.Child(j)
				if paramList.Type() == "parameter_list" {
					// Check each parameter
					for k := 0; k < int(paramList.ChildCount()); k++ {
						param := paramList.Child(k)
						if param.Type() == "parameter_declaration" {
							t.checkParameterPointerConst(param, source, violations)
						}
					}
				}
			}
		}
	}
}

// checkParameterPointerConst checks if a parameter has proper const for pointers
func (t *TypesChecker) checkParameterPointerConst(paramNode *sitter.Node, source []byte, violations *[]Violation) {
	// Look for pointer_declarator in parameter
	hasPointer := false
	hasConst := false
	var ptrNode *sitter.Node

	for i := 0; i < int(paramNode.ChildCount()); i++ {
		child := paramNode.Child(i)
		if child.Type() == "type_qualifier" && child.Content(source) == "const" {
			hasConst = true
		}
		if child.Type() == "pointer_declarator" {
			hasPointer = true
			ptrNode = child
		}
	}

	// If we have a pointer without const, warn
	if hasPointer && !hasConst && ptrNode != nil {
		paramName := t.extractPointerName(ptrNode, source)
		*violations = append(*violations, Violation{
			Rule:   ConstPointersRule,
			File:   "",
			Line:   int(paramNode.StartPoint().Row) + 1,
			Column: int(paramNode.StartPoint().Column) + 1,
			Message: fmt.Sprintf(
				"Pointer parameter '%s' should consider const qualifier for safety (const uint8_t* instead of uint8_t*)",
				paramName,
			),
			Suggestion: fmt.Sprintf("Add const qualifier: const ... *%s", paramName),
		})
	}
}

// suggestStdintReplacement provides a suggestion for replacing a type
func (t *TypesChecker) suggestStdintReplacement(oldType string, varName string) string {
	replacement, ok := typeReplacements[oldType]
	if !ok {
		return ""
	}

	// Handle platform-dependent types
	if strings.Contains(replacement, "or") {
		return fmt.Sprintf("Replace '%s %s' with appropriate stdint.h type (%s)", oldType, varName, replacement)
	}

	return fmt.Sprintf("Replace '%s %s' with '%s %s'", oldType, varName, replacement, varName)
}
