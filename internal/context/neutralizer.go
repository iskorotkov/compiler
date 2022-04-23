package context

import (
	"github.com/iskorotkov/compiler/internal/module/neutralizer"
)

var _ NeutralizerContext = (*neutralizerContext)(nil)

type neutralizerContext struct {
	neutralizer neutralizer.Neutralizer
}

func (n *neutralizerContext) Neutralizer() neutralizer.Neutralizer {
	return n.neutralizer
}
