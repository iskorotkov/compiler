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
	MarkerMultiplicativeOp
	MarkerAdditiveOp
	MarkerLogicOp
	MarkerCompareOp

	// Operators.

	MarkerAssign
	MarkerLeftSide
	MarkerRightSide

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

	for k, v := range other {
		ms[k] = v
	}

	return ms
}
