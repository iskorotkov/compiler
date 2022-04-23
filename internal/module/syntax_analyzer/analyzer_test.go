package syntax_analyzer_test

import (
	stdcontext "context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channel"
	"github.com/iskorotkov/compiler/internal/fn/slice"
	"github.com/iskorotkov/compiler/internal/module/syntax_analyzer"
	"github.com/iskorotkov/compiler/internal/snapshot"
)

var (
	header = []token.Token{
		{ID: token.Program},
		{ID: token.UserDefined},
		{ID: token.Semicolon},
		{ID: token.Begin},
	}
	footer = []token.Token{
		{ID: token.End},
		{ID: token.Period},
		{ID: token.EOF},
	}
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	type Test struct {
		name   string
		tokens []token.Token
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
				[]token.Token{
					{ID: token.UserDefined},
					{ID: token.Assign},
					{ID: token.IntLiteral},
					{ID: token.Semicolon},
				},
				footer,
			),
		},
		{
			name: "signed int literal assignment",
			tokens: slice.Flatten(
				header,
				[]token.Token{
					{ID: token.UserDefined},
					{ID: token.Assign},
					{ID: token.Minus},
					{ID: token.IntLiteral},
					{ID: token.Semicolon},
				},
				footer,
			),
		},
		{
			name: "unsigned double literal assignment",
			tokens: slice.Flatten(
				header,
				[]token.Token{
					{ID: token.UserDefined},
					{ID: token.Assign},
					{ID: token.DoubleLiteral},
					{ID: token.Semicolon},
				},
				footer,
			),
		},
		{
			name: "signed double literal assignment",
			tokens: slice.Flatten(
				header,
				[]token.Token{
					{ID: token.UserDefined},
					{ID: token.Assign},
					{ID: token.Minus},
					{ID: token.DoubleLiteral},
					{ID: token.Semicolon},
				},
				footer,
			),
		},
		{
			name: "bool literal assignment",
			tokens: slice.Flatten(
				header,
				[]token.Token{
					{ID: token.UserDefined},
					{ID: token.Assign},
					{ID: token.BoolLiteral},
					{ID: token.Semicolon},
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
