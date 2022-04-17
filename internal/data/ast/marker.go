package ast

const (
	// Declarations.

	MarkerVarDecl Marker = iota
	MarkerConstDecl
	MarkerTypeDecl
	MarkerFuncDecl
	MarkerFuncName
	MarkerParamGroupDecl
	MarkerReturnType

	// Expressions.

	MarkerExpr
	MarkerAdditionalOperands
	MarkerMultiplicativeOp
	MarkerAdditiveOp
	MarkerLogicOp
	MarkerCompareOp

	// Operators.

	MarkerAssign
	MarkerLeftSide
	MarkerRightSide

	// Functions.

	MarkerFuncCall
	MarkerFuncArg

	// Control flow.

	MarkerForHeader
	MarkerIfExpr
	MarkerWhileExpr
	MarkerRepeatExpr

	// Blocks.

	MarkerBlock
	MarkerProgramBlock
	MarkerDeclarations
	MarkerOperators

	// Common.

	MarkerName
	MarkerType
	MarkerValue
)

type Marker int

type Markers map[Marker]bool

func (ms Markers) Has(m Marker) bool {
	if ms == nil {
		return false
	}

	_, ok := ms[m]
	return ok
}

func (ms Markers) Merge(other Markers) Markers {
	if ms == nil {
		return other
	}

	if other == nil {
		return ms
	}

	m := make(Markers)
	for k, v := range ms {
		m[k] = v
	}

	for k, v := range other {
		m[k] = v
	}

	return m
}
