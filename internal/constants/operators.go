package constants

// Operators are taken from https://wiki.freepascal.org/ToID.
//
//   const ops = [...document.querySelectorAll('table > tbody > tr > td:first-child')].map(_ => _.innerText).filter(_ => _.indexOf(' ') === -1);
//   const filtered = ops.filter((_, i) => ops.indexOf(_) === i);
//   `map[string]OperatorID{${filtered.map((_, i) => `"${_}": ${i+1}`)}}`
var Operators = map[string]ID{
	"=":       1,
	"<>":      2,
	"<":       3,
	">":       4,
	"<=":      5,
	">=":      6,
	"in":      7,
	"+":       8,
	"-":       9,
	"*":       10,
	"/":       11,
	"div":     12,
	"mod":     13,
	"not":     14,
	"and":     15,
	"or":      16,
	"xor":     17,
	"shl":     18,
	"shr":     19,
	"<<":      20,
	">>":      21,
	"><":      22,
	"include": 23,
	"exclude": 24,
	"is":      25,
	"as":      26,
	"^":       27,
	":=":      28,
}
