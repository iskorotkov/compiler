package typechecker

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
	"github.com/iskorotkov/compiler/internal/data/symbol"
	"github.com/iskorotkov/compiler/internal/data/token"
)

var opGroups = []opGroup{
	{
		tokens: map[token.ID]bool{
			token.And: true,
			token.Or:  true,
			token.Xor: true,
		},
		expectedType:  symbol.BuiltinTypeBool,
		resultingType: symbol.BuiltinTypeBool,
	},
	{
		tokens: map[token.ID]bool{
			token.Eq:  true,
			token.Ne:  true,
			token.Lt:  true,
			token.Lte: true,
			token.Gt:  true,
			token.Gte: true,
		},
		expectedType:  symbol.BuiltinTypeUnknown,
		resultingType: symbol.BuiltinTypeBool,
	},
}

type opGroup struct {
	// Which tokens are in the group.
	tokens map[token.ID]bool
	// Expected type of the left and right operands.
	expectedType symbol.BuiltinType
	// Resulting type of the operation.
	resultingType symbol.BuiltinType
}

type TypeResolver struct {
	converter TypeConverter
}

func NewTypeResolver(converter TypeConverter) TypeResolver {
	return TypeResolver{converter: converter}
}

func (r TypeResolver) Resolve(
	ctx interface {
		context.LoggerContext
		context.ErrorsContext
		context.NeutralizerContext
	},
	scope symbol.Scope,
	expr ast.Node,
) (symbol.BuiltinType, error) {
	leafs := linearizeExpressionTokens(expr)
	return r.resolveLinearExpr(ctx, scope, leafs)
}

func (r TypeResolver) resolveLinearExpr(
	ctx interface {
		context.LoggerContext
		context.ErrorsContext
		context.NeutralizerContext
	},
	scope symbol.Scope,
	leafs []*ast.Leaf,
) (symbol.BuiltinType, error) {
	for _, group := range opGroups {
		for i, leaf := range leafs {
			if leaf.ID == token.Not {
				t, err := getLeafType(ctx, scope, leafs[i+1])
				if err != nil {
					return symbol.BuiltinTypeUnknown, err
				}

				if t != symbol.BuiltinTypeBool {
					return symbol.BuiltinTypeUnknown, fmt.Errorf("unsupported type for not: %s", t)
				}
			}

			if !group.tokens[leaf.ID] {
				continue
			}

			leftType, err := r.resolveLinearExpr(ctx, scope, leafs[:i])
			if err != nil {
				return leftType, err
			}

			rightType, err := r.resolveLinearExpr(ctx, scope, leafs[i+1:])
			if err != nil {
				return rightType, err
			}

			// Check if left operand has type compatible with given operation.
			if _, ok := r.converter.IsAssignable(leftType, group.expectedType); !ok {
				return symbol.BuiltinTypeUnknown, fmt.Errorf("left operand has incompatible type %s", leftType)
			}

			// Check if right operand has type compatible with given operation.
			if _, ok := r.converter.IsAssignable(leftType, group.expectedType); !ok {
				return symbol.BuiltinTypeUnknown, fmt.Errorf("right operand has incompatible type %s", rightType)
			}

			// Check if operands have compatible types (can be cast one to another).
			if _, ok := r.converter.IsAssignable(rightType, leftType); !ok {
				return symbol.BuiltinTypeUnknown, fmt.Errorf("operands have incompatible types %s and %s", leftType, rightType)
			}

			if group.resultingType != symbol.BuiltinTypeUnknown {
				return group.resultingType, nil
			}

			return leftType, nil
		}
	}

	var valuesOnly []*ast.Leaf
	for _, item := range leafs {
		switch item.ID {
		case token.UserDefined, token.IntLiteral, token.DoubleLiteral, token.BoolLiteral:
			valuesOnly = append(valuesOnly, item)
		}
	}

	return r.subExpressionType(ctx, scope, valuesOnly)
}

func (r TypeResolver) subExpressionType(ctx interface {
	context.ErrorsContext
	context.NeutralizerContext
},
	scope symbol.Scope,
	linearized []*ast.Leaf,
) (symbol.BuiltinType, error) {
	currentType := symbol.BuiltinTypeUnknown
	for _, leaf := range linearized {
		leafType, err := getLeafType(ctx, scope, leaf)
		if err != nil {
			return symbol.BuiltinTypeUnknown, err
		}

		newType, ok := r.converter.IsAssignable(leafType, currentType)
		if !ok {
			return symbol.BuiltinTypeUnknown, fmt.Errorf("can't cast type %s to %s", leafType, currentType)
		}

		currentType = newType
	}

	return currentType, nil
}

func getLeafType(
	ctx interface {
		context.ErrorsContext
		context.NeutralizerContext
	},
	scope symbol.Scope,
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
		s, err := ctx.Neutralizer().NeutralizeUserDefined(scope, leaf.Value)
		if err != nil {
			return symbol.BuiltinTypeUnknown, err
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
	switch expr := expr.(type) {
	case *ast.Branch:
		// DFS.
		var tokens []*ast.Leaf
		for _, item := range expr.Items {
			tokens = append(tokens, linearizeExpressionTokens(item)...)
		}

		return tokens
	case *ast.Leaf:
		return []*ast.Leaf{expr}
	default:
		panic(fmt.Errorf("unexpected expression type %T", expr))
	}
}
