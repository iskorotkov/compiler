package wasm

import (
	"fmt"
	"os/exec"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/ast"
	"github.com/iskorotkov/compiler/internal/data/symbol"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
	"github.com/iskorotkov/compiler/internal/module/typechecker"
)

var opsByPriority = []map[token.ID]bool{
	{
		token.Xor: true,
		token.Or:  true,
	},
	{
		token.And: true,
	},
	{
		token.Eq:  true,
		token.Ne:  true,
		token.Lt:  true,
		token.Lte: true,
		token.Gt:  true,
		token.Gte: true,
	},
	{
		token.Plus:  true,
		token.Minus: true,
	},
	{
		token.Multiply: true,
		token.Divide:   true,
		token.Mod:      true,
	},
}

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(
	ctx interface {
		context.LoggerContext
	},
	input <-chan option.Option[typechecker.Result],
) chan option.Option[interface{}] {
	ch := make(chan option.Option[interface{}])

	go func() {
		defer close(ch)

		for opt := range input {
			program, err := opt.Unwrap()
			if err != nil {
				ch <- option.Err[interface{}](err)
				return
			}

			var globals []Global
			for _, s := range program.Scope.Symbols() {
				switch s := s.(type) {
				case *symbol.Type:
					// We don't declare typedefs in generated code.
					continue
				case *symbol.Const:
					globals = append(globals, Global{
						Name:    s.Value,
						Type:    g.convertToWASMType(s.Type.BuiltinType),
						Value:   s.RawValue,
						Mutable: false,
					})
				case *symbol.Func:
					ctx.Logger().Errorf("%s: function declaration in wasm module is not supported for now", s.Value)
					continue
				case *symbol.Var:
					globals = append(globals, Global{
						Name:    s.Value,
						Type:    g.convertToWASMType(s.Type.BuiltinType),
						Value:   g.defaultValue(s.Type.BuiltinType),
						Mutable: true,
					})
				default:
					panic("unknown symbol type")
				}
			}

			m := Module{
				Imports: []Import{
					// Add writeln function that invokes console.log imported from JS.
					// We use f64 as a param type because we don't need string support,
					// and i32 can be converted to f64 without loss of precision.
					{
						Path: []string{"console", "log"},
						Name: "writeln_i32",
						Params: []Param{
							{
								Name: "value",
								Type: TypeI32,
							},
						},
						Return: nil,
					},
					{
						Path: []string{"console", "log"},
						Name: "writeln_f64",
						Params: []Param{
							{
								Name: "value",
								Type: TypeF64,
							},
						},
						Return: nil,
					},
				},
				Globals: globals,
				Funcs: []Func{
					// Main function (program execution starts here).
					{
						Name:   "main",
						Params: nil,
						Return: nil,
						Body:   g.buildFuncBody(ctx, program.Scope, program.Node),
					},
					// TODO: Add support for other functions
				},
			}

			s := m.String()
			fmt.Println(s)
		}
	}()

	return ch
}

func (g *Generator) buildFuncBody(
	ctx interface {
		context.LoggerContext
	},
	scope symbol.Scope,
	node ast.Node,
) []Statement {
	nodes := node.Query(ast.QueryTypeTop,
		ast.MarkerAssign,
		ast.MarkerIf,
		ast.MarkerFor,
		ast.MarkerWhile,
		ast.MarkerRepeat,
		ast.MarkerFuncCall,
	)

	var statements []Statement
	for _, node := range nodes {
		switch {
		case node.Has(ast.MarkerAssign):
			variable := node.
				Query(ast.QueryTypeOne, ast.MarkerLeftSide)[0].(*ast.Leaf)

			expr := node.
				Query(ast.QueryTypeOne, ast.MarkerRightSide)[0].
				Query(ast.QueryTypeOne, ast.MarkerExpr)[0]

			wasmExpr, _, err := g.buildExpression(scope, expr)
			if err != nil {
				ctx.Logger().Errorf("%s: %s", variable.Value, err)
				continue
			}

			statements = append(statements, &GlobalSet{
				Name: variable.Value,
				Expr: wasmExpr,
			})
		case node.Has(ast.MarkerIf):
			expr := node.
				Query(ast.QueryTypeOne, ast.MarkerExpr)[0]

			wasmExpr, _, err := g.buildExpression(scope, expr)
			if err != nil {
				ctx.Logger().Errorf("%v: %s", expr, err)
				continue
			}

			operators := node.
				Query(ast.QueryTypeTop, ast.MarkerBlock)

			var falseBody []Statement
			if len(operators) > 1 {
				falseBody = g.buildFuncBody(ctx, scope, operators[1])
			}

			statements = append(statements, &If{
				Cond:      wasmExpr,
				TrueBody:  g.buildFuncBody(ctx, scope, operators[0]),
				FalseBody: falseBody,
			})
		case node.Has(ast.MarkerFor):
			// TODO: Add support for for loops.
			panic("for loop is not supported yet")
		case node.Has(ast.MarkerWhile):
			expr := node.
				Query(ast.QueryTypeOne, ast.MarkerExpr)[0]

			wasmExpr, _, err := g.buildExpression(scope, expr)
			if err != nil {
				ctx.Logger().Errorf("%v: %s", expr, err)
				continue
			}

			body := node.
				Query(ast.QueryTypeOne, ast.MarkerBlock)[0]

			statements = append(statements, &Loop{
				PreCond: &BinaryOp{
					Type: TypeI32,
					Op:   OpEq,
					Left: wasmExpr,
					Right: &Const{
						Type:  TypeI32,
						Value: "0",
					},
				},
				Body: g.buildFuncBody(ctx, scope, body),
			})
		case node.Has(ast.MarkerRepeat):
			// TODO: Add support for repeat loops.
			panic("repeat loop is not supported yet")
		case node.Has(ast.MarkerFuncCall):
			name := node.Query(ast.QueryTypeOne, ast.MarkerName)[0].(*ast.Leaf)
			if name.Value != "writeln" {
				panic("only writeln function is supported")
			}

			args := node.Query(ast.QueryTypeTop, ast.MarkerFuncArg)
			if len(args) != 1 {
				panic("only one argument is supported")
			}

			argExpr, builtinType, err := g.buildExpression(scope, args[0])
			if err != nil {
				ctx.Logger().Errorf("%v: %s", args[0], err)
				continue
			}

			funcName := fmt.Sprintf("%s_%s", name.Value, g.convertToWASMType(builtinType))

			statements = append(statements, &FuncCall{
				Name: funcName,
				Args: []Expr{argExpr},
			})
		default:
			panic("unknown node marker")
		}
	}

	return statements
}

type exprTree struct {
	builtinType symbol.BuiltinType
	node        *ast.Leaf
	left, right *exprTree
}

func (e exprTree) Leaf() bool {
	return e.left == nil && e.right == nil
}

func (g *Generator) buildExpression(scope symbol.Scope, node ast.Node) (Expr, symbol.BuiltinType, error) {
	leafs := g.linearizeExpression(node)
	tree := g.buildExprTree(scope, leafs)
	expr, err := g.convertToWASMExpr(scope, *tree)
	return expr, tree.builtinType, err
}

func (g *Generator) convertToWASMExpr(scope symbol.Scope, tree exprTree) (Expr, error) {
	if tree.Leaf() {
		if tree.node.ID == token.UserDefined {
			// TODO: Add support for function params and locals.
			return &GlobalGet{
				Name: tree.node.Value,
			}, nil
		}

		switch tree.builtinType {
		case symbol.BuiltinTypeBool:
			value := 0
			if tree.node.Value == "true" {
				value = 1
			}

			return &Const{
				Type:  TypeI32,
				Value: fmt.Sprint(value),
			}, nil
		case symbol.BuiltinTypeInt:
			return &Const{
				Type:  TypeI32,
				Value: tree.node.Value,
			}, nil
		case symbol.BuiltinTypeDouble:
			return &Const{
				Type:  TypeF64,
				Value: tree.node.Value,
			}, nil
		default:
			panic("unexpected token id")
		}
	}

	left, err := g.convertToWASMExpr(scope, *tree.left)
	if err != nil {
		return nil, err
	}

	right, err := g.convertToWASMExpr(scope, *tree.right)
	if err != nil {
		return nil, err
	}

	left, right, exprType := g.wrapWithTypeConversions(left, right, tree)

	op, err := MapTokenToWASMOp(tree.node.ID, exprType)
	if err != nil {
		return nil, err
	}

	return &BinaryOp{
		Type:  g.convertToWASMType(exprType),
		Op:    op,
		Left:  left,
		Right: right,
	}, nil
}

func (g *Generator) wrapWithTypeConversions(left, right Expr, tree exprTree) (Expr, Expr, symbol.BuiltinType) {
	if tree.left.builtinType != symbol.BuiltinTypeDouble && tree.right.builtinType == symbol.BuiltinTypeDouble {
		return &Conversion{
			ResultingType: TypeF64,
			Expr:          left,
		}, right, symbol.BuiltinTypeDouble
	}

	if tree.left.builtinType == symbol.BuiltinTypeDouble && tree.right.builtinType != symbol.BuiltinTypeDouble {
		return left, &Conversion{
			ResultingType: TypeF64,
			Expr:          right,
		}, symbol.BuiltinTypeDouble
	}

	return left, right, tree.left.builtinType
}

func (g *Generator) buildExprTree(scope symbol.Scope, leafs []*ast.Leaf) *exprTree {
	if len(leafs) == 1 {
		var builtinType symbol.BuiltinType
		switch leafs[0].ID {
		case token.BoolLiteral:
			builtinType = symbol.BuiltinTypeBool
		case token.IntLiteral:
			builtinType = symbol.BuiltinTypeInt
		case token.DoubleLiteral:
			builtinType = symbol.BuiltinTypeDouble
		case token.UserDefined:
			s, ok := scope.Lookup(&symbol.Name{Name: leafs[0].Value})
			if !ok {
				panic(fmt.Sprintf("%s: symbol not found", leafs[0].Value))
			}

			switch s := s.(type) {
			case *symbol.Const:
				builtinType = s.Type.BuiltinType
			case *symbol.Var:
				builtinType = s.Type.BuiltinType
			case *symbol.Func:
				panic("function call not implemented yet")
			default:
				panic("unknown symbol type")
			}
		default:
			panic(fmt.Sprintf("unexpected token id: %s", leafs[0].ID))
		}

		return &exprTree{
			builtinType: builtinType,
			node:        leafs[0],
		}
	}

	for _, ops := range opsByPriority {
		for i, leaf := range leafs {
			if !ops[leaf.ID] {
				continue
			}

			left := g.buildExprTree(scope, leafs[:i])
			right := g.buildExprTree(scope, leafs[i+1:])

			builtinType := left.builtinType
			if left.builtinType == symbol.BuiltinTypeDouble || right.builtinType == symbol.BuiltinTypeDouble {
				builtinType = symbol.BuiltinTypeDouble
			}

			return &exprTree{
				builtinType: builtinType,
				node:        leaf,
				left:        left,
				right:       right,
			}
		}
	}

	panic("found several consecutive non-operator leafs in expressions")
}

func (g *Generator) linearizeExpression(node ast.Node) []*ast.Leaf {
	switch node := node.(type) {
	case *ast.Leaf:
		return []*ast.Leaf{node}
	case *ast.Branch:
		var res []*ast.Leaf
		for _, items := range node.Items {
			res = append(res, g.linearizeExpression(items)...)
		}

		return res
	default:
		panic("unknown node type")
	}
}

func (g *Generator) convertToWASMType(t symbol.BuiltinType) Type {
	switch t {
	case symbol.BuiltinTypeInt:
		return TypeI32
	case symbol.BuiltinTypeDouble:
		return TypeF64
	case symbol.BuiltinTypeString:
		panic("string usage in WASM is not supported for now")
	case symbol.BuiltinTypeBool:
		return TypeI32
	default:
		panic("unknown type")
	}
}

func (g *Generator) defaultValue(t symbol.BuiltinType) string {
	switch t {
	case symbol.BuiltinTypeInt:
		return "0"
	case symbol.BuiltinTypeDouble:
		return "0.0"
	case symbol.BuiltinTypeString:
		panic("strings are not supported")
	case symbol.BuiltinTypeBool:
		return "0"
	default:
		panic("unknown type")
	}
}

func (g *Generator) wat2Wasm(watFilename string) (string, error) {
	file, err := exec.LookPath("wat2wasm")
	if err != nil {
		return "", err
	}

	cmd := exec.Command(file, watFilename)
	if err := cmd.Run(); err != nil {
		b, _ := cmd.CombinedOutput()
		return string(b), err
	}

	return "", nil
}
