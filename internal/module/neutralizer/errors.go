package neutralizer

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/data/token"
)

var (
	_ error = (*FixableError)(nil)
	_ error = (*UnfixableError)(nil)
)

type FixableError struct {
	Expected token.ID
	Actual   token.Token
}

func (e *FixableError) Error() string {
	return fmt.Sprintf("replaced %v with %v", e.Actual, e.Expected)
}

func (e *FixableError) Is(other error) bool {
	_, ok := other.(*FixableError)
	return ok
}

type UnfixableError struct {
	Expected token.ID
	Actual   token.Token
}

func (e *UnfixableError) Error() string {
	return fmt.Sprintf("unfixable syntax error: expected %v, got %v", e.Expected, e.Actual)
}

func (e *UnfixableError) Is(other error) bool {
	_, ok := other.(*UnfixableError)
	return ok
}
