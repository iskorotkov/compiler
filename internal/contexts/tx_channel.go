package contexts

import (
	"github.com/iskorotkov/compiler/internal/data/token"
	"github.com/iskorotkov/compiler/internal/fn/channels"
	"github.com/iskorotkov/compiler/internal/fn/options"
)

var _ TxChannelContext = (*txChannelContext)(nil)

type txChannelContext struct {
	ch *channels.TxChannel[options.Option[token.Token]]
}

func NewTxChannelContext(ch *channels.TxChannel[options.Option[token.Token]]) TxChannelContext {
	return &txChannelContext{ch: ch}
}

func (t txChannelContext) TxChannel() *channels.TxChannel[options.Option[token.Token]] {
	return t.ch
}
