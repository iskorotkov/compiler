package token

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/data/literal"
	"github.com/iskorotkov/compiler/internal/fn/option"
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
	return fmt.Sprintf("%v: %v", t.ID, t.Literal)
}

type Option = option.Option[Token, error]

func Ok(token Token) Option {
	return option.Ok[Token, error](token)
}

func Err(err error) option.Option[Token, error] {
	return option.Err[Token](err)
}
