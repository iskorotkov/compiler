package context

import (
	"github.com/iskorotkov/compiler/internal/module/syntax_neutralizer"
)

var _ NeutralizerContext = (*neutralizerContext)(nil)

type neutralizerContext struct {
	neutralizer syntax_neutralizer.Neutralizer
}

func (n *neutralizerContext) Neutralizer() syntax_neutralizer.Neutralizer {
	return n.neutralizer
}
