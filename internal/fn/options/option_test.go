package options_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iskorotkov/compiler/internal/fn/options"
)

func TestOptionWithOk(t *testing.T) {
	t.Parallel()

	opt := options.Ok(123)

	assert.Equal(t, "ok: 123", opt.String())
}

func TestOptionWithErr(t *testing.T) {
	t.Parallel()

	opt := options.Err[int](fmt.Errorf("test error"))

	assert.Equal(t, "err: test error", opt.String())
}
