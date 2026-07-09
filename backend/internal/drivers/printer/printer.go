package printer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/atoscerebro/eviden-petshop/internal/entities"
)

type ConsolePrinter struct {
}

func (p *ConsolePrinter) Print(pets []entities.Pet) {
	for _, pet := range pets {
		fmt.Printf("%s (%s): %.2f\n", pet.Name, pet.Type, float64(pet.Price))
	}
}

type JSONPrinter struct {
	Path string
}

func (p *JSONPrinter) Print(pets []entities.Pet) {
	path := p.Path

	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}

	jsonData, err := json.Marshal(pets)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		os.Exit(1)
	}

	file.Write(jsonData)
}
