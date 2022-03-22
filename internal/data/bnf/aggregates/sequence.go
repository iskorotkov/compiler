package aggregates

import (
	"fmt"
	"strings"

	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/bnf"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var _ bnf.BNF = &Sequence{}

type Sequence struct {
	Name string
	BNFs []bnf.BNF
}

func (s Sequence) Accept(tokensCh *channel.TransactionChannel[option.Option[token.Token]]) error {
	defer tokensCh.Rollback()

	if s.Name != "" {
		log.Printf("%s: expecting %v", strings.ToUpper(s.Name), s)
	} else {
		log.Printf("expecting %v", s)
	}

	for i, item := range s.BNFs {
		if err := item.Accept(tokensCh.StartTx()); err != nil {
			return fmt.Errorf("error in composite at index %d: %w", i, err)
		}
	}

	tokensCh.Commit()

	return nil
}

func (s Sequence) String() string {
	if len(s.BNFs) == 0 {
		return "<empty>"
	}

	var values []string
	for _, value := range s.BNFs {
		values = append(values, value.String())
	}

	return fmt.Sprintf("sequence of (%s)", strings.Join(values, "; "))
}
