package symbol

import (
	"fmt"

	"github.com/iskorotkov/compiler/internal/data/literal"
	"github.com/iskorotkov/compiler/internal/data/token"
)

const (
	BuiltinTypeUnknown BuiltinType = iota
	BuiltinTypeVoid
	BuiltinTypeInt
	BuiltinTypeDouble
	BuiltinTypeBool
	BuiltinTypeString
)

var (
	_ Symbol = (*Name)(nil)
	_ Symbol = (*Type)(nil)
	_ Symbol = (*Func)(nil)
	_ Symbol = (*Var)(nil)
	_ Symbol = (*Const)(nil)
)

type BuiltinType int

func (t BuiltinType) String() string {
	switch t {
	case BuiltinTypeUnknown:
		return "unknown"
	case BuiltinTypeInt:
		return "int"
	case BuiltinTypeDouble:
		return "double"
	case BuiltinTypeBool:
		return "bool"
	default:
		panic(fmt.Sprintf("unknown builtin type: %d", t))
	}
}

func builtinToken(name string) token.Token {
	return token.Token{
		ID: token.UserDefined,
		Literal: literal.Literal{
			Value:    name,
			Position: literal.Position{},
		},
	}
}

type Symbol interface {
	Hash() int
	fmt.Stringer
}

type Name struct {
	Name string
	hash hasher
}

func (n *Name) Hash() int {
	return n.hash.Hash(n.Name)
}

func (n *Name) String() string {
	return n.Name
}

type Type struct {
	token.Token // Only for user-defined symbols.
	hash        hasher
	Alias       *Type // Only for user-defined types.
	BuiltinType BuiltinType
}

func (t *Type) Hash() int {
	return t.hash.Hash(t.Value)
}

func (t Type) String() string {
	return fmt.Sprintf("type %v", t.Value)
}

type Var struct {
	token.Token // Only for user-defined symbols.
	hash        hasher
	Type        Type
	Initialized bool
}

func (v *Var) Hash() int {
	return v.hash.Hash(v.Value)
}

func (v Var) String() string {
	return fmt.Sprintf("var %v", v.Value)
}

type Const struct {
	token.Token // Only for user-defined symbols.
	hash        hasher
	Type        Type
	RawValue    string
}

func (c *Const) Hash() int {
	return c.hash.Hash(c.Value)
}

func (c Const) String() string {
	return fmt.Sprintf("const %v", c.Value)
}

type Func struct {
	token.Token // Only for user-defined symbols.
	hash        hasher
	Params      []Var
	ReturnType  Type
}

func (f *Func) Hash() int {
	return f.hash.Hash(f.Value)
}

func (f Func) String() string {
	return fmt.Sprintf("func %v", f.Value)
}
