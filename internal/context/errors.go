package context

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/iskorotkov/compiler/internal/data/literal"
)

const (
	ErrorSourceReader    ErrorSource = "reader"
	ErrorSourceScanner   ErrorSource = "scanner"
	ErrorSourceSyntax    ErrorSource = "syntax"
	ErrorSourceTypecheck ErrorSource = "typecheck"
	ErrorSourceCodegen   ErrorSource = "codegen"
	ErrorSourceInternal  ErrorSource = "internal"
)

type ErrorSource string

var _ ErrorsContext = (*errorsContext)(nil)

type Error struct {
	Source   ErrorSource
	Position literal.Position
	Err      error
}

func (e Error) Error() string {
	return fmt.Sprintf("%s %s: %v", strings.ToUpper(string(e.Source)), e.Position, e.Err)
}

type errorsContext struct {
	errors []Error
	m      sync.Mutex
}

func (e *errorsContext) AddError(source ErrorSource, position literal.Position, err error) {
	e.m.Lock()
	defer e.m.Unlock()

	e.errors = append(e.errors, Error{
		Source:   source,
		Position: position,
		Err:      err,
	})
}

func (e *errorsContext) Errors() []Error {
	e.m.Lock()
	defer e.m.Unlock()

	sort.SliceStable(e.errors, func(i, j int) bool {
		return e.errors[i].Position.Before(e.errors[j].Position)
	})

	return e.errors
}
