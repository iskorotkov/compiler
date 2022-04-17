package ast

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/data/literal"
	"github.com/iskorotkov/compiler/internal/data/token"
)

const (
	QueryTypeOne QueryType = iota
	QueryTypeTop
	QueryTypeRecursive
)

var (
	_ Node = (*Leaf)(nil)
	_ Node = (*Branch)(nil)
)

type QueryType int

type Node interface {
	Query(queryType QueryType, markers ...Marker) []Node
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

func (l *Leaf) Query(_ QueryType, markers ...Marker) []Node {
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

func (b *Branch) Query(queryType QueryType, markers ...Marker) []Node {
	var res []Node
	for _, marker := range markers {
		if b.Has(marker) {
			res = append(res, b)

			// If we want single node or top nodes only, we can stop here.
			if queryType == QueryTypeOne || queryType == QueryTypeTop {
				return res
			}

			break
		}
	}

	for _, item := range b.Items {
		res = append(res, item.Query(queryType, markers...)...)

		// If we want single node only, we can stop here.
		if len(res) != 0 && queryType == QueryTypeOne {
			return res
		}
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
