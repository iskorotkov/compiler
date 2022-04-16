package typechecker

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
	"github.com/iskorotkov/compiler/internal/data/symbol"
	"github.com/iskorotkov/compiler/internal/data/token"
)

type TypeResolver struct {
	converter TypeConverter
}

func NewTypeResolver(converter TypeConverter) TypeResolver {
	return TypeResolver{converter: converter}
}

func (r TypeResolver) ResolveExpr(
	ctx interface {
		context.LoggerContext
		context.SymbolScopeContext
	},
	expr ast.Node,
) (symbol.BuiltinType, error) {
	// TODO: Split expression into subexpressions when &&, ||, <, >, <=, >=, =, <>, etc. are encountered.
	linearized := linearizeExpressionTokens(expr)

	var valuesOnly []*ast.Leaf
	for _, item := range linearized {
		switch item.ID {
		case token.UserDefined, token.IntLiteral, token.DoubleLiteral, token.BoolLiteral:
			valuesOnly = append(valuesOnly, item)
		}
	}

	return subExpressionType(ctx, r.converter, valuesOnly)
}

func subExpressionType(
	ctx interface {
		context.LoggerContext
		context.SymbolScopeContext
	},
	converter TypeConverter,
	linearized []*ast.Leaf,
) (symbol.BuiltinType, error) {
	currentType := symbol.BuiltinTypeUnknown
	for _, leaf := range linearized {
		leafType, err := getLeafType(ctx, leaf)
		if err != nil {
			return symbol.BuiltinTypeUnknown, err
		}

		newType, ok := converter.Convert(ctx, leafType, currentType)
		if !ok {
			return symbol.BuiltinTypeUnknown, fmt.Errorf("can't assign type %s to %s", leafType, currentType)
		}

		currentType = newType
	}

	return currentType, nil
}

func getLeafType(
	ctx interface {
		context.SymbolScopeContext
	},
	leaf *ast.Leaf,
) (symbol.BuiltinType, error) {
	switch leaf.ID {
	case token.IntLiteral:
		return symbol.BuiltinTypeInt, nil
	case token.DoubleLiteral:
		return symbol.BuiltinTypeDouble, nil
	case token.BoolLiteral:
		return symbol.BuiltinTypeBool, nil
	case token.UserDefined:
		s, ok := ctx.SymbolScope().Lookup(&symbol.Name{Name: leaf.Value})
		if !ok {
			return symbol.BuiltinTypeUnknown, fmt.Errorf("symbol %s not found", leaf.Value)
		}

		switch s := s.(type) {
		case *symbol.Var:
			return s.Type.BuiltinType, nil
		case *symbol.Const:
			return s.Type.BuiltinType, nil
		default:
			return symbol.BuiltinTypeUnknown, fmt.Errorf("expected Type, got %T", s)
		}
	default:
		return symbol.BuiltinTypeUnknown, fmt.Errorf("unexpected token id %v", leaf.ID)
	}
}

func linearizeExpressionTokens(expr ast.Node) []*ast.Leaf {
	var tokens []*ast.Leaf

	// DFS.
	stack := []ast.Node{expr}
	for len(stack) != 0 {
		var generation []ast.Node
		for _, e := range stack {
			switch t := e.(type) {
			case *ast.Branch:
				// Traverse all branch children.
				for _, item := range t.Items {
					generation = append(generation, item)
				}
			case *ast.Leaf:
				tokens = append(tokens, t)
			}
		}

		stack = generation
	}

	return tokens
}
