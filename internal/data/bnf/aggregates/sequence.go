package aggregates

import (
	"fmt"

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

	log.Print(s)

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
		return "empty"
	}

	if s.Name != "" {
		return fmt.Sprintf("sequence %q", s.Name)
	} else {
		return fmt.Sprintf("sequence")
	}
}
