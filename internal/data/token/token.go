package token

import (
	"fmt"
	"reflect"

	"github.com/iskorotkov/compiler/internal/constants"
	"github.com/iskorotkov/compiler/internal/data/literal"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

const (
	TypeUnknown Type = iota
	TypeConstant
	TypeUserIdentifier
	TypeKeyword
	TypePunctuation
	TypeOperator
)

type Type int

type Token struct {
	Type    Type
	ID      constants.ID
	Literal literal.Literal
	Value   *reflect.Value
}

func New(tokenType Type, id constants.ID, lit literal.Literal, value *reflect.Value) Token {
	return Token{
		Type:    tokenType,
		ID:      id,
		Literal: lit,
		Value:   value,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%d.%d: %v", t.Type, t.ID, t.Literal)
}

type Option = option.Option[Token, error]

func Ok(token Token) Option {
	return option.Ok[Token, error](token)
}

func Err(err error) option.Option[Token, error] {
	return option.Err[Token](err)
}
