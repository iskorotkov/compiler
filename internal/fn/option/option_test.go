package option_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iskorotkov/compiler/internal/fn/option"
)

func TestOptionWithOk(t *testing.T) {
	t.Parallel()

	opt := option.Ok(123)

	assert.Equal(t, "ok: 123", opt.String())
}

func TestOptionWithErr(t *testing.T) {
	t.Parallel()

	opt := option.Err[int](fmt.Errorf("test error"))

	assert.Equal(t, "err: test error", opt.String())
}
