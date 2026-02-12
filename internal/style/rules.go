package style

// BARR-C Embedded C Coding Standard Rules
// Based on: https://barrgroup.com/embedded-systems/books/embedded-c-coding-standard

// Naming Convention Rules (BARR-C Rule 2)
const (
	RuleFunctionNaming = "barrc-naming-function"
	RuleVariableNaming = "barrc-naming-variable"
	RuleConstantNaming = "barrc-naming-constant"
)

// Type Safety Rules (BARR-C Rule 3)
const (
	RuleStdintTypes = "barrc-types-stdint"
)

// Magic Number Rules (BARR-C Rule 4)
const (
	RuleMagicNumbers = "barrc-safety-magic-numbers"
)

// Const Correctness Rules (BARR-C Rule 5)
const (
	RuleConstPointers = "barrc-safety-const-pointers"
)

// Rule definitions with descriptive metadata
var (
	// Function naming: Module_Function() format
	FunctionNamingRule = Rule{
		ID:       RuleFunctionNaming,
		Name:     "Function Naming Convention",
		Severity: "error",
		Category: "naming",
	}

	// Variable naming: snake_case
	VariableNamingRule = Rule{
		ID:       RuleVariableNaming,
		Name:     "Variable Naming Convention",
		Severity: "error",
		Category: "naming",
	}

	// Constant naming: UPPER_SNAKE
	ConstantNamingRule = Rule{
		ID:       RuleConstantNaming,
		Name:     "Constant Naming Convention",
		Severity: "error",
		Category: "naming",
	}

	// Type safety: Prefer stdint.h types (uint8_t over unsigned char)
	StdintTypesRule = Rule{
		ID:       RuleStdintTypes,
		Name:     "Stdint Type Usage",
		Severity: "error",
		Category: "types",
	}

	// No magic numbers: All constants must be #define (except 0, 1)
	MagicNumbersRule = Rule{
		ID:       RuleMagicNumbers,
		Name:     "Magic Number Detection",
		Severity: "warning",
		Category: "safety",
	}

	// Const correctness: const uint8_t* not uint8_t*
	ConstPointersRule = Rule{
		ID:       RuleConstPointers,
		Name:     "Const Pointer Correctness",
		Severity: "warning",
		Category: "safety",
	}
)

// AllRules returns all defined BARR-C rules
func AllRules() []Rule {
	return []Rule{
		FunctionNamingRule,
		VariableNamingRule,
		ConstantNamingRule,
		StdintTypesRule,
		MagicNumbersRule,
		ConstPointersRule,
	}
}
