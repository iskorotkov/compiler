package wasm

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/data/symbol"
	"github.com/iskorotkov/compiler/internal/data/token"
)

const (
	// Math.

	OpAdd Op = "add"
	OpSub Op = "sub"
	OpMul Op = "mul"
	OpDiv Op = "div_s"

	// Logical.

	OpAnd Op = "and"
	OpXor Op = "xor"

	// Comparison.

	OpEq Op = "eq"
	OpNe Op = "ne"

	// Comparison (float only).

	OpLt Op = "lt"
	OpGt Op = "gt"
	OpLe Op = "le"
	OpGe Op = "ge"

	// Comparison (int only).

	OpLtSigned Op = "lt_s"
	OpGtSigned Op = "gt_s"
	OpLeSigned Op = "le_s"
	OpGeSigned Op = "ge_s"

	// Conversion.

	OpTruncateF64 Op = "trunc_f64_s"
	OpConvertI32  Op = "convert_i32_s"
)

const (
	TypeI32 Type = "i32"
	TypeF64 Type = "f64"
)

var (
	tokensToWASMOps = map[token.ID]Op{
		token.Plus:     OpAdd,
		token.Minus:    OpSub,
		token.Multiply: OpMul,
		token.Divide:   OpDiv,
		token.And:      OpAnd,
		token.Xor:      OpXor,
		token.Eq:       OpEq,
		token.Ne:       OpNe,
	}
	tokensToWASMOpsInt = map[token.ID]Op{
		token.Lt:  OpLtSigned,
		token.Gt:  OpGtSigned,
		token.Lte: OpLeSigned,
		token.Gte: OpGeSigned,
	}
	tokensToWASMOpsFloat = map[token.ID]Op{
		token.Lt:  OpLt,
		token.Gt:  OpGt,
		token.Lte: OpLe,
		token.Gte: OpGe,
	}
)

func MapTokenToWASMOp(t token.ID, builtinType symbol.BuiltinType) (Op, error) {
	if op, ok := tokensToWASMOps[t]; ok {
		return op, nil
	}

	switch builtinType {
	case symbol.BuiltinTypeInt, symbol.BuiltinTypeBool:
		if op, ok := tokensToWASMOpsInt[t]; ok {
			return op, nil
		}
	case symbol.BuiltinTypeDouble:
		if op, ok := tokensToWASMOpsFloat[t]; ok {
			return op, nil
		}
	default:
		panic(fmt.Sprintf("unsupported builtin type: %v", builtinType))
	}

	return "", fmt.Errorf("unsupported token %s", t)
}

type Op string

type Type string

type Expr interface {
	String() string
}

type BinaryOp struct {
	Type  Type
	Op    Op
	Left  Expr
	Right Expr
}

func (o *BinaryOp) String() string {
	return fmt.Sprintf("(%s.%s %s %s)",
		o.Type, o.Op, o.Left, o.Right)
}

type Conversion struct {
	ResultingType Type
	Expr          Expr
}

func (o *Conversion) String() string {
	switch o.ResultingType {
	case TypeI32:
		return fmt.Sprintf("(i32.%s %v)", OpTruncateF64, o.Expr)
	case TypeF64:
		return fmt.Sprintf("(f64.%s %v)", OpConvertI32, o.Expr)
	default:
		panic(fmt.Sprintf("unsupported conversion type: %v", o.ResultingType))
	}
}

type Const struct {
	Type  Type
	Value string
}

func (c *Const) String() string {
	return fmt.Sprintf("(%s.const %s)", c.Type, c.Value)
}

type LocalGet struct {
	Name string
}

func (l *LocalGet) String() string {
	return fmt.Sprintf("(local.get $%s)", l.Name)
}

type GlobalGet struct {
	Name string
}

func (g *GlobalGet) String() string {
	return fmt.Sprintf("(global.get $%s)", g.Name)
}
