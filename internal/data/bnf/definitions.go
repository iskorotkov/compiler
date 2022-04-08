package bnf

import (
	"github.com/iskorotkov/compiler/internal/data/token"
)

// Conditions.

var (
	If Sequence
)

func init() {
	If = Sequence{"if", []BNF{
		Token{token.If},
		&Expression,
		Token{token.Then},
		&Operator,
		Optional{"", Sequence{"", []BNF{
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
	Sign = Either{"sign", []BNF{
		Token{token.Plus},
		Token{token.Minus},
		&Empty,
	}}

	IntLiteral = Token{token.IntLiteral}

	DoubleLiteral = Token{token.DoubleLiteral}

	BoolLiteral = Token{token.BoolLiteral}

	Constant = Either{"constant", []BNF{
		Sequence{"", []BNF{
			Optional{"", Sign},
			Either{"", []BNF{
				Either{"", []BNF{
					&IntLiteral,
					&DoubleLiteral,
				}},
				Token{token.UserDefined},
			}},
		}},
		&BoolLiteral,
	}}

	ConstantDefinition = Sequence{"constant-definition", []BNF{
		Token{token.UserDefined},
		Token{token.Eq},
		&Constant,
	}}

	Constants = Optional{"constants", Sequence{"", []BNF{
		Token{token.Const},
		&ConstantDefinition,
		Token{token.Semicolon},
		Several{"", Sequence{"", []BNF{
			&ConstantDefinition,
			Token{token.Semicolon},
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
	Expression = Sequence{"expression", []BNF{
		&SimpleExpression,
		Optional{"", Sequence{"", []BNF{
			&RelationOperation,
			&SimpleExpression,
		}}},
	}}

	SimpleExpression = Sequence{"simple-expression", []BNF{
		&Sign,
		&AdditiveOperand,
		Several{"", Sequence{"", []BNF{
			&AdditiveOperation,
			&AdditiveOperand,
		}}},
	}}

	RelationOperation = Either{"relation-operation", []BNF{
		Token{token.Eq},
		Token{token.Ne},
		Token{token.Lt},
		Token{token.Lte},
		Token{token.Gt},
		Token{token.Gte},
		Token{token.In},
	}}

	AdditiveOperation = Either{"additive-operation", []BNF{
		Token{token.Plus},
		Token{token.Minus},
		Token{token.Or},
	}}

	AdditiveOperand = Sequence{"additive-operand", []BNF{
		&MultiplicativeOperand,
		Several{"", Sequence{"", []BNF{
			&MultiplicativeOperation,
			&MultiplicativeOperand,
		}}},
	}}

	MultiplicativeOperation = Either{"multiplicative-operation", []BNF{
		Token{token.Multiply},
		Token{token.Divide},
		Token{token.Div},
		Token{token.Mod},
		Token{token.And},
	}}

	MultiplicativeOperand = Either{"multiplicative-operand", []BNF{
		&Variable,
		&Constant,
		Sequence{"", []BNF{
			Token{token.OpeningParenthesis},
			&Expression,
			Token{token.ClosingParenthesis},
		}},
		&FunctionUsage,
		Sequence{"", []BNF{
			Token{token.Not},
			FunctionUsage,
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
	FunctionName = Token{token.UserDefined}

	ParameterGroup = Sequence{"parameter-group", []BNF{
		Token{token.UserDefined},
		Several{"", Sequence{"", []BNF{
			Token{token.Comma},
			Token{token.UserDefined},
		}}},
		Token{token.Colon},
		&Type,
	}}

	FormalParameters = Either{"formal-parameters", []BNF{
		&ParameterGroup,
		Sequence{"", []BNF{
			Token{token.Var},
			&ParameterGroup,
		}},
		Sequence{"", []BNF{
			Token{token.Function},
			&ParameterGroup,
		}},
	}}

	FactualParameter = Either{"factual-parameter", []BNF{
		&Expression,
		&Variable,
		&FunctionName,
	}}

	FunctionHeader = Sequence{"function-header", []BNF{
		Token{token.Function},
		&FunctionName,
		Optional{"", Sequence{"", []BNF{
			&FormalParameters,
			Several{"", Sequence{"", []BNF{
				Token{token.Colon},
				&FormalParameters,
			}}},
		}}},
		Token{token.Colon},
		&Type,
	}}

	FunctionDefinition = Sequence{"function-definition", []BNF{
		&FunctionHeader,
		&Block,
	}}

	Functions = Several{"functions", &FunctionDefinition}

	FunctionUsage = Sequence{"function-usage", []BNF{
		&FunctionName,
		Optional{"", Sequence{"", []BNF{
			Token{token.OpeningParenthesis},
			&FactualParameter,
			Several{"", Sequence{"", []BNF{
				Token{token.Comma},
				&FactualParameter,
			}}},
			Token{token.ClosingParenthesis},
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
	While = Sequence{"while", []BNF{
		Token{token.While},
		&Expression,
		Token{token.Do},
		&Operator,
	}}

	Repeat = Sequence{"repeat", []BNF{
		Token{token.Repeat},
		&Operator,
		Several{"", Sequence{"", []BNF{
			Token{token.Semicolon},
			&Operator,
		}}},
		Token{token.Until},
		&Expression,
	}}

	Direction = Either{"direction", []BNF{
		Token{token.To},
		Token{token.Downto},
	}}

	For = Sequence{"for", []BNF{
		Token{token.For},
		Token{token.UserDefined},
		Token{token.Assign},
		&Expression,
		&Direction,
		&Expression,
		Token{token.Do},
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
	Operator = Either{"operator", []BNF{
		&SimpleOperator,
		&ComplexOperator,
		&Empty,
	}}

	SimpleOperator = Either{"simple-operator", []BNF{
		&AssignmentOperator,
	}}

	CompositeOperator = Sequence{"composite-operator", []BNF{
		Token{token.Begin},
		&Operator,
		Several{"", Sequence{"", []BNF{
			Token{token.Semicolon},
			&Operator,
		}}},
		Token{token.End},
	}}

	ComplexOperator = Either{"complex-operator", []BNF{
		&CompositeOperator,
		&AssignmentOperator,
		&ConditionOperator,
		&LoopOperator,
	}}

	ConditionOperator = Either{"condition operator", []BNF{
		&If,
	}}

	LoopOperator = Either{"loop operator", []BNF{
		&For,
		&While,
		&Repeat,
	}}

	Operators = Sequence{"operators", []BNF{
		&CompositeOperator,
	}}

	AssignmentOperator = Sequence{"assignment-operator", []BNF{
		Either{"", []BNF{
			&Variable,
			&FunctionName,
		}},
		Token{token.Assign},
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
	Type = Token{token.UserDefined}

	TypeDefinition = Sequence{"type-definition", []BNF{
		Token{token.UserDefined},
		Token{token.Eq},
		&Type,
	}}

	Types = Optional{"types", Sequence{"", []BNF{
		Token{token.Type},
		&TypeDefinition,
		Token{token.Semicolon},
		Several{"", Sequence{"", []BNF{
			&TypeDefinition,
			Token{token.Semicolon},
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
	VariableName = Sequence{"variable-name", []BNF{
		Token{token.UserDefined},
	}}

	FullVariable = Sequence{"full variable", []BNF{
		&VariableName,
	}}

	Variable = Either{"variable", []BNF{
		&FullVariable,
	}}

	SameTypeVariables = Sequence{"same-type-variables", []BNF{
		Token{token.UserDefined},
		Several{"", Sequence{"", []BNF{
			Token{token.Comma},
			Token{token.UserDefined},
		}}},
		Token{token.Colon},
		&Type,
	}}

	Variables = Optional{"variables", Sequence{"", []BNF{
		Token{token.Var},
		&SameTypeVariables,
		Token{token.Semicolon},
		Several{"", Sequence{"", []BNF{
			&SameTypeVariables,
			Token{token.Semicolon},
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

	Program = Sequence{"program", []BNF{
		Token{token.Program},
		Token{token.UserDefined},
		Token{token.Semicolon},
		&Block,
		Token{token.Period},
		Token{token.EOF},
	}}
}
