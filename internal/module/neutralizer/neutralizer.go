package neutralizer

import (
	"github.com/agnivade/levenshtein"

	"github.com/iskorotkov/compiler/internal/data/token"
)

type Neutralizer struct {
	maxDistance int
}

func New(maxDistance int) Neutralizer {
	return Neutralizer{
		maxDistance: maxDistance,
	}
}

func (n Neutralizer) Neutralize(expected token.ID, actual token.Token) (token.Token, error) {
	// Actual token matches expected token.
	if expected == actual.ID {
		return actual, nil
	}

	switch expected {
	case token.Unknown, token.UserDefined, token.EOF, token.IntLiteral, token.DoubleLiteral, token.BoolLiteral:
		return actual, &UnfixableError{
			Expected: expected,
			Actual:   actual,
		}
	default:
		expectedTokenValue := token.ByID(expected)

		// Avoid overwriting entire token value.
		if len(expectedTokenValue) < 3 || len(expectedTokenValue) <= n.maxDistance {
			return actual, &UnfixableError{
				Expected: expected,
				Actual:   actual,
			}
		}

		dist := levenshtein.ComputeDistance(expectedTokenValue, actual.Value)
		if dist > n.maxDistance {
			return actual, &UnfixableError{
				Expected: expected,
				Actual:   actual,
			}
		}

		actual.ID = expected
		actual.Value = expectedTokenValue

		return actual, &FixableError{
			Expected: expected,
			Actual:   actual,
		}
	}
}
