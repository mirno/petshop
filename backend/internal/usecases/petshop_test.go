package usecases_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/mirno/petshop/internal/drivers/inmemorykvstore"
	"github.com/mirno/petshop/internal/entities"
	"github.com/mirno/petshop/internal/testdata"
	"github.com/mirno/petshop/internal/usecases"
)

func TestNewPetshopSeedsPetsWithIDs(t *testing.T) {
	store := inmemorykvstore.NewInMemoryKVStore[entities.Pet]()

	usecases.NewPetshop(store, testdata.PetRex, entities.Pet{Name: "No ID"})

	pet, err := store.Get(testdata.IDRex.String())
	if err != nil {
		t.Fatalf("expected seeded pet: %v", err)
	}
	if pet.Name != testdata.PetRex.Name {
		t.Fatalf("expected seeded pet %q, got %q", testdata.PetRex.Name, pet.Name)
	}

	keys, err := store.Keys()
	if err != nil {
		t.Fatalf("keys: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("expected only pets with IDs to be seeded, got keys %#v", keys)
	}
}

func TestUpdatePetReplacesExistingPet(t *testing.T) {
	store := inmemorykvstore.NewInMemoryKVStore[entities.Pet]()
	shop := usecases.NewPetshop(store, testdata.PetRex)
	updated := testdata.PetRex
	updated.Name = "Max"
	updated.Price = 95

	if err := shop.UpdatePet(updated); err != nil {
		t.Fatalf("update pet: %v", err)
	}

	got, err := store.Get(updated.Id.String())
	if err != nil {
		t.Fatalf("get updated pet: %v", err)
	}
	if got.Name != "Max" || got.Price != 95 {
		t.Fatalf("expected updated pet, got %#v", got)
	}
}

func TestUpdatePetRequiresExistingID(t *testing.T) {
	store := inmemorykvstore.NewInMemoryKVStore[entities.Pet]()
	shop := usecases.NewPetshop(store)

	err := shop.UpdatePet(entities.Pet{Name: "No ID"})
	if !errors.Is(err, usecases.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument for empty ID, got %v", err)
	}

	err = shop.UpdatePet(entities.Pet{Id: uuid.New(), Name: "Missing"})
	if !errors.Is(err, usecases.ErrKeyNotFoundError) {
		t.Fatalf("expected key not found for missing pet, got %v", err)
	}
}

func TestDeletePetRemovesExistingPet(t *testing.T) {
	store := inmemorykvstore.NewInMemoryKVStore[entities.Pet]()
	shop := usecases.NewPetshop(store, testdata.PetRex)

	if err := shop.DeletePet(testdata.PetRex); err != nil {
		t.Fatalf("delete pet: %v", err)
	}

	_, err := store.Get(testdata.IDRex.String())
	if !errors.Is(err, usecases.ErrKeyNotFoundError) {
		t.Fatalf("expected deleted pet to be missing, got %v", err)
	}
}

func TestDeletePetRequiresExistingID(t *testing.T) {
	store := inmemorykvstore.NewInMemoryKVStore[entities.Pet]()
	shop := usecases.NewPetshop(store)

	err := shop.DeletePet(entities.Pet{Name: "No ID"})
	if !errors.Is(err, usecases.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument for empty ID, got %v", err)
	}

	err = shop.DeletePet(entities.Pet{Id: uuid.New(), Name: "Missing"})
	if !errors.Is(err, usecases.ErrKeyNotFoundError) {
		t.Fatalf("expected key not found for missing pet, got %v", err)
	}
}
