package neutralizer

import (
	"fmt"

	"github.com/agnivade/levenshtein"

	"github.com/iskorotkov/compiler/internal/data/symbol"
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

func (n Neutralizer) NeutralizeKeyword(expected token.ID, actual token.Token) (token.Token, error) {
	// Actual token matches expected token.
	if expected == actual.ID {
		return actual, nil
	}

	switch expected {
	case token.Unknown, token.UserDefined, token.EOF, token.IntLiteral, token.DoubleLiteral, token.BoolLiteral:
		return actual, &UnfixableKeywordError{
			Expected: expected,
			Actual:   actual,
		}
	default:
		expectedTokenValue := token.ByID(expected)

		// Avoid overwriting entire token value.
		if len(expectedTokenValue) < 3 || len(expectedTokenValue) <= n.maxDistance {
			return actual, &UnfixableKeywordError{
				Expected: expected,
				Actual:   actual,
			}
		}

		dist := levenshtein.ComputeDistance(expectedTokenValue, actual.Value)
		if dist > n.maxDistance {
			return actual, &UnfixableKeywordError{
				Expected: expected,
				Actual:   actual,
			}
		}

		neutralized := actual
		neutralized.ID = expected
		neutralized.Value = expectedTokenValue

		return neutralized, &FixableKeywordError{
			Expected: expected,
			Actual:   actual,
		}
	}
}

func (n Neutralizer) NeutralizeUserDefined(scope symbol.Scope, actual string) (symbol.Symbol, error) {
	s, ok := scope.Lookup(&symbol.Name{Name: actual})
	if ok {
		return s, nil
	}

	allSymbols := scope.Symbols()
	for _, savedSymbol := range allSymbols {
		var symbolName string
		switch savedSymbol := savedSymbol.(type) {
		case *symbol.Const:
			symbolName = savedSymbol.RawValue
		case *symbol.Var:
			symbolName = savedSymbol.Value
		case *symbol.Func:
			symbolName = savedSymbol.Value
		case *symbol.Type:
			symbolName = savedSymbol.Value
		default:
			return nil, fmt.Errorf("unexpected symbol type: %T", savedSymbol)
		}

		dist := levenshtein.ComputeDistance(actual, symbolName)
		if dist > n.maxDistance {
			continue
		}

		return savedSymbol, fmt.Errorf("replace %s with %s", actual, symbolName)
	}

	return nil, fmt.Errorf("no symbol found for %s", actual)
}
