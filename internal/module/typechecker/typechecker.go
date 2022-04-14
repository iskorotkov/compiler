package typechecker

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
	"github.com/iskorotkov/compiler/internal/data/symbol"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
	"github.com/iskorotkov/compiler/internal/module/syntax_neutralizer"
)

type TypeChecker struct {
	buffer      int
	neutralizer syntax_neutralizer.Neutralizer
}

func NewTypeChecker(buffer int) *TypeChecker {
	return &TypeChecker{
		buffer: buffer,
	}
}

func (a TypeChecker) Check(
	ctx interface {
		context.LoggerContext
		context.ErrorsContext
		context.NeutralizerContext
		context.SymbolScopeContext
	},
	input <-chan option.Option[ast.Node],
) <-chan option.Option[interface{}] {
	ch := make(chan option.Option[interface{}], a.buffer)

	go func() {
		defer close(ch)

		ctx.Logger().Infof("type checking started")

		for opt := range input {
			program, err := opt.Unwrap()
			if err != nil {
				ctx.Logger().Error(err)
				continue
			}

			addTypeDecls(ctx, program)
			addConstDecls(ctx, program)
			addVarDecls(ctx, program)
			addFuncDecls(ctx, program)

			// TODO: Remove expressions in function names.
			expressions := program.Query(ast.MarkerExpr)
			_ = expressions

			assignments := program.Query(ast.MarkerAssign)
			_ = assignments

			ctx.Logger().Infof("type checker found %d symbols", len(ctx.SymbolScope().Symbols()))
			for _, s := range ctx.SymbolScope().Symbols() {
				ctx.Logger().Infof("%d: %v", s.Hash(), s)
			}
		}

		ctx.Logger().Infof("type checking succeeded")
	}()

	return ch
}

func addFuncDecls(
	ctx interface {
		context.LoggerContext
		context.ErrorsContext
		context.NeutralizerContext
		context.SymbolScopeContext
	}, program ast.Node,
) {
	for _, decl := range program.Query(ast.MarkerFuncDecl) {
		name := decl.Query(ast.MarkerFuncName)[0].(ast.Leaf)
		returnType := decl.Query(ast.MarkerReturnType)[0].Query(ast.MarkerType)[0].(ast.Leaf)

		returnTypeSymbol, ok := ctx.SymbolScope().Lookup(&symbol.Name{Name: returnType.Value})
		if !ok {
			ctx.AddError(returnType.Position, fmt.Errorf("type %s not found", returnType.Value))
			continue
		}

		if _, ok := returnTypeSymbol.(*symbol.Type); !ok {
			ctx.AddError(returnType.Position, fmt.Errorf("symbol %s is not a type", returnType.Value))
			continue
		}

		var params []symbol.Var
		for _, param := range decl.Query(ast.MarkerParamGroupDecl) {
			for _, name := range param.Query(ast.MarkerName) {
				paramName := name.(ast.Leaf)
				paramType := param.Query(ast.MarkerType)[0].(ast.Leaf)
				paramTypeSymbol, ok := ctx.SymbolScope().Lookup(&symbol.Name{Name: paramType.Value})
				if !ok {
					ctx.AddError(paramType.Position, fmt.Errorf("type %s not found", paramType.Value))
					continue
				}

				if _, ok := paramTypeSymbol.(*symbol.Type); !ok {
					ctx.AddError(paramType.Position, fmt.Errorf("symbol %s is not a type", paramType.Value))
					continue
				}

				params = append(params, symbol.Var{
					Token:       paramName.Token,
					Type:        *paramTypeSymbol.(*symbol.Type),
					Initialized: false,
				})
			}
		}

		if err := ctx.SymbolScope().Add(&symbol.Func{
			Token:      name.Token,
			Params:     params,
			ReturnType: *returnTypeSymbol.(*symbol.Type),
		}); err != nil {
			ctx.AddError(name.Position, err)
			continue
		}

		// TODO: Analyze variables, constants and types defined in function.
	}
}

func addTypeDecls(
	ctx interface {
		context.LoggerContext
		context.ErrorsContext
		context.NeutralizerContext
		context.SymbolScopeContext
	}, program ast.Node,
) {
	for _, decl := range program.Query(ast.MarkerTypeDecl) {
		name := decl.Query(ast.MarkerName)[0].(ast.Leaf)
		typeName := decl.Query(ast.MarkerType)[0].(ast.Leaf)

		typeNameSymbol, ok := ctx.SymbolScope().Lookup(&symbol.Name{Name: typeName.Value})
		if !ok {
			ctx.AddError(typeName.Position, fmt.Errorf("type %s not found", typeName.Value))
			continue
		}

		if _, ok := typeNameSymbol.(*symbol.Type); !ok {
			ctx.AddError(typeName.Position, fmt.Errorf("symbol %s is not a type", typeName.Value))
			continue
		}

		if err := ctx.SymbolScope().Add(&symbol.Type{
			Token: name.Token,
			Alias: typeNameSymbol.(*symbol.Type),
		}); err != nil {
			ctx.AddError(name.Position, err)
			continue
		}
	}
}

func addConstDecls(
	ctx interface {
		context.LoggerContext
		context.ErrorsContext
		context.NeutralizerContext
		context.SymbolScopeContext
	}, program ast.Node,
) {
	for _, decl := range program.Query(ast.MarkerConstDecl) {
		name := decl.Query(ast.MarkerName)[0].(ast.Leaf)
		valueNode := decl.Query(ast.MarkerValue)[0].(ast.Leaf)

		var typeName string
		switch valueNode.ID {
		case token.IntLiteral:
			typeName = "integer"
		case token.DoubleLiteral:
			typeName = "real"
		case token.BoolLiteral:
			typeName = "boolean"
		default:
			ctx.AddError(valueNode.Position, fmt.Errorf("unsupported constant type %s", valueNode.ID))
			continue
		}

		typeSymbol, ok := ctx.SymbolScope().Lookup(&symbol.Name{Name: typeName})
		if !ok {
			ctx.AddError(valueNode.Position, fmt.Errorf("type %s not found", typeName))
			continue
		}

		if _, ok := typeSymbol.(*symbol.Type); !ok {
			ctx.AddError(valueNode.Position, fmt.Errorf("symbol %s is not a type", typeName))
			continue
		}

		if err := ctx.SymbolScope().Add(&symbol.Const{
			Token:    name.Token,
			Type:     *typeSymbol.(*symbol.Type),
			RawValue: valueNode.Value,
		}); err != nil {
			ctx.AddError(name.Position, err)
			continue
		}
	}
}

func addVarDecls(
	ctx interface {
		context.LoggerContext
		context.ErrorsContext
		context.NeutralizerContext
		context.SymbolScopeContext
	}, program ast.Node,
) {
	for _, decl := range program.Query(ast.MarkerVarDecl) {
		typeName := decl.Query(ast.MarkerType)[0].(ast.Leaf)
		typeNameSymbol, ok := ctx.SymbolScope().Lookup(&symbol.Name{Name: typeName.Value})
		if !ok {
			ctx.AddError(typeName.Position, fmt.Errorf("type %s not found", typeName.Value))
			continue
		}

		if _, ok := typeNameSymbol.(*symbol.Type); !ok {
			ctx.AddError(typeName.Position, fmt.Errorf("symbol %s is not a type", typeName.Value))
			continue
		}

		for _, name := range decl.Query(ast.MarkerName) {
			name := name.(ast.Leaf)

			if err := ctx.SymbolScope().Add(&symbol.Var{
				Token: name.Token,
				Type:  *typeNameSymbol.(*symbol.Type),
			}); err != nil {
				ctx.AddError(name.Position, err)
				continue
			}
		}
	}
}
