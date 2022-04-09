package contexts

import (
	"github.com/iskorotkov/compiler/internal/modules/syntax_neutralizer"
)

var _ NeutralizerContext = (*neutralizerContext)(nil)

type neutralizerContext struct {
	neutralizer syntax_neutralizer.Neutralizer
}

func NewNeutralizerContext(neutralizer syntax_neutralizer.Neutralizer) NeutralizerContext {
	return &neutralizerContext{neutralizer: neutralizer}
}

func (n *neutralizerContext) Neutralizer() syntax_neutralizer.Neutralizer {
	return n.neutralizer
}
