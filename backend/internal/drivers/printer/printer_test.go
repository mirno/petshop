package printer

import (
	"testing"

	"github.com/mirno/petshop/internal/entities"
)

func TestConsolePrinterFormatsPriceAndMetadata(t *testing.T) {
	printer := NewConsolePrinter(WithMetadata())
	pet := entities.Pet{
		Name:  "Rex",
		Type:  entities.Dog,
		Price: 80,
		Metadata: entities.Metadata{
			"Region": "Europe",
		},
	}

	got := printer.formatPet(pet)
	want := "Rex (dog): 80.00 [Region: Europe]"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestConsolePrinterOmitsMetadataByDefault(t *testing.T) {
	printer := NewConsolePrinter()
	pet := entities.Pet{
		Name:  "Rex",
		Type:  entities.Dog,
		Price: 80,
		Metadata: entities.Metadata{
			"Region": "Europe",
		},
	}

	got := printer.formatPet(pet)
	want := "Rex (dog): 80.00"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
