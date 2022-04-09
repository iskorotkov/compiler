package context

import (
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channel"
	"github.com/iskorotkov/compiler/internal/fn/option"
)

var _ TxChannelContext = (*txChannelContext)(nil)

type txChannelContext struct {
	ch *channel.TxChannel[option.Option[token.Token]]
}

func NewTxChannelContext(ch *channel.TxChannel[option.Option[token.Token]]) TxChannelContext {
	return &txChannelContext{ch: ch}
}

func (t txChannelContext) TxChannel() *channel.TxChannel[option.Option[token.Token]] {
	return t.ch
}
