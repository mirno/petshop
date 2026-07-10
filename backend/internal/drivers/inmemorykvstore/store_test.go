package inmemorykvstore_test

import (
	"testing"

	"github.com/mirno/petshop/internal/drivers/inmemorykvstore"
	"github.com/mirno/petshop/internal/entities"
	"github.com/mirno/petshop/internal/testdata"
	"github.com/mirno/petshop/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreSavesAndGetsNilPetPointer(t *testing.T) {
	store := inmemorykvstore.NewInMemoryKVStore[*entities.Pet]()

	err := store.Save("empty", nil)
	require.NoError(t, err)

	got, err := store.Get("empty")
	require.NoError(t, err)
	require.Empty(t, got)

	keys, err := store.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.NotEmpty(t, keys[0])
}

func TestStoreSupportsPetEntityGeneric(t *testing.T) {
	var _ usecases.Store[entities.Pet] = inmemorykvstore.NewInMemoryKVStore[entities.Pet]()

	store := inmemorykvstore.NewInMemoryKVStore[entities.Pet]()

	want := testdata.PetRex
	want.Metadata.Set("Region", "Europe")

	err := store.Save(want.Name, want)
	require.NoError(t, err)

	got, err := store.Get(want.Name)
	require.NoError(t, err)

	assert.Equal(t, got, want)

	err = store.Delete(want.Name)
	require.NoError(t, err)

	_, err = store.Get(want.Name)
	require.ErrorIs(t, err, usecases.ErrKeyNotFoundError)
}
