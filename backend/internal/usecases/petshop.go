package usecases

import (
	"github.com/atoscerebro/eviden-petshop/internal/entities"
)

type Petshop interface {
	ListPets() []entities.Pet
	AddPet(entities.Pet) error
}

type PetShopWithPrinter struct {
	pets []entities.Pet
	Printer
}

// NewPetshop creates a new instance of PetShop
func NewPetshop(printer Printer, pets ...entities.Pet) *PetShopWithPrinter {
	return &PetShopWithPrinter{
		pets:    pets,
		Printer: printer,
	}
}

func (shop *PetShopWithPrinter) ListPets() []entities.Pet {
	return shop.pets
}

func (shop *PetShopWithPrinter) AddPet(pet entities.Pet) error {
	shop.pets = append(shop.pets, pet)
	return nil
}

func (shop *PetShopWithPrinter) PrintPets() {
	shop.Printer.Print(shop.pets)
}
