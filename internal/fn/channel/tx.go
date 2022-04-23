package channel

import (
	"sync"
)

type TxChannel[T any] struct {
	sourceCh <-chan T

	transactions [][]T
	rollback     []T
	closed       bool
	m            sync.Mutex
}

func NewTxChannel[T any](ch <-chan T) *TxChannel[T] {
	return &TxChannel[T]{
		sourceCh: ch,
	}
}

// Commit commits all reads and starts a new transaction.
func (c *TxChannel[T]) Commit() {
	c.m.Lock()
	defer c.m.Unlock()

	if len(c.transactions) == 0 {
		return
	}

	c.transactions = c.transactions[:len(c.transactions)-1]
}

// Rollback adds values from ongoing transaction to the rollback list and starts a new transaction.
func (c *TxChannel[T]) Rollback() {
	c.m.Lock()
	defer c.m.Unlock()

	if len(c.transactions) == 0 {
		return
	}

	c.rollback = append(c.transactions[len(c.transactions)-1], c.rollback...)
	c.transactions = c.transactions[:len(c.transactions)-1]

	if len(c.rollback) > 0 {
		c.closed = false
	}
}

// Read reads from the rollback list or from source channel.
func (c *TxChannel[T]) Read() T {
	c.m.Lock()
	defer c.m.Unlock()

	var (
		nextElement  T
		fromRollback bool
	)
	if len(c.rollback) > 0 {
		nextElement, c.rollback = c.rollback[0], c.rollback[1:]
		fromRollback = true
	}

	if !fromRollback {
		var open bool
		nextElement, open = <-c.sourceCh
		if !open {
			c.closed = true
		}
	}

	if len(c.transactions) == 0 {
		c.StartTx()
	}

	c.transactions[len(c.transactions)-1] = append(c.transactions[len(c.transactions)-1], nextElement)

	return nextElement
}

// Open checks whether the channel has rollback values or if the source channel is open.
func (c *TxChannel[T]) Open() bool {
	c.m.Lock()
	defer c.m.Unlock()

	if c.closed {
		return false
	}

	if len(c.rollback) > 0 {
		return true
	}

	item, open := <-c.sourceCh
	if open {
		c.rollback = append([]T{item}, c.rollback...)
	} else {
		c.closed = true
	}

	return !c.closed
}

// StartTx starts a new nested transaction.
func (c *TxChannel[T]) StartTx() *TxChannel[T] {
	c.transactions = append(c.transactions, []T{})
	return c
}
