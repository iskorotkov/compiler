package literal

import (
	"fmt"
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

func (l Literal) String() string {
	return fmt.Sprintf("%q at %v", l.Value, l.Position)
}
