package aggregates

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/bnf"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var _ bnf.BNF = &Either{}

type Either struct {
	Name string
	BNFs []bnf.BNF
}

func (e Either) Accept(tokensCh *channel.TransactionChannel[option.Option[token.Token]]) error {
	defer tokensCh.Rollback()

	log.Print(e)

	var lastError error
	for _, item := range e.BNFs {
		if err := item.Accept(tokensCh.StartTx()); errors.Is(err, bnf.ErrUnexpectedToken) {
			lastError = err
			continue
		} else if err != nil {
			return fmt.Errorf("error in optional: %w", err)
		}

		tokensCh.Commit()

		return nil
	}

	return fmt.Errorf("the token is not in a list of expected tokens: %w", lastError)
}

func (e Either) String() string {
	if e.Name != "" {
		return fmt.Sprintf("either %q", e.Name)
	} else {
		return fmt.Sprintf("either")
	}
}
