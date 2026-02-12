package style

import (
	"context"
	"errors"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/c"
)

// Parser wraps tree-sitter C parser for analyzing C code
type Parser struct {
	parser *sitter.Parser
}

// NewParser creates a new C parser instance
func NewParser() *Parser {
	parser := sitter.NewParser()
	parser.SetLanguage(c.GetLanguage())
	return &Parser{parser: parser}
}

// Parse parses C source code and returns the AST
func (p *Parser) Parse(source []byte) (*sitter.Tree, error) {
	tree, err := p.parser.ParseCtx(context.TODO(), nil, source)
	if err != nil {
		return nil, err
	}
	if tree == nil {
		return nil, errors.New("failed to parse C code")
	}
	return tree, nil
}

// WalkTree traverses the AST and calls the visitor function for each node
// Uses a cursor-based traversal for efficient tree navigation
func (p *Parser) WalkTree(tree *sitter.Tree, visitor func(*sitter.Node)) {
	cursor := sitter.NewTreeCursor(tree.RootNode())
	defer cursor.Close()

	walkNode := cursor.CurrentNode()
	p.walkHelper(walkNode, visitor)
}

// walkHelper recursively walks the tree and calls visitor on each node
func (p *Parser) walkHelper(node *sitter.Node, visitor func(*sitter.Node)) {
	visitor(node)

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			p.walkHelper(child, visitor)
		}
	}
}

// NodeType represents common C node types we care about for BARR-C checking
type NodeType string

const (
	NodeFunction    NodeType = "function_definition"
	NodeDeclaration NodeType = "declaration"
	NodeIdentifier  NodeType = "identifier"
	NodeNumber      NodeType = "number_literal"
	NodePointer     NodeType = "pointer_declarator"
	NodeParameter   NodeType = "parameter_declaration"
)

// FindNodesByType finds all nodes of a specific type in the tree
// Uses tree-sitter query syntax for efficient pattern matching
//
// Example queries:
//   - Function declarations: "(function_definition) @func"
//   - Variable declarations: "(declaration) @var"
//   - Numeric literals: "(number_literal) @num"
//   - Pointer declarations: "(pointer_declarator) @ptr"
func (p *Parser) FindNodesByType(tree *sitter.Tree, nodeType NodeType) []*sitter.Node {
	var nodes []*sitter.Node

	p.WalkTree(tree, func(node *sitter.Node) {
		if node.Type() == string(nodeType) {
			nodes = append(nodes, node)
		}
	})

	return nodes
}

// GetFunctionName extracts the function name from a function_definition node
// Tree structure: function_definition -> function_declarator -> identifier
func (p *Parser) GetFunctionName(funcNode *sitter.Node, source []byte) string {
	if funcNode.Type() != string(NodeFunction) {
		return ""
	}

	// Find function_declarator child
	for i := 0; i < int(funcNode.ChildCount()); i++ {
		child := funcNode.Child(i)
		if child.Type() == "function_declarator" {
			// Find identifier within declarator
			for j := 0; j < int(child.ChildCount()); j++ {
				grandchild := child.Child(j)
				if grandchild.Type() == string(NodeIdentifier) {
					return grandchild.Content(source)
				}
			}
		}
	}

	return ""
}

// GetVariableName extracts the variable name from a declaration node
// Tree structure: declaration -> declarator -> identifier
func (p *Parser) GetVariableName(declNode *sitter.Node, source []byte) string {
	if declNode.Type() != string(NodeDeclaration) {
		return ""
	}

	// Find declarator child
	for i := 0; i < int(declNode.ChildCount()); i++ {
		child := declNode.Child(i)
		if child.Type() == "init_declarator" || child.Type() == "identifier" {
			// Handle both "int x = 5" and "int x"
			if child.Type() == "identifier" {
				return child.Content(source)
			}
			// For init_declarator, find the identifier
			for j := 0; j < int(child.ChildCount()); j++ {
				grandchild := child.Child(j)
				if grandchild.Type() == string(NodeIdentifier) {
					return grandchild.Content(source)
				}
			}
		}
	}

	return ""
}

// GetTypeSpecifier extracts the type specifier from a declaration node
// Used to detect "unsigned char" vs "uint8_t" etc.
func (p *Parser) GetTypeSpecifier(declNode *sitter.Node, source []byte) string {
	if declNode.Type() != string(NodeDeclaration) {
		return ""
	}

	// Find type specifier (primitive_type, sized_type_specifier, etc.)
	for i := 0; i < int(declNode.ChildCount()); i++ {
		child := declNode.Child(i)
		if child.Type() == "primitive_type" ||
			child.Type() == "sized_type_specifier" ||
			child.Type() == "type_identifier" {
			return child.Content(source)
		}
	}

	return ""
}

// IsConstPointer checks if a pointer declaration includes const qualifier
func (p *Parser) IsConstPointer(ptrNode *sitter.Node) bool {
	// Look for type_qualifier with "const" before the pointer
	parent := ptrNode.Parent()
	if parent == nil {
		return false
	}

	for i := 0; i < int(parent.ChildCount()); i++ {
		child := parent.Child(i)
		if child.Type() == "type_qualifier" {
			// Found a qualifier, could be const
			return true
		}
	}

	return false
}
