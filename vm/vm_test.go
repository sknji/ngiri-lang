package vm

import (
	"fmt"
	"github.com/marmotini/monkey-lang/ast"
	"github.com/marmotini/monkey-lang/compiler"
	"github.com/marmotini/monkey-lang/lexer"
	"github.com/marmotini/monkey-lang/object"
	"github.com/marmotini/monkey-lang/parser"
	"testing"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
	}

	runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		comp := compiler.NewCompiler()
		err := comp.Compile(parse(tt.input))
		if err != nil {
			t.Fatalf("Compiler error: %s", err)
		}

		vm := NewVM(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.StackTop()

		testExpectedObject(t, tt.expected, stackElem)
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	}
}
func parse(input string) *ast.Program {
	return parser.NewParser(lexer.NewLexer(input)).ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	return nil
}
