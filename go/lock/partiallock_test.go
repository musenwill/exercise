package lock_test

import (
	"testing"

	"github.com/musenwill/exercise/lock"
	"github.com/stretchr/testify/require"
)

func TestPartiallockOK(t *testing.T) {
	pre := lock.NewPartialLock()
	post := lock.NewPartialLock()
	pre.WithPostLock(post)

	require.NoError(t, pre.Lock())
	require.NoError(t, pre.Unlock())

	require.NoError(t, post.Lock())
	require.NoError(t, post.Unlock())

	require.NoError(t, pre.Lock())
	require.NoError(t, pre.Unlock())

	require.NoError(t, pre.Lock())
	require.NoError(t, post.Lock())
	require.NoError(t, pre.Unlock())
	require.NoError(t, post.Unlock())

	require.NoError(t, pre.Lock())
	require.NoError(t, post.Lock())
	require.NoError(t, post.Unlock())
	require.NoError(t, pre.Unlock())
}

func TestPartiallockFail(t *testing.T) {
	pre := lock.NewPartialLock()
	post := lock.NewPartialLock()
	pre.WithPostLock(post)

	require.NoError(t, post.Lock())
	require.Error(t, post.Lock())
	require.Error(t, pre.Lock())
	require.Error(t, pre.Unlock())

	third := lock.NewPartialLock()
	require.Error(t, third.WithPostLock(post))
}
