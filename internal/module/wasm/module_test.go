package wasm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModule_String(t *testing.T) {
	t.Parallel()

	type test struct {
		name     string
		expected string
		module   Module
	}

	tests := []test{
		{
			name:     "empty module",
			expected: "(module)",
			module: Module{
				Imports: nil,
				Globals: nil,
				Funcs:   nil,
			},
		},
		{
			name: "non-empty module",
			expected: `(module
  (import "console" "log" (func $log_i32 (param $0 i32)))
  (global $g (mut i32) (i32.const 42))
  (global $pi f64 (f64.const 3.141592653589793))
  (global $e f64 (f64.const 2.718281828459045))
)`,
			module: Module{
				Imports: []Import{
					{
						Path:   []string{"console", "log"},
						Name:   "log_i32",
						Params: []Param{{"0", "i32"}},
						Return: nil,
					},
				},
				Globals: []Global{
					{
						Name:    "g",
						Type:    "i32",
						Value:   "42",
						Mutable: true,
					},
					{
						Name:    "pi",
						Type:    "f64",
						Value:   "3.141592653589793",
						Mutable: false,
					},
					{
						Name:    "e",
						Type:    "f64",
						Value:   "2.718281828459045",
						Mutable: false,
					},
				},
				Funcs: nil,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, test.module.String())
		})
	}
}
