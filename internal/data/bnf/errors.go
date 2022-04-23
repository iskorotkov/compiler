package bnf

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/data/token"
)

var _ error = (*UnexpectedTokenError)(nil)

type UnexpectedTokenError struct {
	Expected token.ID
	Actual   token.Token
}

func (e *UnexpectedTokenError) Error() string {
	return fmt.Sprintf("unexpected token: expected %q, got %q", e.Expected, e.Actual.Value)
}

func (e *UnexpectedTokenError) Is(other error) bool {
	switch other.(type) {
	case *UnexpectedTokenError:
		return true
	default:
		return false
	}
}
