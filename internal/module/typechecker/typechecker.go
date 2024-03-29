package typechecker

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
	"github.com/iskorotkov/compiler/internal/data/symbol"
	"github.com/iskorotkov/compiler/internal/data/token"
)

type Result struct {
	Node  ast.Node
	Scope symbol.Scope
}

type TypeChecker struct {
	buffer    int
	converter TypeConverter
	resolver  TypeResolver
}

func NewTypeChecker(buffer int) *TypeChecker {
	converter := NewTypeConverter()
	return &TypeChecker{
		buffer:    buffer,
		converter: converter,
		resolver:  NewTypeResolver(converter),
	}
}

func (c TypeChecker) Check(
	ctx interface {
		context.LoggerContext
		context.ErrorsContext
		context.NeutralizerContext
	},
	input <-chan ast.Node,
) <-chan Result {
	ch := make(chan Result, c.buffer)

	go func() {
		defer close(ch)

		ctx.Logger().Infof("type checking started")

		for program := range input {
			block := program.Query(ast.QueryTypeOne, ast.MarkerProgramBlock)[0]
			scope := symbol.NewScope()

			c.checkBlock(ctx, scope, block)

			ch <- Result{
				Node:  program,
				Scope: scope,
			}
		}

		ctx.Logger().Infof("type checking succeeded")
	}()

	return ch
}

func (c TypeChecker) checkBlock(
	ctx interface {
		context.LoggerContext
		context.ErrorsContext
		context.NeutralizerContext
	},
	scope symbol.Scope,
	block ast.Node,
) {
	c.addTypeDecls(ctx, scope, block)
	c.addConstDecls(ctx, scope, block)
	c.addVarDecls(ctx, scope, block)
	c.addFuncDecls(ctx, scope, block)

	// TODO: Find a better way to find main program block.
	operatorsBlocks := block.Query(ast.QueryTypeTop, ast.MarkerOperators)
	mainBlock := operatorsBlocks[len(operatorsBlocks)-1]

	c.checkAssignments(ctx, scope, mainBlock)
	c.checkFlowOperators(ctx, scope, mainBlock)
	c.checkFunctionCalls(ctx, scope, mainBlock)

	ctx.Logger().Infof("type checker found %d symbols in current scope", len(scope.Symbols()))
	for _, s := range scope.Symbols() {
		ctx.Logger().Infof("%d: %v", s.Hash(), s)
	}
}

func (c TypeChecker) checkFunctionCalls(
	ctx interface {
		context.LoggerContext
		context.ErrorsContext
		context.NeutralizerContext
	},
	scope symbol.Scope,
	program ast.Node,
) {
	calls := program.Query(ast.QueryTypeTop, ast.MarkerFuncCall)
	for _, call := range calls {
		name := call.Query(ast.QueryTypeOne, ast.MarkerFuncName)[0].(*ast.Leaf)

		sym, err := ctx.Neutralizer().NeutralizeUserDefined(scope, name.Value)
		if err != nil {
			ctx.AddError(context.ErrorSourceTypecheck, name.Position(), err)
			continue
		}

		args := call.Query(ast.QueryTypeTop, ast.MarkerFuncArg)
		switch sym := sym.(type) {
		case *symbol.Var, *symbol.Const:
			continue
		case *symbol.Func:
			if len(args) != len(sym.Params) {
				ctx.AddError(context.ErrorSourceTypecheck, name.Position(), fmt.Errorf("wrong number of arguments for function %s", name.Value))
				continue
			}

			for i, arg := range args {
				argType, err := c.resolver.Resolve(ctx, scope, arg)
				if err != nil {
					ctx.AddError(context.ErrorSourceTypecheck, arg.Position(), err)
					continue
				}

				if _, ok := c.converter.IsAssignable(argType, sym.Params[i].Type.BuiltinType); ok {
					ctx.AddError(context.ErrorSourceTypecheck, arg.Position(), fmt.Errorf("wrong type of argument %d for function %s", i, name.Value))
					continue
				}
			}
		default:
			panic(fmt.Errorf("unexpected symbol type %T", sym))
		}
	}
}

func (c TypeChecker) checkFlowOperators(
	ctx interface {
		context.LoggerContext
		context.ErrorsContext
		context.NeutralizerContext
	},
	scope symbol.Scope,
	program ast.Node,
) {
	forHeaders := program.Query(ast.QueryTypeTop, ast.MarkerForHeader)
	for _, header := range forHeaders {
		expressions := header.Query(ast.QueryTypeTop, ast.MarkerExpr)
		fromExpr, toExpr := expressions[0], expressions[1]

		fromType, err := c.resolver.Resolve(ctx, scope, fromExpr)
		if err != nil {
			ctx.AddError(context.ErrorSourceTypecheck, fromExpr.Position(), err)
			continue
		}

		toType, err := c.resolver.Resolve(ctx, scope, toExpr)
		if err != nil {
			ctx.AddError(context.ErrorSourceTypecheck, toExpr.Position(), err)
			continue
		}

		if fromType != symbol.BuiltinTypeInt || toType != symbol.BuiltinTypeInt {
			ctx.AddError(context.ErrorSourceTypecheck, fromExpr.Position().Join(toExpr.Position()), fmt.Errorf("range in for loop must have int type"))
			continue
		}
	}

	conditions := program.Query(ast.QueryTypeTop,
		ast.MarkerIfExpr,
		ast.MarkerWhileExpr,
		ast.MarkerRepeatExpr,
	)
	for _, condition := range conditions {
		conditionType, err := c.resolver.Resolve(ctx, scope, condition)
		if err != nil {
			ctx.AddError(context.ErrorSourceTypecheck, condition.Position(), err)
			continue
		}

		if conditionType != symbol.BuiltinTypeBool {
			ctx.AddError(context.ErrorSourceTypecheck, condition.Position(), fmt.Errorf("condition in if/while/repeat statement must have bool type"))
			continue
		}
	}
}

func (c TypeChecker) checkAssignments(
	ctx interface {
		context.LoggerContext
		context.ErrorsContext
		context.NeutralizerContext
	},
	scope symbol.Scope,
	program ast.Node,
) {
	assignments := program.Query(ast.QueryTypeRecursive, ast.MarkerAssign)
	for _, a := range assignments {
		v := a.Query(ast.QueryTypeOne, ast.MarkerLeftSide)[0].(*ast.Leaf)

		vSymbol, err := ctx.Neutralizer().NeutralizeUserDefined(scope, v.Value)
		if err != nil {
			ctx.AddError(context.ErrorSourceTypecheck, v.Position(), err)
			continue
		}

		if _, ok := vSymbol.(*symbol.Var); !ok {
			ctx.AddError(context.ErrorSourceTypecheck, v.Position(), fmt.Errorf("symbol is not a variable: %s", v.Value))
			continue
		}

		expr := a.Query(ast.QueryTypeOne, ast.MarkerRightSide)[0].Query(ast.QueryTypeOne, ast.MarkerExpr)

		exprType, err := c.resolver.Resolve(ctx, scope, expr[0])
		if err != nil {
			ctx.AddError(context.ErrorSourceTypecheck, a.Position(), err)
			continue
		}

		if _, ok := c.converter.IsAssignable(exprType, vSymbol.(*symbol.Var).Type.BuiltinType); !ok {
			ctx.AddError(context.ErrorSourceTypecheck, a.Position(), fmt.Errorf("type mismatch: %s", vSymbol.(*symbol.Var).Type.BuiltinType))
			continue
		}
	}
}

func (c TypeChecker) addFuncDecls(
	ctx interface {
		context.LoggerContext
		context.ErrorsContext
		context.NeutralizerContext
	},
	scope symbol.Scope,
	program ast.Node,
) {
	for _, decl := range program.Query(ast.QueryTypeTop, ast.MarkerFuncDecl) {
		name := decl.Query(ast.QueryTypeOne, ast.MarkerFuncName)[0].(*ast.Leaf)
		returnType := decl.Query(ast.QueryTypeOne, ast.MarkerReturnType)[0].Query(ast.QueryTypeOne, ast.MarkerType)[0].(*ast.Leaf)

		returnTypeSymbol, err := ctx.Neutralizer().NeutralizeUserDefined(scope, returnType.Value)
		if err != nil {
			ctx.AddError(context.ErrorSourceTypecheck, returnType.Position(), err)
			continue
		}

		if _, ok := returnTypeSymbol.(*symbol.Type); !ok {
			ctx.AddError(context.ErrorSourceTypecheck, returnType.Position(), fmt.Errorf("symbol %s is not a type", returnType.Value))
			continue
		}

		var params []symbol.Var
		for _, param := range decl.Query(ast.QueryTypeTop, ast.MarkerParamGroupDecl) {
			paramName := param.Query(ast.QueryTypeOne, ast.MarkerName)[0].(*ast.Leaf)
			paramType := param.Query(ast.QueryTypeOne, ast.MarkerType)[0].(*ast.Leaf)

			paramTypeSymbol, err := ctx.Neutralizer().NeutralizeUserDefined(scope, paramType.Value)
			if err != nil {
				ctx.AddError(context.ErrorSourceTypecheck, paramType.Position(), err)
				continue
			}

			if _, ok := paramTypeSymbol.(*symbol.Type); !ok {
				ctx.AddError(context.ErrorSourceTypecheck, paramType.Position(), fmt.Errorf("symbol %s is not a type", paramType.Value))
				continue
			}

			params = append(params, symbol.Var{
				Token:       paramName.Token,
				Type:        *paramTypeSymbol.(*symbol.Type),
				Initialized: false,
			})
		}

		functionSymbol := &symbol.Func{
			Token:      name.Token,
			Params:     params,
			ReturnType: *returnTypeSymbol.(*symbol.Type),
		}

		if err := scope.Add(functionSymbol); err != nil {
			ctx.AddError(context.ErrorSourceTypecheck, name.Position(), err)
			continue
		}

		var functionSymbols []symbol.Symbol
		for _, param := range params {
			functionSymbols = append(functionSymbols, &param)
		}

		functionSymbols = append(functionSymbols, &symbol.Var{
			Token:       functionSymbol.Token,
			Type:        functionSymbol.ReturnType,
			Initialized: false,
		})

		functionScope := scope.SubScope(functionSymbols)
		functionBlock := decl.Query(ast.QueryTypeOne, ast.MarkerFunctionBlock)[0]

		c.checkBlock(ctx, functionScope, functionBlock)
	}
}

func (c TypeChecker) addTypeDecls(
	ctx interface {
		context.LoggerContext
		context.ErrorsContext
		context.NeutralizerContext
	},
	scope symbol.Scope,
	program ast.Node,
) {
	for _, decl := range program.Query(ast.QueryTypeTop, ast.MarkerTypeDecl) {
		name := decl.Query(ast.QueryTypeOne, ast.MarkerName)[0].(*ast.Leaf)
		typeName := decl.Query(ast.QueryTypeOne, ast.MarkerType)[0].(*ast.Leaf)

		typeNameSymbol, err := ctx.Neutralizer().NeutralizeUserDefined(scope, typeName.Value)
		if err != nil {
			ctx.AddError(context.ErrorSourceTypecheck, typeName.Position(), err)
			continue
		}

		if _, ok := typeNameSymbol.(*symbol.Type); !ok {
			ctx.AddError(context.ErrorSourceTypecheck, typeName.Position(), fmt.Errorf("symbol %s is not a type", typeName.Value))
			continue
		}

		if err := scope.Add(&symbol.Type{
			Token:       name.Token,
			Alias:       typeNameSymbol.(*symbol.Type),
			BuiltinType: typeNameSymbol.(*symbol.Type).BuiltinType,
		}); err != nil {
			ctx.AddError(context.ErrorSourceTypecheck, name.Position(), err)
			continue
		}
	}
}

func (c TypeChecker) addConstDecls(
	ctx interface {
		context.LoggerContext
		context.ErrorsContext
		context.NeutralizerContext
	},
	scope symbol.Scope,
	program ast.Node,
) {
	for _, decl := range program.Query(ast.QueryTypeTop, ast.MarkerConstDecl) {
		name := decl.Query(ast.QueryTypeOne, ast.MarkerName)[0].(*ast.Leaf)
		valueNode := decl.Query(ast.QueryTypeOne, ast.MarkerValue)[0].(*ast.Leaf)

		var typeName string
		switch valueNode.ID {
		case token.IntLiteral:
			typeName = "integer"
		case token.DoubleLiteral:
			typeName = "real"
		case token.BoolLiteral:
			typeName = "boolean"
		default:
			ctx.AddError(context.ErrorSourceTypecheck, valueNode.Position(), fmt.Errorf("unsupported constant type %s", valueNode.ID))
			continue
		}

		typeSymbol, err := ctx.Neutralizer().NeutralizeUserDefined(scope, typeName)
		if err != nil {
			ctx.AddError(context.ErrorSourceTypecheck, valueNode.Position(), err)
			continue
		}

		if _, ok := typeSymbol.(*symbol.Type); !ok {
			ctx.AddError(context.ErrorSourceTypecheck, valueNode.Position(), fmt.Errorf("symbol %s is not a type", typeName))
			continue
		}

		if err := scope.Add(&symbol.Const{
			Token:    name.Token,
			Type:     *typeSymbol.(*symbol.Type),
			RawValue: valueNode.Value,
		}); err != nil {
			ctx.AddError(context.ErrorSourceTypecheck, name.Position(), err)
			continue
		}
	}
}

func (c TypeChecker) addVarDecls(
	ctx interface {
		context.LoggerContext
		context.ErrorsContext
		context.NeutralizerContext
	},
	scope symbol.Scope,
	program ast.Node,
) {
	for _, decl := range program.Query(ast.QueryTypeTop, ast.MarkerVarDecl) {
		typeName := decl.Query(ast.QueryTypeOne, ast.MarkerType)[0].(*ast.Leaf)

		typeNameSymbol, err := ctx.Neutralizer().NeutralizeUserDefined(scope, typeName.Value)
		if err != nil {
			ctx.AddError(context.ErrorSourceTypecheck, typeName.Position(), err)
			continue
		}

		if _, ok := typeNameSymbol.(*symbol.Type); !ok {
			ctx.AddError(context.ErrorSourceTypecheck, typeName.Position(), fmt.Errorf("symbol %s is not a type", typeName.Value))
			continue
		}

		for _, name := range decl.Query(ast.QueryTypeTop, ast.MarkerName) {
			name := name.(*ast.Leaf)

			if err := scope.Add(&symbol.Var{
				Token: name.Token,
				Type:  *typeNameSymbol.(*symbol.Type),
			}); err != nil {
				ctx.AddError(context.ErrorSourceTypecheck, name.Position(), err)
				continue
			}
		}
	}
}
