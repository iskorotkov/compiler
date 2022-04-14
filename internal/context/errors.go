package context

import (
	"sync"

	"github.com/iskorotkov/compiler/internal/data/literal"
)

var _ ErrorsContext = (*errorsContext)(nil)

type Error struct {
	literal.Position
	error
}

type errorsContext struct {
	errors []Error
	m      sync.Mutex
}

func (e *errorsContext) AddError(position literal.Position, err error) {
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

	return e.errors
}
