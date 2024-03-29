package bnf

import (
	"github.com/iskorotkov/compiler/internal/data/ast"
	"github.com/iskorotkov/compiler/internal/data/token"
)

// Conditions.

var (
	If Sequence
)

func init() {
	If = Sequence{Name: "if", BNFs: []BNF{
		Token{ID: token.If},
		Sequence{BNFs: []BNF{&Expression}, Markers: ast.Markers{ast.MarkerIfExpr: true}},
		Token{ID: token.Then},
		&Block,
		Optional{BNF: Sequence{BNFs: []BNF{
			Token{ID: token.Else},
			&Block,
		}}},
	}, Markers: ast.Markers{ast.MarkerIf: true}}
}

// Constants.

var (
	IntLiteral         Token
	DoubleLiteral      Token
	BoolLiteral        Token
	Constant           Either
	ConstantDefinition Sequence
	Constants          Optional
	Sign               Optional
)

func init() {
	IntLiteral = Token{ID: token.IntLiteral}
	DoubleLiteral = Token{ID: token.DoubleLiteral}
	BoolLiteral = Token{ID: token.BoolLiteral}

	Sign = Optional{Name: "sign", BNF: Either{BNFs: []BNF{
		Token{ID: token.Plus},
		Token{ID: token.Minus},
	}}}

	Constant = Either{Name: "constant", BNFs: []BNF{
		Sequence{BNFs: []BNF{
			&Sign,
			Either{BNFs: []BNF{
				Either{BNFs: []BNF{
					&IntLiteral,
					&DoubleLiteral,
				}},
				Token{ID: token.UserDefined},
			}},
		}},
		&BoolLiteral,
	}, Markers: ast.Markers{ast.MarkerValue: true}}

	ConstantDefinition = Sequence{Name: "constant-definition", BNFs: []BNF{
		Token{ID: token.UserDefined, Markers: ast.Markers{ast.MarkerName: true}},
		Token{ID: token.Eq},
		&Constant,
	}, Markers: ast.Markers{ast.MarkerConstDecl: true}}

	Constants = Optional{Name: "constants", BNF: Sequence{BNFs: []BNF{
		Token{ID: token.Const},
		&ConstantDefinition,
		Token{ID: token.Semicolon},
		Several{BNF: Sequence{BNFs: []BNF{
			&ConstantDefinition,
			Token{ID: token.Semicolon},
		}}},
	}}}
}

// Expressions.

var (
	Expression              Sequence
	SimpleExpression        Sequence
	RelationOperation       Either
	AdditiveOperation       Either
	AdditiveOperand         Sequence
	MultiplicativeOperation Either
	MultiplicativeOperand   Either
)

func init() {
	Expression = Sequence{Name: "expression", BNFs: []BNF{
		&SimpleExpression,
		Optional{BNF: Sequence{BNFs: []BNF{
			&RelationOperation,
			&SimpleExpression,
		}}},
	}, Markers: ast.Markers{ast.MarkerExpr: true}}

	SimpleExpression = Sequence{Name: "simple-expression", BNFs: []BNF{
		&Sign,
		&AdditiveOperand,
		Several{BNF: Sequence{BNFs: []BNF{
			&AdditiveOperation,
			&AdditiveOperand,
		}}, Markers: ast.Markers{ast.MarkerAdditionalOperands: true}},
	}}

	RelationOperation = Either{Name: "relation-operation", BNFs: []BNF{
		Token{ID: token.Eq},
		Token{ID: token.Ne},
		Token{ID: token.Lt},
		Token{ID: token.Lte},
		Token{ID: token.Gt},
		Token{ID: token.Gte},
		Token{ID: token.In},
	}, Markers: ast.Markers{ast.MarkerCompareOp: true}}

	AdditiveOperation = Either{Name: "additive-operation", BNFs: []BNF{
		Token{ID: token.Plus, Markers: ast.Markers{ast.MarkerAdditiveOp: true}},
		Token{ID: token.Minus, Markers: ast.Markers{ast.MarkerAdditiveOp: true}},
		Token{ID: token.Or, Markers: ast.Markers{ast.MarkerLogicOp: true}},
	}}

	AdditiveOperand = Sequence{Name: "additive-operand", BNFs: []BNF{
		&MultiplicativeOperand,
		Several{BNF: Sequence{BNFs: []BNF{
			&MultiplicativeOperation,
			&MultiplicativeOperand,
		}}, Markers: ast.Markers{ast.MarkerAdditionalOperands: true}},
	}}

	MultiplicativeOperation = Either{Name: "multiplicative-operation", BNFs: []BNF{
		Token{ID: token.Multiply, Markers: ast.Markers{ast.MarkerMultiplicativeOp: true}},
		Token{ID: token.Divide, Markers: ast.Markers{ast.MarkerMultiplicativeOp: true}},
		Token{ID: token.Div, Markers: ast.Markers{ast.MarkerMultiplicativeOp: true}},
		Token{ID: token.Mod, Markers: ast.Markers{ast.MarkerMultiplicativeOp: true}},
		Token{ID: token.And, Markers: ast.Markers{ast.MarkerLogicOp: true}},
	}}

	MultiplicativeOperand = Either{Name: "multiplicative-operand", BNFs: []BNF{
		&Variable,
		&Constant,
		&FunctionCall,
		Sequence{BNFs: []BNF{
			Token{ID: token.OpeningParenthesis},
			&Expression,
			Token{ID: token.ClosingParenthesis},
		}},
		Sequence{BNFs: []BNF{
			Token{ID: token.Not},
			&MultiplicativeOperand,
		}},
	}}
}

// Functions.

var (
	FunctionName       Token
	ParameterGroup     Sequence
	FormalParameters   Either
	FactualParameter   Either
	FunctionHeader     Sequence
	FunctionDefinition Sequence
	Functions          Several
	FunctionCall       Sequence
	FunctionReturnType Sequence
)

func init() {
	Functions = Several{Name: "functions", BNF: &FunctionDefinition}
	FunctionReturnType = Sequence{Name: "function-return-type", BNFs: []BNF{&Type}, Markers: ast.Markers{ast.MarkerReturnType: true}}

	FunctionName = Token{ID: token.UserDefined, Markers: ast.Markers{
		ast.MarkerName:     true,
		ast.MarkerFuncName: true,
	}}

	ParameterGroup = Sequence{Name: "parameter-group", BNFs: []BNF{
		Token{ID: token.UserDefined, Markers: ast.Markers{ast.MarkerName: true}},
		Several{BNF: Sequence{BNFs: []BNF{
			Token{ID: token.Comma},
			Token{ID: token.UserDefined, Markers: ast.Markers{ast.MarkerName: true}},
		}}},
		Token{ID: token.Colon},
		&Type,
	}, Markers: ast.Markers{ast.MarkerParamGroupDecl: true}}

	FormalParameters = Either{Name: "formal-parameters", BNFs: []BNF{
		&ParameterGroup,
		Sequence{BNFs: []BNF{
			Token{ID: token.Var},
			&ParameterGroup,
		}},
		Sequence{BNFs: []BNF{
			Token{ID: token.Function},
			&ParameterGroup,
		}},
	}}

	FactualParameter = Either{Name: "factual-parameter", BNFs: []BNF{
		&Expression,
		&Variable,
		&FunctionName,
	}, Markers: ast.Markers{ast.MarkerFuncArg: true}}

	FunctionHeader = Sequence{Name: "function-header", BNFs: []BNF{
		Token{ID: token.Function},
		&FunctionName,
		Optional{BNF: Sequence{BNFs: []BNF{
			Token{ID: token.OpeningParenthesis},
			Optional{BNF: Sequence{BNFs: []BNF{
				&FormalParameters,
				Several{BNF: Sequence{BNFs: []BNF{
					Token{ID: token.Semicolon},
					&FormalParameters,
				}}},
			}}},
			Token{ID: token.ClosingParenthesis},
		}}},
		Token{ID: token.Colon},
		&FunctionReturnType,
		Token{ID: token.Semicolon},
	}}

	FunctionDefinition = Sequence{Name: "function-definition", BNFs: []BNF{
		&FunctionHeader,
		&FunctionBlock,
	}, Markers: ast.Markers{ast.MarkerFuncDecl: true}}

	FunctionCall = Sequence{Name: "function-call", BNFs: []BNF{
		&FunctionName,
		Token{ID: token.OpeningParenthesis},
		Optional{BNF: Sequence{BNFs: []BNF{
			&FactualParameter,
			Several{BNF: Sequence{BNFs: []BNF{
				Token{ID: token.Comma},
				&FactualParameter,
			}}},
		}}},
		Token{ID: token.ClosingParenthesis},
	}, Markers: ast.Markers{ast.MarkerFuncCall: true}}
}

// Loops.

var (
	Repeat    Sequence
	Direction Either
	For       Sequence
	While     Sequence
)

func init() {
	While = Sequence{Name: "while", BNFs: []BNF{
		Token{ID: token.While},
		Sequence{BNFs: []BNF{&Expression}, Markers: ast.Markers{ast.MarkerWhileExpr: true}},
		Token{ID: token.Do},
		&Block,
	}, Markers: ast.Markers{ast.MarkerWhile: true}}

	Repeat = Sequence{Name: "repeat", BNFs: []BNF{
		Token{ID: token.Repeat},
		&Block,
		Token{ID: token.Until},
		Sequence{BNFs: []BNF{&Expression}, Markers: ast.Markers{ast.MarkerRepeatExpr: true}},
	}, Markers: ast.Markers{ast.MarkerRepeat: true}}

	Direction = Either{Name: "direction", BNFs: []BNF{
		Token{ID: token.To},
		Token{ID: token.Downto},
	}}

	For = Sequence{Name: "for", BNFs: []BNF{
		Sequence{BNFs: []BNF{
			Token{ID: token.For},
			Token{ID: token.UserDefined},
			Token{ID: token.Assign},
			&Expression,
			&Direction,
			&Expression,
			Token{ID: token.Do},
		}, Markers: ast.Markers{ast.MarkerForHeader: true}},
		&Block,
	}, Markers: ast.Markers{ast.MarkerFor: true}}
}

// Operators.

var (
	Operator           Optional
	Block              Sequence
	SimpleOperator     Either
	CompositeOperator  Sequence
	ComplexOperator    Either
	ConditionOperator  Either
	LoopOperator       Either
	Operators          Sequence
	AssignmentOperator Sequence
)

func init() {
	ConditionOperator = Either{Name: "condition-operator", BNFs: []BNF{&If}}
	Operators = Sequence{Name: "operators", BNFs: []BNF{&CompositeOperator}}
	Block = Sequence{Name: "block", BNFs: []BNF{&Operator}, Markers: ast.Markers{ast.MarkerBlock: true}}

	SimpleOperator = Either{Name: "simple-operator", BNFs: []BNF{
		&FunctionCall,
		&AssignmentOperator,
	}}

	// TODO: Syntax analyzer is very sensitive to extra semicolons.
	Operator = Optional{Name: "operator", BNF: Either{BNFs: []BNF{
		&SimpleOperator,
		&ComplexOperator,
	}}}

	CompositeOperator = Sequence{Name: "composite-operator", BNFs: []BNF{
		Token{ID: token.Begin},
		&Operator,
		Several{BNF: Sequence{BNFs: []BNF{
			Token{ID: token.Semicolon},
			&Operator,
		}}},
		Token{ID: token.End},
	}}

	ComplexOperator = Either{Name: "complex-operator", BNFs: []BNF{
		&CompositeOperator,
		&AssignmentOperator,
		&ConditionOperator,
		&LoopOperator,
	}}

	LoopOperator = Either{Name: "loop-operator", BNFs: []BNF{
		&For,
		&While,
		&Repeat,
	}}

	AssignmentOperator = Sequence{Name: "assignment-operator", BNFs: []BNF{
		Either{BNFs: []BNF{
			&Variable,
			&FunctionName,
		}, Markers: ast.Markers{ast.MarkerLeftSide: true}},
		Token{ID: token.Assign},
		Sequence{BNFs: []BNF{&Expression}, Markers: ast.Markers{ast.MarkerRightSide: true}},
	}, Markers: ast.Markers{ast.MarkerAssign: true}}
}

// Types.

var (
	TypeDefinition Sequence
	Types          Optional
	Type           Token
)

func init() {
	Type = Token{ID: token.UserDefined, Markers: ast.Markers{ast.MarkerType: true}}

	TypeDefinition = Sequence{Name: "type-definition", BNFs: []BNF{
		Token{ID: token.UserDefined, Markers: ast.Markers{ast.MarkerName: true}},
		Token{ID: token.Eq},
		&Type,
	}, Markers: ast.Markers{ast.MarkerTypeDecl: true}}

	Types = Optional{Name: "types", BNF: Sequence{BNFs: []BNF{
		Token{ID: token.Type},
		&TypeDefinition,
		Token{ID: token.Semicolon},
		Several{BNF: Sequence{BNFs: []BNF{
			&TypeDefinition,
			Token{ID: token.Semicolon},
		}}},
	}}}
}

// Variables.

var (
	VariableName      Sequence
	FullVariable      Sequence
	Variable          Sequence
	SameTypeVariables Sequence
	Variables         Optional
)

func init() {
	VariableName = Sequence{Name: "variable-name", BNFs: []BNF{Token{ID: token.UserDefined}}}
	FullVariable = Sequence{Name: "full variable", BNFs: []BNF{&VariableName}}
	Variable = Sequence{Name: "variable", BNFs: []BNF{&FullVariable}}

	SameTypeVariables = Sequence{Name: "same-type-variables", BNFs: []BNF{
		Token{ID: token.UserDefined, Markers: ast.Markers{ast.MarkerName: true}},
		Several{BNF: Sequence{BNFs: []BNF{
			Token{ID: token.Comma},
			Token{ID: token.UserDefined, Markers: ast.Markers{ast.MarkerName: true}},
		}}},
		Token{ID: token.Colon},
		&Type,
	}, Markers: ast.Markers{ast.MarkerVarDecl: true}}

	Variables = Optional{Name: "variables", BNF: Sequence{BNFs: []BNF{
		Token{ID: token.Var},
		&SameTypeVariables,
		Token{ID: token.Semicolon},
		Several{BNF: Sequence{BNFs: []BNF{
			&SameTypeVariables,
			Token{ID: token.Semicolon},
		}}},
	}}}
}

// Special.

var (
	FunctionBlock Sequence
	Program       Sequence
)

func init() {
	FunctionBlock = Sequence{Name: "function-block", BNFs: []BNF{
		Sequence{BNFs: []BNF{
			&Constants,
			&Types,
			&Variables,
			&Functions,
		}, Markers: ast.Markers{ast.MarkerDeclarations: true}},
		Sequence{BNFs: []BNF{&Operators}, Markers: ast.Markers{ast.MarkerOperators: true}},
	}, Markers: ast.Markers{ast.MarkerFunctionBlock: true}}

	Program = Sequence{Name: "program", BNFs: []BNF{
		Token{ID: token.Program},
		Token{ID: token.UserDefined},
		Token{ID: token.Semicolon},
		Sequence{BNFs: []BNF{&FunctionBlock}, Markers: ast.Markers{ast.MarkerProgramBlock: true}},
		Token{ID: token.Period},
		Token{ID: token.EOF},
	}}
}
