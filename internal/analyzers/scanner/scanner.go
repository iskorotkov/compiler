package scanner

import (
	"fmt"
	"regexp"

	"github.com/iskorotkov/compiler/internal/constants"
	"github.com/iskorotkov/compiler/internal/data/literal"
	"github.com/iskorotkov/compiler/internal/data/token"
)

var (
	intConstantRegex    = regexp.MustCompile(`^\d+$`)
	doubleConstantRegex = regexp.MustCompile(`^\d+\.\d+$`)
	boolConstantRegex   = regexp.MustCompile(`^true|false$`)
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

func (l Scanner) Scan(input <-chan literal.Option) <-chan token.Option {
	ch := make(chan token.Option, l.buffer)

	go func() {
		defer close(ch)

		for item := range input {
			lit, err := item.Unwrap()
			if err != nil {
				ch <- token.Err(fmt.Errorf("error passed from reader: %w", err))
				continue
			}

			if id := constants.Keywords[lit.Value]; id != constants.None {
				ch <- token.Ok(token.New(token.TypeKeyword, id, lit, nil))
			} else if id := constants.Operators[lit.Value]; id != constants.None {
				ch <- token.Ok(token.New(token.TypeOperator, id, lit, nil))
			} else if id := constants.Punctuation[lit.Value]; id != constants.None {
				ch <- token.Ok(token.New(token.TypePunctuation, id, lit, nil))
			} else if intConstantRegex.Match([]byte(lit.Value)) || doubleConstantRegex.Match([]byte(lit.Value)) || boolConstantRegex.Match([]byte(lit.Value)) {
				// TODO: Pass value to next analyzers.
				ch <- token.Ok(token.New(token.TypeConstant, 0, lit, nil))
			} else if userIdentifierRegex.Match([]byte(lit.Value)) {
				ch <- token.Ok(token.New(token.TypeUserIdentifier, 0, lit, nil))
			} else {
				ch <- token.Err(fmt.Errorf("unknown token %s at position %v", lit.Value, lit.Position))
			}
		}
	}()

	return ch
}
