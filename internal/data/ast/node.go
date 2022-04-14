package ast

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/data/token"
)

var (
	_ Node = Leaf{}
	_ Node = Branch{}
)

type Node interface {
	Query(markers ...Marker) []Node
	Has(marker Marker) bool
	fmt.Stringer
}

func Token(t token.Token, markers Markers) Node {
	return Leaf{
		Token:   t,
		Markers: markers,
	}
}

func Wrap(node Node, markers Markers) Node {
	switch node := node.(type) {
	case Leaf:
		return Leaf{
			Token:   node.Token,
			Markers: node.Markers.Merge(markers),
		}
	case Branch:
		return Branch{
			Items:   node.Items,
			Markers: node.Markers.Merge(markers),
		}
	default:
		panic("unknown node type")
	}
}

func WrapSlice(nodes []Node, markers Markers) Node {
	if len(nodes) == 0 {
		return nil
	}

	if len(nodes) == 1 {
		return Wrap(nodes[0], markers)
	}

	return Branch{
		Items:   nodes,
		Markers: markers,
	}
}

type Leaf struct {
	token.Token
	Markers
}

func (t Leaf) Query(markers ...Marker) []Node {
	for _, marker := range markers {
		if t.Has(marker) {
			return []Node{t}
		}
	}

	return nil
}

type Branch struct {
	Items []Node
	Markers
}

func (l Branch) Query(markers ...Marker) []Node {
	for _, marker := range markers {
		if l.Has(marker) {
			// By default, we don't want to descend into the children of a branch.
			return []Node{l}
		}
	}

	var res []Node
	for _, item := range l.Items {
		res = append(res, item.Query(markers...)...)
	}

	return res
}

func (l Branch) String() string {
	return fmt.Sprintf("list of %d items", len(l.Items))
}
