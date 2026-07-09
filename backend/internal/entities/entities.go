package entities

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	Snake AnimalType = "snake"
	Dog   AnimalType = "dog"
)

type AnimalType string

type Metadata map[string]string

type Pet struct {
	ID       uuid.UUID
	Name     string     `json:"name"`
	Type     AnimalType `json:"type"`
	Price    uint16     `json:"price"`
	Metadata Metadata
}

func (p *Pet) SetRegion(regionName string) {
	if p.Metadata == nil {
		p.defineMetadata() // Or it will panic
	}

	p.Metadata["Region"] = regionName
}

func (p *Pet) GetRegion() string {
	if p.Metadata == nil {
		p.defineMetadata()
	}

	return p.Metadata["Region"]
}

// PrintPrice with decimals by converting to float
func (p *Pet) PrintPrice() {
	fmt.Printf("%.2f\n", float64(p.Price))
}

// defineMetadata is a helper function to avoid panics if no metadata map is defined.
func (p *Pet) defineMetadata() {
	p.Metadata = make(Metadata, 0)
}
