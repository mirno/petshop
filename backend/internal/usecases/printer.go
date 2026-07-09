package usecases

import "github.com/atoscerebro/eviden-petshop/internal/entities"

// TODO: Use []bytes
type Printer interface {
	Print([]entities.Pet) // Ugly implementations, scales better if we move to []byte
}
