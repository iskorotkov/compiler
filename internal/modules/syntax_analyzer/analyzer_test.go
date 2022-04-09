package syntax_analyzer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iskorotkov/compiler/internal/contexts"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channels"
	"github.com/iskorotkov/compiler/internal/fn/options"
	"github.com/iskorotkov/compiler/internal/fn/slices"
	"github.com/iskorotkov/compiler/internal/modules/syntax_analyzer"
	"github.com/iskorotkov/compiler/internal/modules/syntax_neutralizer"
	"github.com/iskorotkov/compiler/internal/snapshots"
)

var (
	header = []options.Option[token.Token]{
		options.Ok(token.Token{ID: token.Program}),
		options.Ok(token.Token{ID: token.UserDefined}),
		options.Ok(token.Token{ID: token.Semicolon}),
		options.Ok(token.Token{ID: token.Begin}),
	}
	footer = []options.Option[token.Token]{
		options.Ok(token.Token{ID: token.End}),
		options.Ok(token.Token{ID: token.Period}),
		options.Ok(token.Token{ID: token.EOF}),
	}
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	type Test struct {
		name   string
		tokens []options.Option[token.Token]
	}

	tests := []Test{
		{
			name: "empty program",
			tokens: slices.Flatten(
				header,
				footer,
			),
		},
		{
			name: "unsigned int literal assignment",
			tokens: slices.Flatten(
				header,
				[]options.Option[token.Token]{
					options.Ok(token.Token{ID: token.UserDefined}),
					options.Ok(token.Token{ID: token.Assign}),
					options.Ok(token.Token{ID: token.IntLiteral}),
					options.Ok(token.Token{ID: token.Semicolon}),
				},
				footer,
			),
		},
		{
			name: "signed int literal assignment",
			tokens: slices.Flatten(
				header,
				[]options.Option[token.Token]{
					options.Ok(token.Token{ID: token.UserDefined}),
					options.Ok(token.Token{ID: token.Assign}),
					options.Ok(token.Token{ID: token.Minus}),
					options.Ok(token.Token{ID: token.IntLiteral}),
					options.Ok(token.Token{ID: token.Semicolon}),
				},
				footer,
			),
		},
		{
			name: "unsigned double literal assignment",
			tokens: slices.Flatten(
				header,
				[]options.Option[token.Token]{
					options.Ok(token.Token{ID: token.UserDefined}),
					options.Ok(token.Token{ID: token.Assign}),
					options.Ok(token.Token{ID: token.DoubleLiteral}),
					options.Ok(token.Token{ID: token.Semicolon}),
				},
				footer,
			),
		},
		{
			name: "signed double literal assignment",
			tokens: slices.Flatten(
				header,
				[]options.Option[token.Token]{
					options.Ok(token.Token{ID: token.UserDefined}),
					options.Ok(token.Token{ID: token.Assign}),
					options.Ok(token.Token{ID: token.Minus}),
					options.Ok(token.Token{ID: token.DoubleLiteral}),
					options.Ok(token.Token{ID: token.Semicolon}),
				},
				footer,
			),
		},
		{
			name: "bool literal assignment",
			tokens: slices.Flatten(
				header,
				[]options.Option[token.Token]{
					options.Ok(token.Token{ID: token.UserDefined}),
					options.Ok(token.Token{ID: token.Assign}),
					options.Ok(token.Token{ID: token.BoolLiteral}),
					options.Ok(token.Token{ID: token.Semicolon}),
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

			resCh := analyzer.Analyze(struct {
				contexts.LoggerContext
				contexts.NeutralizerContext
			}{
				LoggerContext:      contexts.NewEnvContext(context.Background()),
				NeutralizerContext: contexts.NewNeutralizerContext(syntax_neutralizer.New(0)),
			}, channels.FromSlice(test.tokens))

			res := channels.ToSlice(resCh)

			actual := snapshots.NewSlice(res)
			expected := snapshots.Load(test.name)

			if !expected.Available() {
				actual.Save(test.name)
				return
			}

			assert.Equal(t, expected, actual)
		})
	}
}
