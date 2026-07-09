package usecases

import (
	"github.com/atoscerebro/eviden-petshop/internal/entities"
)

// TODO: Refactor Petshop, since we don't need an interface to meet the functions.
// Use an adapter pattern to wire up the Petshop + Printer. To comply with Single-responsibilty principle.

// type Petshop interface {
// 	ListPets() []entities.Pet
// 	AddPet(entities.Pet) error
// }

type Petshop struct {
	pets []entities.Pet
}

// NewPetshop creates a new instance of PetShop
func NewPetshop(pets ...entities.Pet) *Petshop {
	return &Petshop{
		pets: pets,
	}
}

func (shop *Petshop) ListPets() []entities.Pet {
	return shop.pets
}

func (shop *Petshop) AddPet(pet entities.Pet) error {
	shop.pets = append(shop.pets, pet)
	return nil
}
