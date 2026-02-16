//go:build linux && cgo

package style

import (
	"testing"
)

func TestParserBasic(t *testing.T) {
	parser := NewParser()

	// Simple C program
	source := []byte(`int main() { return 0; }`)

	tree, err := parser.Parse(source)
	if err != nil {
		t.Fatalf("Failed to parse simple C code: %v", err)
	}

	if tree == nil {
		t.Fatal("Tree is nil")
	}

	root := tree.RootNode()
	if root == nil {
		t.Fatal("Root node is nil")
	}

	if root.Type() != "translation_unit" {
		t.Errorf("Expected root type 'translation_unit', got %s", root.Type())
	}
}

func TestParserFindFunctions(t *testing.T) {
	parser := NewParser()

	source := []byte(`
void setup() {
	return;
}

int main() {
	return 0;
}
`)

	tree, err := parser.Parse(source)
	if err != nil {
		t.Fatalf("Failed to parse C code: %v", err)
	}

	functions := parser.FindNodesByType(tree, NodeFunction)
	if len(functions) != 2 {
		t.Errorf("Expected 2 functions, found %d", len(functions))
	}

	// Test function name extraction
	name1 := parser.GetFunctionName(functions[0], source)
	if name1 != "setup" {
		t.Errorf("Expected function name 'setup', got '%s'", name1)
	}

	name2 := parser.GetFunctionName(functions[1], source)
	if name2 != "main" {
		t.Errorf("Expected function name 'main', got '%s'", name2)
	}
}

func TestParserFindVariables(t *testing.T) {
	parser := NewParser()

	source := []byte(`
int global_var = 42;
unsigned char status;
`)

	tree, err := parser.Parse(source)
	if err != nil {
		t.Fatalf("Failed to parse C code: %v", err)
	}

	decls := parser.FindNodesByType(tree, NodeDeclaration)
	if len(decls) < 1 {
		t.Fatalf("Expected at least 1 declaration, found %d", len(decls))
	}

	// Test type extraction
	type1 := parser.GetTypeSpecifier(decls[0], source)
	if type1 != "int" {
		t.Errorf("Expected type 'int', got '%s'", type1)
	}
}
