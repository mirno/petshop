package usecases

import (
	"errors"

	"github.com/google/uuid"
	"github.com/mirno/petshop/internal/entities"
)

// TODO: Refactor Petshop, since we don't need an interface to meet the functions.
// Use an adapter pattern to wire up the Petshop + Printer. To comply with Single-responsibilty principle.

// type Petshop interface {
// 	ListPets() []entities.Pet
// 	AddPet(entities.Pet) error
// }

type Petshop struct {
	store Store[entities.Pet]
}

// NewPetshop creates a new instance of PetShop
func NewPetshop(store Store[entities.Pet], pets ...entities.Pet) *Petshop {
	for _, pet := range pets {
		if pet.Id != uuid.Nil {
			store.Save(pet.Id.String(), pet) //nolint: errcheck // TODO: fix; skip now due to lack of err return on the constructor
		}
	}

	return &Petshop{
		store: store,
	}
}

// ListPets is a User function to list all Pets
func (shop *Petshop) ListPets() []entities.Pet {
	var petList []entities.Pet

	// Return nothing if errors out for now
	keys, err := shop.store.Keys()
	if err != nil {
		return nil
	}

	for _, key := range keys {
		pet, _ := shop.store.Get(key) //nolint: errcheck ignore err and list other pets

		petList = append(petList, pet)
	}

	return petList
}

// AddPet will be an administrative function to add new pets to the shop
func (shop *Petshop) AddPet(pet entities.Pet) (uuid.UUID, error) {
	// TODO: Implement store.Has() to validate if the key already exists
	if pet.Id != uuid.Nil {
		return uuid.Nil, errors.Join(errors.New("enter a pet without id"), ErrInvalidArgument)
	}

	pet.Id = uuid.New()

	if err := shop.store.Save(pet.Id.String(), pet); err != nil {
		return uuid.Nil, err
	}

	return pet.Id, nil
}

// UpdatePet will be an administrative function to update existing pets in the shop.
func (shop *Petshop) UpdatePet(pet entities.Pet) error {
	if pet.Id == uuid.Nil {
		return errors.Join(errors.New("enter a pet with id"), ErrInvalidArgument)
	}

	if _, err := shop.store.Get(pet.Id.String()); err != nil {
		return err
	}

	return shop.store.Save(pet.Id.String(), pet)
}

// DeletePet will be an administrative function to delete existing pets from the shop.
func (shop *Petshop) DeletePet(pet entities.Pet) error {
	if pet.Id == uuid.Nil {
		return errors.Join(errors.New("enter a pet with id"), ErrInvalidArgument)
	}

	if _, err := shop.store.Get(pet.Id.String()); err != nil {
		return err
	}

	return shop.store.Delete(pet.Id.String())
}
