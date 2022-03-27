package syntax_analyzer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iskorotkov/compiler/internal/analyzers/syntax_analyzer"
	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
	"github.com/iskorotkov/compiler/internal/snapshot"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	type Test struct {
		name   string
		tokens []option.Option[token.Token]
	}

	tests := []Test{
		{
			name: "empty program",
			tokens: []option.Option[token.Token]{
				option.Ok(token.Token{ID: token.Program}),
				option.Ok(token.Token{ID: token.UserDefined}),
				option.Ok(token.Token{ID: token.Semicolon}),
				option.Ok(token.Token{ID: token.Begin}),
				option.Ok(token.Token{ID: token.End}),
				option.Ok(token.Token{ID: token.Period}),
				option.Ok(token.Token{ID: token.EOF}),
			},
		},
	}

	analyzer := syntax_analyzer.New(0)

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			resCh := analyzer.Analyze(channel.FromSlice(test.tokens))
			res := channel.ToSlice(resCh)

			actual := snapshot.NewSlice(res)
			expected := snapshot.Load(test.name)

			if !expected.Available() {
				actual.Save(test.name)
				return
			}

			assert.Equal(t, expected, actual)
		})
	}
}
