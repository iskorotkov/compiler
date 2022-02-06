package literal

import (
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var Factory = option.Factory[Literal, error]{}

type Option = option.Option[Literal, error]
