package literal

type LineNumber int

type ColNumber int

type Position struct {
	Line     LineNumber
	StartCol ColNumber
	EndCol   ColNumber
}

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
