package entities

import (
	"github.com/mirno/petshop/pkg/petshop"
)

const (
	Snake petshop.AnimalType = "snake"
	Dog   petshop.AnimalType = "dog"
)

type Metadata map[string]string

type Pet = petshop.Pet // is Alias; define as type `Pet = petshop.Pet` to add functions.
