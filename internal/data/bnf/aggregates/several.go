package aggregates

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/bnf"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var _ bnf.BNF = &Several{}

type Several struct {
	Name string
	bnf.BNF
}

func (s Several) Accept(tokensCh *channel.TransactionChannel[option.Option[token.Token]]) error {
	defer tokensCh.Rollback()

	log.Print(s)

	for {
		if err := s.BNF.Accept(tokensCh.StartTx()); errors.Is(err, bnf.ErrUnexpectedToken) {
			tokensCh.Commit()
			return nil
		} else if err != nil {
			return fmt.Errorf("error in optional: %w", err)
		}
	}
}

func (s Several) String() string {
	if s.Name != "" {
		return fmt.Sprintf("several %q", s.Name)
	} else {
		return fmt.Sprintf("several")
	}
}
