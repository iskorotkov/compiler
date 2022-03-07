package option_test

import (
	"fmt"
	"testing"

	"github.com/iskorotkov/compiler/internal/fn/option"
	"github.com/stretchr/testify/assert"
)

func TestOptionWithOk(t *testing.T) {
t.Parallel()

	opt := option.Ok[int, error](123)

	assert.Equal(t, "ok: 123", opt.String())
}

func TestOptionWithErr(t *testing.T) {
t.Parallel()

	opt := option.Err[int](fmt.Errorf("test error"))

	assert.Equal(t, "err: test error", opt.String())
}
