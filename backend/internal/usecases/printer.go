package usecases

import "github.com/atoscerebro/eviden-petshop/internal/entities"

type Printer interface {
	Print([]entities.Pet) // Ugly implementations, scales better if we move to []byte
}
