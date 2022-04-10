package bnf

import (
	"github.com/iskorotkov/compiler/internal/data/token"
)

// Conditions.

var (
	If Sequence
)

func init() {
	If = Sequence{Name: "if", BNFs: []BNF{
		Token{ID: token.If},
		&Expression,
		Token{ID: token.Then},
		&Operator,
		Optional{BNF: Sequence{BNFs: []BNF{
			Token{token.Else},
			&Operator,
		}}},
	}}
}

// Constants.

var (
	IntLiteral         Token
	DoubleLiteral      Token
	BoolLiteral        Token
	Constant           Either
	ConstantDefinition Sequence
	Constants          Optional
	Sign               Either
)

func init() {
	Sign = Either{Name: "sign", BNFs: []BNF{
		Token{ID: token.Plus},
		Token{ID: token.Minus},
		&Empty,
	}}

	IntLiteral = Token{ID: token.IntLiteral}

	DoubleLiteral = Token{ID: token.DoubleLiteral}

	BoolLiteral = Token{ID: token.BoolLiteral}

	Constant = Either{Name: "constant", BNFs: []BNF{
		Sequence{BNFs: []BNF{
			Optional{BNF: Sign},
			Either{BNFs: []BNF{
				Either{BNFs: []BNF{
					&IntLiteral,
					&DoubleLiteral,
				}},
				Token{ID: token.UserDefined},
			}},
		}},
		&BoolLiteral,
	}}

	ConstantDefinition = Sequence{Name: "constant-definition", BNFs: []BNF{
		Token{ID: token.UserDefined},
		Token{ID: token.Eq},
		&Constant,
	}}

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
	}}

	SimpleExpression = Sequence{Name: "simple-expression", BNFs: []BNF{
		&Sign,
		&AdditiveOperand,
		Several{BNF: Sequence{BNFs: []BNF{
			&AdditiveOperation,
			&AdditiveOperand,
		}}},
	}}

	RelationOperation = Either{Name: "relation-operation", BNFs: []BNF{
		Token{ID: token.Eq},
		Token{ID: token.Ne},
		Token{ID: token.Lt},
		Token{ID: token.Lte},
		Token{ID: token.Gt},
		Token{ID: token.Gte},
		Token{ID: token.In},
	}}

	AdditiveOperation = Either{Name: "additive-operation", BNFs: []BNF{
		Token{ID: token.Plus},
		Token{ID: token.Minus},
		Token{ID: token.Or},
	}}

	AdditiveOperand = Sequence{Name: "additive-operand", BNFs: []BNF{
		&MultiplicativeOperand,
		Several{BNF: Sequence{"", []BNF{
			&MultiplicativeOperation,
			&MultiplicativeOperand,
		}}},
	}}

	MultiplicativeOperation = Either{Name: "multiplicative-operation", BNFs: []BNF{
		Token{ID: token.Multiply},
		Token{ID: token.Divide},
		Token{ID: token.Div},
		Token{ID: token.Mod},
		Token{ID: token.And},
	}}

	MultiplicativeOperand = Either{Name: "multiplicative-operand", BNFs: []BNF{
		&FunctionUsage,
		&Variable,
		&Constant,
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
	FunctionUsage      Sequence
)

func init() {
	FunctionName = Token{ID: token.UserDefined}

	ParameterGroup = Sequence{Name: "parameter-group", BNFs: []BNF{
		Token{ID: token.UserDefined},
		Several{BNF: Sequence{BNFs: []BNF{
			Token{ID: token.Comma},
			Token{ID: token.UserDefined},
		}}},
		Token{ID: token.Colon},
		&Type,
	}}

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
	}}

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
		&Type,
		Token{ID: token.Semicolon},
	}}

	FunctionDefinition = Sequence{Name: "function-definition", BNFs: []BNF{
		&FunctionHeader,
		&Block,
	}}

	Functions = Several{Name: "functions", BNF: &FunctionDefinition}

	FunctionUsage = Sequence{Name: "function-usage", BNFs: []BNF{
		&FunctionName,
		Optional{BNF: Sequence{BNFs: []BNF{
			Token{ID: token.OpeningParenthesis},
			Optional{BNF: Sequence{BNFs: []BNF{
				&FactualParameter,
				Several{BNF: Sequence{BNFs: []BNF{
					Token{ID: token.Comma},
					&FactualParameter,
				}}},
			}}},
			Token{ID: token.ClosingParenthesis},
		}}},
	}}
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
		&Expression,
		Token{ID: token.Do},
		&Operator,
	}}

	Repeat = Sequence{Name: "repeat", BNFs: []BNF{
		Token{ID: token.Repeat},
		&Operator,
		Several{BNF: Sequence{BNFs: []BNF{
			Token{ID: token.Semicolon},
			&Operator,
		}}},
		Token{ID: token.Until},
		&Expression,
	}}

	Direction = Either{Name: "direction", BNFs: []BNF{
		Token{ID: token.To},
		Token{ID: token.Downto},
	}}

	For = Sequence{Name: "for", BNFs: []BNF{
		Token{ID: token.For},
		Token{ID: token.UserDefined},
		Token{ID: token.Assign},
		&Expression,
		&Direction,
		&Expression,
		Token{ID: token.Do},
		&Operator,
	}}
}

// Operators.

var (
	Operator           Either
	SimpleOperator     Either
	CompositeOperator  Sequence
	ComplexOperator    Either
	ConditionOperator  Either
	LoopOperator       Either
	Operators          Sequence
	AssignmentOperator Sequence
)

func init() {
	// TODO: Syntax analyzer is very sensitive to extra semicolons.
	Operator = Either{Name: "operator", BNFs: []BNF{
		&SimpleOperator,
		&ComplexOperator,
		&Empty,
	}}

	SimpleOperator = Either{Name: "simple-operator", BNFs: []BNF{
		&AssignmentOperator,
	}}

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

	ConditionOperator = Either{Name: "condition-operator", BNFs: []BNF{
		&If,
	}}

	LoopOperator = Either{Name: "loop-operator", BNFs: []BNF{
		&For,
		&While,
		&Repeat,
	}}

	Operators = Sequence{Name: "operators", BNFs: []BNF{
		&CompositeOperator,
	}}

	AssignmentOperator = Sequence{Name: "assignment-operator", BNFs: []BNF{
		Either{BNFs: []BNF{
			&Variable,
			&FunctionName,
		}},
		Token{ID: token.Assign},
		&Expression,
	}}
}

// Types.

var (
	TypeDefinition Sequence
	Types          Optional
	Type           Token
)

func init() {
	Type = Token{ID: token.UserDefined}

	TypeDefinition = Sequence{Name: "type-definition", BNFs: []BNF{
		Token{ID: token.UserDefined},
		Token{ID: token.Eq},
		&Type,
	}}

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
	Variable          Either
	SameTypeVariables Sequence
	Variables         Optional
)

func init() {
	VariableName = Sequence{Name: "variable-name", BNFs: []BNF{
		Token{ID: token.UserDefined},
	}}

	FullVariable = Sequence{Name: "full variable", BNFs: []BNF{
		&VariableName,
	}}

	Variable = Either{Name: "variable", BNFs: []BNF{
		&FullVariable,
	}}

	SameTypeVariables = Sequence{Name: "same-type-variables", BNFs: []BNF{
		Token{ID: token.UserDefined},
		Several{BNF: Sequence{BNFs: []BNF{
			Token{ID: token.Comma},
			Token{ID: token.UserDefined},
		}}},
		Token{ID: token.Colon},
		&Type,
	}}

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
	Block   Sequence
	Program Sequence
	Empty   Sequence
)

func init() {
	Empty = Sequence{}

	Block = Sequence{"block", []BNF{
		&Constants,
		&Types,
		&Variables,
		&Functions,
		&Operators,
	}}

	Program = Sequence{Name: "program", BNFs: []BNF{
		Token{ID: token.Program},
		Token{ID: token.UserDefined},
		Token{ID: token.Semicolon},
		&Block,
		Token{ID: token.Period},
		Token{ID: token.EOF},
	}}
}
