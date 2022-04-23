package scanner

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/data/literal"
	"github.com/iskorotkov/compiler/internal/data/token"
)

var (
	intConstantRegex    = regexp.MustCompile(`^\d+$`)
	doubleConstantRegex = regexp.MustCompile(`^\d+\.\d+$`)
	boolConstantRegex   = regexp.MustCompile(`^true$|^false$`)
	userIdentifierRegex = regexp.MustCompile(`^(?i)[a-z_]\w*$`)
)

type Scanner struct {
	buffer int
}

func New(buffer int) *Scanner {
	return &Scanner{
		buffer: buffer,
	}
}

func (l Scanner) Scan(ctx interface{ context.ErrorsContext }, input <-chan literal.Literal) <-chan token.Token {
	ch := make(chan token.Token, l.buffer)

	go func() {
		defer close(ch)

		for lit := range input {
			addTypedToken(ctx, lit, ch)
		}

		ch <- token.Token{ID: token.EOF}
	}()

	return ch
}

func addTypedToken(ctx interface{ context.ErrorsContext }, lit literal.Literal, ch chan<- token.Token) {
	// Whitespace only - skip it.
	if len(strings.TrimSpace(lit.Value)) == 0 {
		return
	}

	id := token.GetID(lit.Value)
	if id != token.Unknown {
		ch <- token.New(id, lit)
		return
	}

	// Constants.
	if intConstantRegex.MatchString(lit.Value) {
		ch <- token.New(token.IntLiteral, lit)
		return
	}

	if doubleConstantRegex.MatchString(lit.Value) {
		ch <- token.New(token.DoubleLiteral, lit)
		return
	}

	if boolConstantRegex.MatchString(lit.Value) {
		ch <- token.New(token.BoolLiteral, lit)
		return
	}

	// User identifiers.
	if userIdentifierRegex.MatchString(lit.Value) {
		ch <- token.New(token.UserDefined, lit)
		return
	}

	ctx.AddError(context.ErrorSourceScanner, lit.Position, fmt.Errorf("unrecognized token: %s", lit.Value))
}
