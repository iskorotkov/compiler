package literal

import (
	"github.com/iskorotkov/compiler/internal/fn/option"
)

type Literal struct {
	Value    string
	Position Position
}

func New(value string, line LineNumber, start, end ColNumber) Literal {
	return Literal{
		Value: value,
		Position: Position{
			Line:     line,
			StartCol: start,
			EndCol:   end,
		},
	}
}

type Option = option.Option[Literal, error]

func Ok(literal Literal) Option {
	return option.Ok[Literal, error](literal)
}

func Err(err error) option.Option[Literal, error] {
	return option.Err[Literal](err)
}
