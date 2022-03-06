package reader_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/literal"
	"github.com/iskorotkov/compiler/internal/reader"
	"github.com/iskorotkov/compiler/internal/snapshot"
	"github.com/iskorotkov/compiler/testdata"
)

func TestReader_Read(t *testing.T) {
	t.Parallel()

	type Test struct {
		name     string
		input    string
		expected []literal.Option
	}

	type Config struct {
		name       string
		bufferSize int
	}

	tests := []Test{
		{
			name:  "sequence without whitespace",
			input: "123;,zcxc,t,,czxc,",
			expected: []literal.Option{
				literal.Ok(literal.New("123", 1, 0, 3)),
				literal.Ok(literal.New(";", 1, 3, 4)),
				literal.Ok(literal.New(",", 1, 4, 5)),
				literal.Ok(literal.New("zcxc", 1, 5, 9)),
				literal.Ok(literal.New(",", 1, 9, 10)),
				literal.Ok(literal.New("t", 1, 10, 11)),
				literal.Ok(literal.New(",", 1, 11, 12)),
				literal.Ok(literal.New(",", 1, 12, 13)),
				literal.Ok(literal.New("czxc", 1, 13, 17)),
				literal.Ok(literal.New(",", 1, 17, 18)),
				literal.Ok(literal.New("\n", 1, 18, 19)),
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
			expected: []literal.Option{
				literal.Ok(literal.New("asd", 1, 0, 3)),
				literal.Ok(literal.New("\t", 1, 3, 4)),
				literal.Ok(literal.New("\n", 1, 4, 5)),
				literal.Ok(literal.New("sa", 2, 0, 2)),
				literal.Ok(literal.New("\n", 2, 2, 3)),
				literal.Ok(literal.New("a__sd21s", 3, 0, 8)),
				literal.Ok(literal.New("\v", 3, 8, 9)),
				literal.Ok(literal.New("123", 3, 9, 12)),
				literal.Ok(literal.New("\n", 3, 12, 13)),
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

					actual := channel.ToSlice(r.Read(strings.NewReader(test.input)))
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
			name:  "sample program 1",
			input: testdata.File1,
		},
		{
			name:  "sample program 2",
			input: testdata.File2,
		},
	}

	r := reader.New(0)

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := channel.ToSlice(r.Read(strings.NewReader(test.input)))
			s := snapshot.NewSlice(actual)

			expected := snapshot.Load(test.name)
			if expected == "" {
				s.Save(test.name)
				return
			}

			assert.Equal(t, expected, s)
		})
	}
}
