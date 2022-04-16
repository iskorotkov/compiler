package literal

import "fmt"

type LineNumber int

type ColNumber int

type Position struct {
	Line     LineNumber
	StartCol ColNumber
	EndCol   ColNumber
}

func (p Position) String() string {
	if p.EndCol == p.StartCol+1 {
		return fmt.Sprintf("%d:%d", p.Line, p.StartCol)
	}

	return fmt.Sprintf("%d:%d-%d", p.Line, p.StartCol, p.EndCol)
}

func (p Position) Before(other Position) bool {
	return p.Line < other.Line || (p.Line == other.Line && p.StartCol < other.StartCol)
}

func (p Position) After(other Position) bool {
	return p.Line > other.Line || (p.Line == other.Line && p.StartCol > other.StartCol)
}

func (p Position) Join(other Position) Position {
	return Position{
		Line:     p.Line,
		StartCol: p.StartCol,
		EndCol:   other.EndCol,
	}
}
