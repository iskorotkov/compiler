package bnf

import (
	"errors"
	"fmt"
	"strings"

	"github.com/iskorotkov/compiler/internal/channel"
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var _ BNF = &Several{}

type Several struct {
	Name string
	BNF
}

func (s Several) Accept(tokensCh *channel.TransactionChannel[option.Option[token.Token]]) error {
	defer tokensCh.Rollback()

	if s.Name != "" {
		log.Printf("%s: expecting %v", strings.ToUpper(s.Name), s)
	} else {
		log.Printf("expecting %v", s)
	}

	for {
		if err := s.BNF.Accept(tokensCh.StartTx()); errors.Is(err, ErrUnexpectedToken) {
			tokensCh.Commit()
			return nil
		} else if err != nil {
			return fmt.Errorf("error in optional: %w", err)
		}
	}
}

func (s Several) String() string {
	return fmt.Sprintf("several of %v", s.BNF)
}
