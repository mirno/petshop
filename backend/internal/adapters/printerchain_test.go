package adapters

import (
	"testing"

	"github.com/atoscerebro/eviden-petshop/internal/entities"
)

type recordingPrinter struct {
	calls int
	pets  []entities.Pet
}

func (printer *recordingPrinter) Print(pets []entities.Pet) {
	printer.calls++
	printer.pets = pets
}

func TestPrinterChainPrintsWithEveryPrinter(t *testing.T) {
	first := &recordingPrinter{}
	second := &recordingPrinter{}
	pets := []entities.Pet{{Name: "Rex"}}

	chain := NewPrinterChain(first, second)
	chain.Print(pets)

	if first.calls != 1 || second.calls != 1 {
		t.Fatalf("expected both printers to be called once, got %d and %d", first.calls, second.calls)
	}
	if first.pets[0].Name != "Rex" || second.pets[0].Name != "Rex" {
		t.Fatalf("expected both printers to receive pets")
	}
}
