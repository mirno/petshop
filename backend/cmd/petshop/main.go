package main

import (
	"github.com/atoscerebro/eviden-petshop/internal/drivers/printer"
	"github.com/atoscerebro/eviden-petshop/internal/entities"
	"github.com/atoscerebro/eviden-petshop/internal/usecases"
	"github.com/google/uuid"
)

var (
	idRex       = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	idSlytherin = uuid.MustParse("00000000-0000-0000-0000-000000000002")

	petFixtures = []entities.Pet{
		{ID: idRex, Name: "Rex", Type: entities.Dog, Price: 80},
		{ID: idSlytherin, Name: "Slytherin", Type: entities.Snake, Price: 400},
	}
)

func main() {
	// printer := printer.JSONPrinter{Path: "pets.json"}
	printer := printer.ConsolePrinter{}
	petshop := usecases.NewPetshop(&printer, petFixtures...)

	petshop.PrintPets()

	pets := petshop.ListPets()
	pets[0].PrintPrice()
}
