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
