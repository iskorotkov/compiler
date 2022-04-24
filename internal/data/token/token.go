package token

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/data/literal"
)

type Token struct {
	ID ID
	literal.Literal
}

func New(id ID, lit literal.Literal) Token {
	return Token{
		ID:      id,
		Literal: lit,
	}
}

func (t Token) String() string {
	switch t.ID {
	case UserDefined, IntLiteral, DoubleLiteral, BoolLiteral, EOF:
		return fmt.Sprintf("%v %v", t.ID, t.Literal)
	default:
		return t.Literal.String()
	}
}
