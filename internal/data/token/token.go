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
	if t.ID == UserDefined {
		return fmt.Sprintf("%v %v", t.ID, t.Literal)
	}

	return t.Literal.String()
}
