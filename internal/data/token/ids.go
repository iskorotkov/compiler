package token

type ID int

const (
	Unknown = ID(iota)
	UserDefined
	EOF

	keywordsStart
	Absolute
	Abstract
	All
	And
	AndThen
	Array
	As
	Asm
	Asmname
	Attribute
	Begin
	Bindable
	C
	Case
	CLanguage
	Class
	Const
	Constructor
	Destructor
	Div
	Do
	Downto
	Else
	End
	Export
	Exports
	External
	Far
	File
	Finalization
	For
	Forward
	Function
	Goto
	If
	Implementation
	Import
	In
	Inherited
	Initialization
	Interface
	Interrupt
	Is
	Label
	Library
	Mod
	Module
	Name
	Near
	Nil
	Not
	Object
	Of
	Only
	Operator
	Or
	OrElse
	Otherwise
	Packed
	Pow
	Private
	Procedure
	Program
	Property
	Protected
	Public
	Published
	Qualified
	Record
	Repeat
	Resident
	Restricted
	Segment
	Set
	Shl
	Shr
	Then
	To
	Type
	Unit
	Until
	Uses
	Value
	Var
	View
	Virtual
	While
	With
	Xor
	keywordsEnd

	operatorsStart
	Eq
	Ne
	Lt
	Gt
	Lte
	Gte
	Plus
	Minus
	Multiply
	Divide
	ShiftLeft
	ShiftRight
	Assign
	operatorsEnd

	punctuationStart
	Semicolon
	Comma
	Period
	Colon
	OpeningParenthesis
	ClosingParenthesis
	OpeningBrace
	ClosingBrace
	OpeningSquareBrace
	ClosingSquareBrace
	punctuationEnd

	whitespaceStart
	Newline
	Space
	Tab
	VerticalTab
	whiteSpaceEnd

	literalsStart
	// TODO: Add literals.
	Literal
	literalsEnd
)

var (
	tokens = [...]string{
		Absolute:       "absolute",
		Abstract:       "abstract",
		All:            "all",
		And:            "and",
		AndThen:        "and_then",
		Array:          "array",
		As:             "as",
		Asm:            "asm",
		Asmname:        "asmname",
		Attribute:      "attribute",
		Begin:          "begin",
		Bindable:       "bindable",
		C:              "c",
		Case:           "case",
		CLanguage:      "c_language",
		Class:          "class",
		Const:          "const",
		Constructor:    "constructor",
		Destructor:     "destructor",
		Div:            "div",
		Do:             "do",
		Downto:         "downto",
		Else:           "else",
		End:            "end",
		Export:         "export",
		Exports:        "exports",
		External:       "external",
		Far:            "far",
		File:           "file",
		Finalization:   "finalization",
		For:            "for",
		Forward:        "forward",
		Function:       "function",
		Goto:           "goto",
		If:             "if",
		Implementation: "implementation",
		Import:         "import",
		In:             "in",
		Inherited:      "inherited",
		Initialization: "initialization",
		Interface:      "interface",
		Interrupt:      "interrupt",
		Is:             "is",
		Label:          "label",
		Library:        "library",
		Mod:            "mod",
		Module:         "module",
		Name:           "name",
		Near:           "near",
		Nil:            "nil",
		Not:            "not",
		Object:         "object",
		Of:             "of",
		Only:           "only",
		Operator:       "operator",
		Or:             "or",
		OrElse:         "or_else",
		Otherwise:      "otherwise",
		Packed:         "packed",
		Pow:            "pow",
		Private:        "private",
		Procedure:      "procedure",
		Program:        "program",
		Property:       "property",
		Protected:      "protected",
		Public:         "public",
		Published:      "published",
		Qualified:      "qualified",
		Record:         "record",
		Repeat:         "repeat",
		Resident:       "resident",
		Restricted:     "restricted",
		Segment:        "segment",
		Set:            "set",
		Shl:            "shl",
		Shr:            "shr",
		Then:           "then",
		To:             "to",
		Type:           "type",
		Unit:           "unit",
		Until:          "until",
		Uses:           "uses",
		Value:          "value",
		Var:            "var",
		View:           "view",
		Virtual:        "virtual",
		While:          "while",
		With:           "with",
		Xor:            "xor",

		Eq:         "=",
		Ne:         "<>",
		Lt:         "<",
		Gt:         ">",
		Lte:        "<=",
		Gte:        ">=",
		Plus:       "+",
		Minus:      "-",
		Multiply:   "*",
		Divide:     "/",
		ShiftLeft:  "<<",
		ShiftRight: ">>",
		Assign:     ":=",

		Semicolon:          ";",
		Comma:              ",",
		Period:             ".",
		Colon:              ":",
		OpeningParenthesis: "(",
		ClosingParenthesis: ")",
		OpeningBrace:       "{",
		ClosingBrace:       "}",
		OpeningSquareBrace: "[",
		ClosingSquareBrace: "]",

		Newline:     "\n",
		Space:       " ",
		Tab:         "\t",
		VerticalTab: "\v",
	}
	ids map[string]ID
)

func init() {
	ids = make(map[string]ID)
	for id, token := range tokens {
		if token != "" {
			ids[token] = ID(id)
		}
	}
}

func (i ID) String() string {
	switch i {
	case Unknown:
		return "<unknown>"
	case UserDefined:
		return "<user defined>"
	case EOF:
		return "<EOF>"
	case Literal:
		return "<literal>"
	default:
		return tokens[i]
	}
}

func (i ID) IsKeyword() bool {
	return i > keywordsStart && i < keywordsEnd
}

func (i ID) IsOperator() bool {
	return i > operatorsStart && i < operatorsEnd
}

func (i ID) IsPunctuation() bool {
	return i > punctuationStart && i < punctuationEnd
}

func (i ID) IsWhitespace() bool {
	return i > whitespaceStart && i < whiteSpaceEnd
}

func (i ID) IsLiteral() bool {
	return i > literalsStart && i < literalsEnd
}

func GetID(token string) ID {
	return ids[token]
}

func ByID(id ID) string {
	return tokens[id]
}
