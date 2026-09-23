package fs_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/richargh/workdown/internal/kernel/fs"
)

func TestMemStoreContract(t *testing.T) {
	testStoreContract(t, func(t *testing.T) fs.Store {
		t.Helper()
		return fs.NewMemStore()
	})
}

func TestOSStoreContract(t *testing.T) {
	testStoreContract(t, func(t *testing.T) fs.Store {
		t.Helper()
		return fs.NewOSStore(t.TempDir())
	})
}

func testStoreContract(t *testing.T, newStore func(t *testing.T) fs.Store) {
	t.Helper()

	t.Run("missing file does not exist", func(t *testing.T) {
		// given
		store := newStore(t)

		// when
		exists, err := store.Exists(context.Background(), "missing.md")

		// then
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("reading missing file returns error", func(t *testing.T) {
		// given
		store := newStore(t)

		// when
		data, err := store.ReadFile(context.Background(), "missing.md")

		// then
		require.Error(t, err)
		require.Nil(t, data)
	})

	t.Run("written file exists and can be read", func(t *testing.T) {
		// given
		store := newStore(t)

		// when
		err := store.WriteFile(context.Background(), "issues/PROJ-1.md", []byte("hello"))
		require.NoError(t, err)

		// then
		exists, err := store.Exists(context.Background(), "issues/PROJ-1.md")
		require.NoError(t, err)
		require.True(t, exists)

		data, err := store.ReadFile(context.Background(), "issues/PROJ-1.md")
		require.NoError(t, err)
		require.Equal(t, []byte("hello"), data)
	})

	t.Run("write overwrites existing file", func(t *testing.T) {
		// given
		store := newStore(t)
		require.NoError(t, store.WriteFile(context.Background(), "issue.md", []byte("old")))

		// when
		require.NoError(t, store.WriteFile(context.Background(), "issue.md", []byte("new")))

		// then
		data, err := store.ReadFile(context.Background(), "issue.md")
		require.NoError(t, err)
		require.Equal(t, []byte("new"), data)
	})

	t.Run("write defensively copies input", func(t *testing.T) {
		// given
		store := newStore(t)
		data := []byte("original")

		// when
		require.NoError(t, store.WriteFile(context.Background(), "issue.md", data))
		data[0] = 'X'

		// then
		stored, err := store.ReadFile(context.Background(), "issue.md")
		require.NoError(t, err)
		require.Equal(t, []byte("original"), stored)
	})

	t.Run("read returns defensive copy", func(t *testing.T) {
		// given
		store := newStore(t)
		require.NoError(t, store.WriteFile(context.Background(), "issue.md", []byte("original")))

		// when
		firstRead, err := store.ReadFile(context.Background(), "issue.md")
		require.NoError(t, err)
		firstRead[0] = 'X'

		// then
		secondRead, err := store.ReadFile(context.Background(), "issue.md")
		require.NoError(t, err)
		require.Equal(t, []byte("original"), secondRead)
	})
}
