package typechecker

import (
	"github.com/iskorotkov/compiler/internal/data/symbol"
)

type TypeConverter struct {
	possibleConversions map[symbol.BuiltinType]map[symbol.BuiltinType]symbol.BuiltinType
}

func NewTypeConverter() TypeConverter {
	return TypeConverter{
		// Type 1 and type 2 -> type 3.
		possibleConversions: map[symbol.BuiltinType]map[symbol.BuiltinType]symbol.BuiltinType{
			symbol.BuiltinTypeInt: {
				// All types are convertible to unknown and produce the same type.
				// It's useful when expression type is not known yet.
				symbol.BuiltinTypeUnknown: symbol.BuiltinTypeInt,
				symbol.BuiltinTypeInt:     symbol.BuiltinTypeInt,
				symbol.BuiltinTypeDouble:  symbol.BuiltinTypeDouble,
			},
			symbol.BuiltinTypeDouble: {
				symbol.BuiltinTypeUnknown: symbol.BuiltinTypeDouble,
				symbol.BuiltinTypeDouble:  symbol.BuiltinTypeDouble,
			},
			symbol.BuiltinTypeBool: {
				symbol.BuiltinTypeUnknown: symbol.BuiltinTypeBool,
				symbol.BuiltinTypeBool:    symbol.BuiltinTypeBool,
			},
		},
	}
}

func (c TypeConverter) Convert(ctx interface{}, from, to symbol.BuiltinType) (symbol.BuiltinType, bool) {
	val, ok := c.possibleConversions[from][to]
	return val, ok
}
