package syntax_analyzer_test

import (
	stdcontext "context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channel"
	"github.com/iskorotkov/compiler/internal/fn/option"
	"github.com/iskorotkov/compiler/internal/fn/slice"
	"github.com/iskorotkov/compiler/internal/module/syntax_analyzer"
	"github.com/iskorotkov/compiler/internal/snapshot"
)

var (
	header = []option.Option[token.Token]{
		option.Ok(token.Token{ID: token.Program}),
		option.Ok(token.Token{ID: token.UserDefined}),
		option.Ok(token.Token{ID: token.Semicolon}),
		option.Ok(token.Token{ID: token.Begin}),
	}
	footer = []option.Option[token.Token]{
		option.Ok(token.Token{ID: token.End}),
		option.Ok(token.Token{ID: token.Period}),
		option.Ok(token.Token{ID: token.EOF}),
	}
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
			tokens: slice.Flatten(
				header,
				footer,
			),
		},
		{
			name: "unsigned int literal assignment",
			tokens: slice.Flatten(
				header,
				[]option.Option[token.Token]{
					option.Ok(token.Token{ID: token.UserDefined}),
					option.Ok(token.Token{ID: token.Assign}),
					option.Ok(token.Token{ID: token.IntLiteral}),
					option.Ok(token.Token{ID: token.Semicolon}),
				},
				footer,
			),
		},
		{
			name: "signed int literal assignment",
			tokens: slice.Flatten(
				header,
				[]option.Option[token.Token]{
					option.Ok(token.Token{ID: token.UserDefined}),
					option.Ok(token.Token{ID: token.Assign}),
					option.Ok(token.Token{ID: token.Minus}),
					option.Ok(token.Token{ID: token.IntLiteral}),
					option.Ok(token.Token{ID: token.Semicolon}),
				},
				footer,
			),
		},
		{
			name: "unsigned double literal assignment",
			tokens: slice.Flatten(
				header,
				[]option.Option[token.Token]{
					option.Ok(token.Token{ID: token.UserDefined}),
					option.Ok(token.Token{ID: token.Assign}),
					option.Ok(token.Token{ID: token.DoubleLiteral}),
					option.Ok(token.Token{ID: token.Semicolon}),
				},
				footer,
			),
		},
		{
			name: "signed double literal assignment",
			tokens: slice.Flatten(
				header,
				[]option.Option[token.Token]{
					option.Ok(token.Token{ID: token.UserDefined}),
					option.Ok(token.Token{ID: token.Assign}),
					option.Ok(token.Token{ID: token.Minus}),
					option.Ok(token.Token{ID: token.DoubleLiteral}),
					option.Ok(token.Token{ID: token.Semicolon}),
				},
				footer,
			),
		},
		{
			name: "bool literal assignment",
			tokens: slice.Flatten(
				header,
				[]option.Option[token.Token]{
					option.Ok(token.Token{ID: token.UserDefined}),
					option.Ok(token.Token{ID: token.Assign}),
					option.Ok(token.Token{ID: token.BoolLiteral}),
					option.Ok(token.Token{ID: token.Semicolon}),
				},
				footer,
			),
		},
	}

	analyzer := syntax_analyzer.New(0)

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			resCh := analyzer.Analyze(context.NewEnvContext(stdcontext.Background()), channel.FromSlice(test.tokens))

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
