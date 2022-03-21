package channel_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iskorotkov/compiler/internal/channel"
)

func TestChannelWrapping(t *testing.T) {
	t.Parallel()

	ch := channel.NewTransactionChannel(channel.FromSlice([]int{1, 2, 3}))
	assert.True(t, ch.Open())

	assert.Equal(t, 1, ch.Read())
	assert.True(t, ch.Open())

	assert.Equal(t, 2, ch.Read())
	assert.True(t, ch.Open())

	assert.Equal(t, 3, ch.Read())
	assert.False(t, ch.Open())
}

func TestCommit(t *testing.T) {
	t.Parallel()

	ch := channel.NewTransactionChannel(channel.FromSlice([]int{1, 2, 3}))
	assert.Equal(t, 1, ch.Read())
	assert.Equal(t, 2, ch.Read())
	assert.Equal(t, 3, ch.Read())
	assert.False(t, ch.Open())

	ch.Commit()
	assert.False(t, ch.Open())
}

func TestRollback(t *testing.T) {
	t.Parallel()

	ch := channel.NewTransactionChannel(channel.FromSlice([]int{1, 2, 3}))
	assert.Equal(t, 1, ch.Read())
	assert.Equal(t, 2, ch.Read())
	assert.Equal(t, 3, ch.Read())
	assert.False(t, ch.Open())

	ch.Rollback()
	assert.True(t, ch.Open())

	assert.Equal(t, 1, ch.Read())
	assert.Equal(t, 2, ch.Read())
	assert.Equal(t, 3, ch.Read())
	assert.False(t, ch.Open())
}

func TestPartialRollback(t *testing.T) {
	t.Parallel()

	ch := channel.NewTransactionChannel(channel.FromSlice([]int{1, 2, 3}))
	assert.Equal(t, 1, ch.Read())
	assert.Equal(t, 2, ch.Read())
	assert.True(t, ch.Open())

	ch.Rollback()
	assert.True(t, ch.Open())

	assert.Equal(t, 1, ch.Read())
	assert.Equal(t, 2, ch.Read())
	assert.Equal(t, 3, ch.Read())
	assert.False(t, ch.Open())
}

func TestCommitRollback(t *testing.T) {
	t.Parallel()

	ch := channel.NewTransactionChannel(channel.FromSlice([]int{1, 2, 3}))
	assert.Equal(t, 1, ch.Read())
	assert.Equal(t, 2, ch.Read())
	assert.True(t, ch.Open())

	ch.Commit()
	assert.True(t, ch.Open())

	ch.Rollback()
	assert.True(t, ch.Open())

	assert.Equal(t, 3, ch.Read())
	assert.False(t, ch.Open())
}

func TestNestedTransactions(t *testing.T) {
	t.Parallel()

	ch := channel.NewTransactionChannel(channel.FromSlice([]int{1, 2, 3}))
	assert.Equal(t, 1, ch.Read())
	assert.True(t, ch.Open())

	ch = ch.StartTx()
	assert.Equal(t, 2, ch.Read())
	assert.True(t, ch.Open())

	ch = ch.StartTx()
	assert.Equal(t, 3, ch.Read())
	assert.False(t, ch.Open())

	ch.Rollback()
	assert.Equal(t, 3, ch.Read())
	assert.False(t, ch.Open())

	ch.Rollback()
	assert.Equal(t, 2, ch.Read())
	assert.True(t, ch.Open())

	assert.Equal(t, 3, ch.Read())
	assert.False(t, ch.Open())

	ch.Commit()
	assert.False(t, ch.Open())
}
