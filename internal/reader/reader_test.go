package reader_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iskorotkov/compiler/internal/fn"
	"github.com/iskorotkov/compiler/internal/literal"
	"github.com/iskorotkov/compiler/internal/reader"
)

func Test(t *testing.T) {
	t.Parallel()

	type Test struct {
		name     string
		input    string
		expected []reader.Element
	}

	type Config struct {
		name       string
		bufferSize int
	}

	tests := []Test{
		{
			name:  "sequence without whitespace",
			input: "123;,zcxc,t,,czxc,",
			expected: []reader.Element{
				{Value: literal.New("123", 1, 0, 3)},
				{Value: literal.New(";", 1, 3, 4)},
				{Value: literal.New(",", 1, 4, 5)},
				{Value: literal.New("zcxc", 1, 5, 9)},
				{Value: literal.New(",", 1, 9, 10)},
				{Value: literal.New("t", 1, 10, 11)},
				{Value: literal.New(",", 1, 11, 12)},
				{Value: literal.New(",", 1, 12, 13)},
				{Value: literal.New("czxc", 1, 13, 17)},
				{Value: literal.New(",", 1, 17, 18)},
				{Value: literal.New("\n", 1, 18, 19)},
			},
		},
		{
			name:     "empty sequence",
			input:    "",
			expected: nil,
		},
		{
			name:  "sequence with strange whitespace",
			input: "asd\t\nsa\r\na__sd21s\v123",
			expected: []reader.Element{
				{Value: literal.New("asd", 1, 0, 3)},
				{Value: literal.New("\t", 1, 3, 4)},
				{Value: literal.New("\n", 1, 4, 5)},
				{Value: literal.New("sa", 2, 0, 2)},
				{Value: literal.New("\n", 2, 2, 3)},
				{Value: literal.New("a__sd21s", 3, 0, 8)},
				{Value: literal.New("\v", 3, 8, 9)},
				{Value: literal.New("123", 3, 9, 12)},
				{Value: literal.New("\n", 3, 12, 13)},
			},
		},
	}

	configs := []Config{
		{
			name:       "running without buffer",
			bufferSize: 0,
		},
		{
			name:       "running with buffer size 1",
			bufferSize: 1,
		},
		{
			name:       "running with buffer size 10",
			bufferSize: 10,
		},
	}

	for _, config := range configs {
		r := reader.New(config.bufferSize)

		t.Run(config.name, func(t *testing.T) {
			t.Parallel()

			for _, test := range tests {
				test := test

				t.Run(test.name, func(t *testing.T) {
					t.Parallel()

					var actual []fn.ValueOrErr[literal.Literal, error]
					for value := range r.Read(strings.NewReader(test.input)) {
						actual = append(actual, value)
					}

					assert.Equal(t, test.expected, actual)
				})
			}
		})
	}
}
