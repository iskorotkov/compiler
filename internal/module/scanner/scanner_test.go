package scanner_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iskorotkov/compiler/internal/data/literal"
	"github.com/iskorotkov/compiler/internal/fn/channel"
	"github.com/iskorotkov/compiler/internal/module/scanner"
	"github.com/iskorotkov/compiler/internal/snapshot"
)

func TestScanner_Scan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []literal.Literal
	}{
		{
			name:  "simple for loop",
			input: fromString("for i := 1 to 10 do writeln ( i )"),
		},
	}

	sc := scanner.New(0)

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := channel.ToSlice(sc.Scan(nil, channel.FromSlice(test.input)))
			s := snapshot.NewSlice(actual)

			expected := snapshot.Load(test.name)
			if !expected.Available() {
				s.Save(test.name)
				return
			}

			assert.Equal(t, expected, s)
		})
	}
}

func fromString(s string) []literal.Literal {
	var literals []literal.Literal
	for _, part := range strings.Split(s, " ") {
		literals = append(literals, literal.New(part, 0, 0, 0))
	}

	return literals
}
