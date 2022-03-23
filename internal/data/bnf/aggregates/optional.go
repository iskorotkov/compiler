package aggregates

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/bnf"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var _ bnf.BNF = &Optional{}

type Optional struct {
	Name string
	bnf.BNF
}

func (o Optional) Accept(tokensCh *channel.TransactionChannel[option.Option[token.Token]]) error {
	defer tokensCh.Rollback()

	log.Print(o)

	if err := o.BNF.Accept(tokensCh.StartTx()); errors.Is(err, bnf.ErrUnexpectedToken) {
		// Return without committing tx.
		return nil
	} else if err != nil {
		return fmt.Errorf("error in optional: %w", err)
	}

	tokensCh.Commit()

	return nil
}

func (o Optional) String() string {
	if o.Name != "" {
		return fmt.Sprintf("optional %q", o.Name)
	} else {
		return fmt.Sprintf("optional")
	}
}
