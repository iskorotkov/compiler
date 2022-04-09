package syntax_neutralizer

import (
	"errors"
	"fmt"
	"math/rand"

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
	case token.Unknown:
		return actual, fmt.Errorf("expecting unknown token: %w", UnfixableError)
	case token.UserDefined:
		actual.ID = token.UserDefined
		actual.Value = fmt.Sprintf("var%d", rand.Int())
		return actual, fmt.Errorf("expected %v, got %v: %w", expected, actual, FixableError)
	case token.EOF:
		return actual, fmt.Errorf("expecting %v, got %v: %w", expected, actual, UnfixableError)
	case token.IntLiteral:
		actual.ID = token.IntLiteral
		actual.Value = "0"
		return actual, fmt.Errorf("expected %v, got %v: %w", expected, actual, FixableError)
	case token.DoubleLiteral:
		actual.ID = token.DoubleLiteral
		actual.Value = "0.0"
		return actual, fmt.Errorf("expected %v, got %v: %w", expected, actual, FixableError)
	case token.BoolLiteral:
		actual.ID = token.BoolLiteral
		actual.Value = "false"
		return actual, fmt.Errorf("expected %v, got %v: %w", expected, actual, FixableError)
	default:
		expectedTokenValue := token.ByID(expected)

		dist := levenshtein.ComputeDistance(expectedTokenValue, actual.Value)
		if dist <= n.maxDistance {
			actual.ID = expected
			actual.Value = expectedTokenValue
			return actual, fmt.Errorf("expected %v, got %v: %w", expected, actual, FixableError)
		}

		return actual, fmt.Errorf("expected %v, got %v: %w", expected, actual, UnfixableError)
	}
}
