package neutralizer

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/data/token"
)

var (
	_ error = (*FixableKeywordError)(nil)
	_ error = (*UnfixableKeywordError)(nil)
)

type FixableKeywordError struct {
	Expected token.ID
	Actual   token.Token
}

func (e *FixableKeywordError) Error() string {
	return fmt.Sprintf("replace %v with %v", e.Actual.Value, e.Expected)
}

func (e *FixableKeywordError) Is(other error) bool {
	_, ok := other.(*FixableKeywordError)
	return ok
}

type UnfixableKeywordError struct {
	Expected token.ID
	Actual   token.Token
}

func (e *UnfixableKeywordError) Error() string {
	return fmt.Sprintf("unfixable syntax error: expected %v, got %v", e.Expected, e.Actual)
}

func (e *UnfixableKeywordError) Is(other error) bool {
	_, ok := other.(*UnfixableKeywordError)
	return ok
}
