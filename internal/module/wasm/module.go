package wasm

import (
	"strings"
)

type Module struct {
	Imports []Import
	Globals []Global
	Funcs   []Func
}

func (m Module) String() string {
	var s strings.Builder

	s.WriteString("(module")
	if len(m.Imports)+len(m.Globals)+len(m.Funcs) > 0 {
		s.WriteString("\n")
	}

	if len(m.Imports) > 0 {
		for _, i := range m.Imports {
			s.WriteString(i.StringIndent(1))
			s.WriteString("\n")
		}
	}

	if len(m.Globals) > 0 {
		for _, g := range m.Globals {
			s.WriteString(g.StringIndent(1))
			s.WriteString("\n")
		}
	}

	if len(m.Funcs) > 0 {
		for _, f := range m.Funcs {
			s.WriteString(f.StringIndent(1))
			s.WriteString("\n")
		}
	}

	s.WriteString(")")
	return s.String()
}
