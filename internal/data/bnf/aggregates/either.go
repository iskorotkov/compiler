package aggregates

import (
	"errors"
	"fmt"
	"strings"

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

	if e.Name != "" {
		log.Printf("%s: expecting %v", strings.ToUpper(e.Name), e)
	} else {
		log.Printf("expecting %v", e)
	}

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
	var values []string
	for _, value := range e.BNFs {
		values = append(values, value.String())
	}

	return fmt.Sprintf("either of (%s)", strings.Join(values, "; "))
}
