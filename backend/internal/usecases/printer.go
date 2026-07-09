package usecases

import "github.com/mirno/petshop/internal/entities"

// TODO: Use []bytes
type Printer interface {
	Print([]entities.Pet) // Ugly implementations, scales better if we move to []byte
}
