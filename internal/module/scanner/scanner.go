package scanner

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/iskorotkov/compiler/internal/data/literal"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
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

func (l Scanner) Scan(ctx interface{}, input <-chan option.Option[literal.Literal]) <-chan option.Option[token.Token] {
	ch := make(chan option.Option[token.Token], l.buffer)

	go func() {
		defer close(ch)

		for item := range input {
			lit, err := item.Unwrap()
			if err != nil {
				ch <- option.Err[token.Token](fmt.Errorf("error passed from reader: %w", err))
				continue
			}

			addTypedToken(ctx, lit, ch)
		}

		ch <- option.Ok(token.Token{ID: token.EOF})
	}()

	return ch
}

//goland:noinspection GoUnusedParameter
func addTypedToken(ctx interface{}, lit literal.Literal, ch chan<- option.Option[token.Token]) {
	// Whitespace only - skip it.
	if len(strings.TrimSpace(lit.Value)) == 0 {
		return
	}

	id := token.GetID(lit.Value)
	if id != token.Unknown {
		ch <- option.Ok(token.New(id, lit))
		return
	}

	// Constants.
	if intConstantRegex.MatchString(lit.Value) {
		ch <- option.Ok(token.New(token.IntLiteral, lit))
		return
	}

	if doubleConstantRegex.MatchString(lit.Value) {
		ch <- option.Ok(token.New(token.DoubleLiteral, lit))
		return
	}

	if boolConstantRegex.MatchString(lit.Value) {
		ch <- option.Ok(token.New(token.BoolLiteral, lit))
		return
	}

	// User identifiers.
	if userIdentifierRegex.MatchString(lit.Value) {
		ch <- option.Ok(token.New(token.UserDefined, lit))
		return
	}

	ch <- option.Err[token.Token](fmt.Errorf("unknown token %s at position %v", lit.Value, lit.Position))
}
