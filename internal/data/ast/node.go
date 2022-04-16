package ast

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/data/literal"
	"github.com/iskorotkov/compiler/internal/data/token"
)

var (
	_ Node = (*Leaf)(nil)
	_ Node = (*Branch)(nil)
)

type Node interface {
	Query(markers ...Marker) []Node
	Has(marker Marker) bool
	Position() literal.Position
	fmt.Stringer
}

func Token(t token.Token, markers Markers) Node {
	return &Leaf{
		Token:   t,
		Markers: markers,
	}
}

func Wrap(node Node, markers Markers) Node {
	switch node := node.(type) {
	case *Leaf:
		return &Leaf{
			Token:   node.Token,
			Markers: node.Markers.Merge(markers),
		}
	case *Branch:
		return &Branch{
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

	return &Branch{
		Items:   nodes,
		Markers: markers,
	}
}

type Leaf struct {
	token.Token
	Markers
}

func (l *Leaf) Query(markers ...Marker) []Node {
	for _, marker := range markers {
		if l.Has(marker) {
			return []Node{l}
		}
	}

	return nil
}

func (l *Leaf) Position() literal.Position {
	return l.Token.Position
}

type Branch struct {
	Items []Node
	Markers
}

func (b *Branch) Query(markers ...Marker) []Node {
	for _, marker := range markers {
		if b.Has(marker) {
			// By default, we don't want to descend into the children of a branch.
			return []Node{b}
		}
	}

	var res []Node
	for _, item := range b.Items {
		res = append(res, item.Query(markers...)...)
	}

	return res
}

func (b *Branch) Position() literal.Position {
	switch len(b.Items) {
	case 0:
		return literal.Position{}
	case 1:
		return b.Items[0].Position()
	default:
		return b.Items[0].Position().Join(b.Items[len(b.Items)-1].Position())
	}
}

func (b *Branch) String() string {
	return fmt.Sprintf("list of %d items", len(b.Items))
}
