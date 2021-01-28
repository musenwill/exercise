package lock_test

import (
	"testing"

	"github.com/musenwill/exercise/lock"
	"github.com/stretchr/testify/require"
)

func TestPartiallockOK(t *testing.T) {
	mu := lock.NewPartialLock()

	require.NoError(t, mu.Prelock())
	require.NoError(t, mu.Preunlock())

	require.NoError(t, mu.Postlock())
	require.NoError(t, mu.Postunlock())

	require.NoError(t, mu.Prelock())
	require.NoError(t, mu.Preunlock())

	require.NoError(t, mu.Prelock())
	require.NoError(t, mu.Postlock())
	require.NoError(t, mu.Preunlock())
	require.NoError(t, mu.Postunlock())

	require.NoError(t, mu.Prelock())
	require.NoError(t, mu.Postlock())
	require.NoError(t, mu.Postunlock())
	require.NoError(t, mu.Preunlock())
}

func TestPartiallockFail(t *testing.T) {
	mu := lock.NewPartialLock()

	require.NoError(t, mu.Postlock())
	require.Error(t, mu.Postlock())
	require.Error(t, mu.Prelock())
	require.Error(t, mu.Preunlock())
}
