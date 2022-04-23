package context

import (
	"fmt"
	"sort"
	"sync"

	"github.com/iskorotkov/compiler/internal/data/literal"
)

const (
	ErrorSourceReader ErrorSource = iota + 1
	ErrorSourceScanner
	ErrorSourceSyntax
	ErrorSourceTypecheck
	ErrorSourceCodegen
	ErrorSourceInternal
)

type ErrorSource int

var _ ErrorsContext = (*errorsContext)(nil)

type Error struct {
	literal.Position
	ErrorSource
	error
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Position, e.error)
}

type errorsContext struct {
	errors []Error
	m      sync.Mutex
}

func (e *errorsContext) AddError(source ErrorSource, position literal.Position, err error) {
	e.m.Lock()
	defer e.m.Unlock()

	e.errors = append(e.errors, Error{
		Position: position,
		error:    err,
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
