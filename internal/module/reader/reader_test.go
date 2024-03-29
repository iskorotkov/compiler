package reader_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iskorotkov/compiler/internal/data/literal"
	"github.com/iskorotkov/compiler/internal/fn/channel"
	"github.com/iskorotkov/compiler/internal/module/reader"
	"github.com/iskorotkov/compiler/internal/snapshot"
	"github.com/iskorotkov/compiler/testdata/programs"
)

func TestReader_Read(t *testing.T) {
	t.Parallel()

	type Test struct {
		name     string
		input    string
		expected []literal.Literal
	}

	type Config struct {
		name       string
		bufferSize int
	}

	tests := []Test{
		{
			name:  "sequence without whitespace",
			input: "123;,zcxc,t,,czxc,",
			expected: []literal.Literal{
				literal.New("123", 1, 1, 4),
				literal.New(";", 1, 4, 5),
				literal.New(",", 1, 5, 6),
				literal.New("zcxc", 1, 6, 10),
				literal.New(",", 1, 10, 11),
				literal.New("t", 1, 11, 12),
				literal.New(",", 1, 12, 13),
				literal.New(",", 1, 13, 14),
				literal.New("czxc", 1, 14, 18),
				literal.New(",", 1, 18, 19),
				literal.New("\n", 1, 19, 20),
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
			expected: []literal.Literal{
				literal.New("asd", 1, 1, 4),
				literal.New("\t", 1, 4, 5),
				literal.New("\n", 1, 5, 6),
				literal.New("sa", 2, 1, 3),
				literal.New("\n", 2, 3, 4),
				literal.New("a__sd21s", 3, 1, 9),
				literal.New("\v", 3, 9, 10),
				literal.New("123", 3, 10, 13),
				literal.New("\n", 3, 13, 14),
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
		config := config
		t.Run(config.name, func(t *testing.T) {
			t.Parallel()

			r := reader.New(config.bufferSize)
			for _, test := range tests {
				test := test
				t.Run(test.name, func(t *testing.T) {
					t.Parallel()

					actual := channel.ToSlice(r.Read(nil, strings.NewReader(test.input)))
					assert.Equal(t, test.expected, actual)
				})
			}
		})
	}
}

func TestReader_ReadWithSnapshots(t *testing.T) {
	t.Parallel()

	type Test struct {
		name  string
		input string
	}

	tests := []Test{
		{
			name:  "assignments program",
			input: programs.Assignments,
		},
		{
			name:  "constants program",
			input: programs.Constants,
		},
	}

	r := reader.New(0)

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := channel.ToSlice(r.Read(nil, strings.NewReader(test.input)))
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
