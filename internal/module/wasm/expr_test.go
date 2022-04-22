package wasm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpr_String(t *testing.T) {
	t.Parallel()

	type test struct {
		name     string
		expected string
		input    Expr
	}

	tests := []test{
		{
			name:     "int constant",
			expected: "(i32.const 42)",
			input:    &Const{"i32", "42"},
		},
		{
			name:     "float constant",
			expected: "(f32.const 3.14)",
			input:    &Const{"f32", "3.14"},
		},
		{
			name: "add ints",
			expected: "(i32.add" +
				" (i32.const 1)" +
				" (i32.const 2))",
			input: &BinaryOp{"i32", OpAdd,
				&Const{"i32", "1"},
				&Const{"i32", "2"}},
		},
		{
			name: "mult ints",
			expected: "(i32.mul" +
				" (i32.const 1)" +
				" (i32.const 2))",
			input: &BinaryOp{"i32", OpMul,
				&Const{"i32", "1"},
				&Const{"i32", "2"}},
		},
		{
			name: "equal ints",
			expected: "(i32.eq" +
				" (i32.const 1)" +
				" (i32.const 2))",
			input: &BinaryOp{"i32", OpEq,
				&Const{"i32", "1"},
				&Const{"i32", "2"}},
		},
		{
			name: "compare ints",
			expected: "(i32.lt_s" +
				" (i32.const 1)" +
				" (i32.const 2))",
			input: &BinaryOp{"i32", OpLtSigned,
				&Const{"i32", "1"},
				&Const{"i32", "2"}},
		},
		{
			name: "add floats",
			expected: "(f32.add" +
				" (f32.const 1.0)" +
				" (f32.const 2.0))",
			input: &BinaryOp{"f32", OpAdd,
				&Const{"f32", "1.0"},
				&Const{"f32", "2.0"}},
		},
		{
			name: "mult floats",
			expected: "(f32.mul" +
				" (f32.const 1.0)" +
				" (f32.const 2.0))",
			input: &BinaryOp{"f32", OpMul,
				&Const{"f32", "1.0"},
				&Const{"f32", "2.0"}},
		},
		{
			name: "equal floats",
			expected: "(f32.eq" +
				" (f32.const 1.0)" +
				" (f32.const 2.0))",
			input: &BinaryOp{"f32", OpEq,
				&Const{"f32", "1.0"},
				&Const{"f32", "2.0"}},
		},
		{
			name: "compare floats",
			expected: "(f32.lt" +
				" (f32.const 1.0)" +
				" (f32.const 2.0))",
			input: &BinaryOp{"f32", OpLt,
				&Const{"f32", "1.0"},
				&Const{"f32", "2.0"}},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, test.input.String())
		})
	}
}
