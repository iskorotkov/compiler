package context

import (
	"context"
	"os"
)

func NewEnvContext(ctx context.Context) FullContext {
	if os.Getenv("DEBUG") == "1" {
		return NewDevContext(ctx)
	}

	return NewProdContext(ctx)
}
