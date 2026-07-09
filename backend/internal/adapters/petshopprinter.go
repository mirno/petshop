package adapters

import "github.com/mirno/petshop/internal/usecases"

// Petshop pinters decorates the Pethop with a printer mechanism, using dependency injection.
type PetshopPrinter struct {
	*usecases.Petshop // Inherit the Petshop functors.
	usecases.Printer
}

func (shop *PetshopPrinter) PrintPets() {
	shop.Print(shop.ListPets())
}
