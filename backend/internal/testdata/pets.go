package testdata

import (
	"github.com/google/uuid"
	"github.com/mirno/petshop/internal/entities"
)

var (
	IDRex       = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	IDSlytherin = uuid.MustParse("00000000-0000-0000-0000-000000000002")

	PetRex       = entities.Pet{Id: IDRex, Name: "Rex", Type: entities.Dog, Price: 80}
	PetSlytherin = entities.Pet{Id: IDSlytherin, Name: "Slytherin", Type: entities.Snake, Price: 400}

	PetFixtures = []entities.Pet{
		PetRex,
		PetSlytherin,
	}
)
