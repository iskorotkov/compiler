package wasm

import (
	"fmt"
	"strings"
)

type Statement interface {
	StringIndent(level int) string
}

type Import struct {
	Path   []string
	Name   string
	Params []Param
	Return *Return
}

func (i *Import) StringIndent(level int) string {
	var pathStrs []string
	for _, p := range i.Path {
		pathStrs = append(pathStrs, fmt.Sprintf("%q", p))
	}

	var paramsStr string
	if len(i.Params) > 0 {
		// We assume there will be no more than 4-5 params,
		// so it's easier and more efficient to use strings instead of string builder.
		paramsStr = " "
		for _, param := range i.Params {
			paramsStr += fmt.Sprintf("(param $%s %s)", param.Name, param.Type)
		}
	}

	if i.Return != nil {
		return fmt.Sprintf("%s(import %s (func $%s%s (return %s)))",
			strings.Repeat("  ", level), strings.Join(pathStrs, " "), i.Name, paramsStr, i.Return.Type)
	} else {
		return fmt.Sprintf("%s(import %s (func $%s%s))",
			strings.Repeat("  ", level), strings.Join(pathStrs, " "), i.Name, paramsStr)
	}
}

type Param struct {
	Name string
	Type Type
}

type Return struct {
	Type Type
}

type Global struct {
	Name    string
	Type    Type
	Value   string
	Mutable bool
}

func (g *Global) StringIndent(level int) string {
	if g.Mutable {
		return fmt.Sprintf("%s(global $%s (mut %[3]s) (%[3]s.const %s))",
			strings.Repeat("  ", level), g.Name, g.Type, g.Value)
	} else {
		return fmt.Sprintf("%s(global $%s %[3]s (%[3]s.const %s))",
			strings.Repeat("  ", level), g.Name, g.Type, g.Value)
	}
}

type Func struct {
	Name   string
	Params []Param
	Return *Return
	Body   []Statement
}

func (f *Func) StringIndent(level int) string {
	var paramsStr string
	if len(f.Params) > 0 {
		paramsStr = " "
		for _, param := range f.Params {
			paramsStr += fmt.Sprintf("(param $%s %s)", param.Name, param.Type)
		}
	}

	var bodyStr string
	for _, stmt := range f.Body {
		bodyStr += stmt.StringIndent(level+1) + "\n"
	}

	if f.Return != nil {
		return fmt.Sprintf(`%s(func (export "%s")%s (return %s)
%s%[1]s)`, strings.Repeat("  ", level), f.Name, paramsStr, f.Return.Type, bodyStr)
	} else {
		return fmt.Sprintf(`%s(func (export "%s")%s
%s%[1]s)`, strings.Repeat("  ", level), f.Name, paramsStr, bodyStr)
	}
}

type If struct {
	Cond      Expr
	TrueBody  []Statement
	FalseBody []Statement
}

func (i *If) StringIndent(level int) string {
	var trueBodyStr string
	for _, stmt := range i.TrueBody {
		trueBodyStr += stmt.StringIndent(level+2) + "\n"
	}

	trueBodyStr = fmt.Sprintf(`
%s(block
%s%[1]s)`, strings.Repeat("  ", level+1), trueBodyStr)

	if len(i.FalseBody) == 0 {
		return fmt.Sprintf(`%s(if %s%s
%[1]s)`,
			strings.Repeat("  ", level), i.Cond, trueBodyStr)
	} else {
		var falseBodyStr string
		for _, stmt := range i.FalseBody {
			falseBodyStr += stmt.StringIndent(level+2) + "\n"
		}

		falseBodyStr = fmt.Sprintf(`%s(block
%s%[1]s)`, strings.Repeat("  ", level+1), falseBodyStr)

		return fmt.Sprintf(`%s(if %s%s
%s
%[1]s)`,
			strings.Repeat("  ", level), i.Cond, trueBodyStr, falseBodyStr)
	}
}

type Loop struct {
	PreCond  Expr
	PostCond Expr
	Body     []Statement
}

func (l *Loop) StringIndent(level int) string {
	var bodyStr string
	for _, stmt := range l.Body {
		bodyStr += stmt.StringIndent(level+2) + "\n"
	}

	if l.PreCond != nil {
		return fmt.Sprintf(`%s(loop
  %[1]s(block
    %[1]s(br_if 0 %[2]s)
%[3]s  %[1]s)
%[1]s)`, strings.Repeat("  ", level), l.PreCond.String(), bodyStr)
	} else {
		return fmt.Sprintf(`%s(loop
  %[1]s(block
%[3]s    %[1]s(br_if 0 %[2]s)
  %[1]s)
%[1]s)`, strings.Repeat("  ", level), l.PostCond.String(), bodyStr)
	}
}

type FuncCall struct {
	Name string
	Args []Expr
}

func (f *FuncCall) StringIndent(level int) string {
	var argsStr string
	if len(f.Args) > 0 {
		argsStr = " "
		for _, arg := range f.Args {
			argsStr += arg.String()
		}
	}

	return fmt.Sprintf("%s(call $%s%s)", strings.Repeat("  ", level), f.Name, argsStr)
}

type LocalSet struct {
	Name string
	Expr Expr
}

func (s *LocalSet) StringIndent(level int) string {
	return fmt.Sprintf("%s(local.set $%s %s)",
		strings.Repeat("  ", level), s.Name, s.Expr.String())
}

type GlobalSet struct {
	Name string
	Expr Expr
}

func (s *GlobalSet) StringIndent(level int) string {
	return fmt.Sprintf("%s(global.set $%s %s)",
		strings.Repeat("  ", level), s.Name, s.Expr.String())
}
