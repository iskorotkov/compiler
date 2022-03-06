package literal

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/fn/option"
)

type LineNumber int

type ColNumber int

type Position struct {
	Line     LineNumber
	StartCol ColNumber
	EndCol   ColNumber
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d-%d", p.Line, p.StartCol, p.EndCol)
}

type Literal struct {
	Value    string
	Position Position
}

type Option = option.Option[Literal, error]

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

func Ok(literal Literal) Option {
	return option.Ok[Literal, error](literal)
}

func Err(err error) option.Option[Literal, error] {
	return option.Err[Literal](err)
}
