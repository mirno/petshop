package adapters

import (
	"github.com/mirno/petshop/internal/entities"
	"github.com/mirno/petshop/internal/usecases"
)

var _ usecases.Printer = (*PrinterChain)(nil)

// PrinterChain can be used to combine different printers
type PrinterChain struct {
	printers []usecases.Printer
}

func NewPrinterChain(printers ...usecases.Printer) *PrinterChain {
	return &PrinterChain{printers: printers}
}

// Print loops over the inserted interface implementations
func (chain *PrinterChain) Print(pets []entities.Pet) {
	for _, printer := range chain.printers {
		printer.Print(pets)
	}
}
