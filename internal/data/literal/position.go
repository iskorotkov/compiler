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
	return fmt.Sprintf("%d:%d-%d", p.Line, p.StartCol, p.EndCol)
}
