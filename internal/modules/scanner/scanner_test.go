package scanner_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iskorotkov/compiler/internal/data/literal"
	"github.com/iskorotkov/compiler/internal/fn/channels"
	"github.com/iskorotkov/compiler/internal/fn/options"
	"github.com/iskorotkov/compiler/internal/modules/scanner"
	"github.com/iskorotkov/compiler/internal/snapshots"
)

func TestScanner_Scan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []options.Option[literal.Literal]
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

			actual := channels.ToSlice(sc.Scan(channels.FromSlice(test.input)))
			s := snapshots.NewSlice(actual)

			expected := snapshots.Load(test.name)
			if !expected.Available() {
				s.Save(test.name)
				return
			}

			assert.Equal(t, expected, s)
		})
	}
}

func fromString(s string) []options.Option[literal.Literal] {
	var literals []options.Option[literal.Literal]
	for _, part := range strings.Split(s, " ") {
		literals = append(literals, options.Ok(literal.New(part, 0, 0, 0)))
	}

	return literals
}
