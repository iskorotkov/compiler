package syntax_neutralizer

import (
	"errors"
	"fmt"

	"github.com/agnivade/levenshtein"

	"github.com/iskorotkov/compiler/internal/data/token"
)

var (
	FixableError   = errors.New("fixable error")
	UnfixableError = errors.New("unfixable error")
)

type Neutralizer struct {
	maxDistance int
}

func New(maxDistance int) Neutralizer {
	return Neutralizer{
		maxDistance: maxDistance,
	}
}

func (n *Neutralizer) Neutralize(expected token.ID, actual token.Token) (token.Token, error) {
	// Actual token matches expected token.
	if expected == actual.ID {
		return actual, nil
	}

	switch expected {
	case token.Unknown, token.UserDefined, token.EOF, token.IntLiteral, token.DoubleLiteral, token.BoolLiteral:
		return actual, fmt.Errorf("expecting %v, got %v: %w", expected, actual, UnfixableError)
	default:
		expectedTokenValue := token.ByID(expected)

		// Avoid overwriting entire token value.
		if len(expectedTokenValue) < 3 || len(expectedTokenValue) <= n.maxDistance {
			return actual, fmt.Errorf("expected %v, got %v: %w", expected, actual, UnfixableError)
		}

		dist := levenshtein.ComputeDistance(expectedTokenValue, actual.Value)
		if dist > n.maxDistance {
			return actual, fmt.Errorf("expected %v, got %v: %w", expected, actual, UnfixableError)
		}

		// TODO: Store info about neutralized errors somewhere.
		// TODO: Fix typos in "end" keyword.
		actual.ID = expected
		actual.Value = expectedTokenValue

		return actual, fmt.Errorf("expected %v, got %v: %w", expected, actual, FixableError)
	}
}
