package printer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/mirno/petshop/internal/entities"
)

type ConsolePrinter struct {
	withMetadata bool
	writer       io.Writer
}

type Option func(*ConsolePrinter)

func NewConsolePrinter(options ...Option) *ConsolePrinter {
	printer := &ConsolePrinter{
		writer: os.Stdout,
	}

	for _, option := range options {
		option(printer)
	}

	return printer
}

func WithMetadata() Option {
	return func(printer *ConsolePrinter) {
		printer.withMetadata = true
	}
}

func WithWriter(writer io.Writer) Option {
	return func(printer *ConsolePrinter) {
		printer.writer = writer
	}
}

func (p *ConsolePrinter) Print(pets []entities.Pet) {
	for _, pet := range pets {
		fmt.Fprintln(p.writer, p.formatPet(pet)) //nolint: errcheck // skipping since I plan to deprecate the Printer in the future.
	}
}

func (p *ConsolePrinter) formatPet(pet entities.Pet) string {
	line := fmt.Sprintf("%s (%s): %.2f", pet.Name, pet.Type, float64(pet.Price))
	if !p.withMetadata || len(pet.Metadata.AdditionalProperties) == 0 {
		return line
	}

	return fmt.Sprintf("%s [%s]", line, formatMetadata(pet.Metadata.AdditionalProperties))
}

func formatMetadata(metadata entities.Metadata) string {
	keys := make([]string, 0, len(metadata))
	for key := range metadata {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	values := make([]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, fmt.Sprintf("%s: %s", key, metadata[key]))
	}

	return strings.Join(values, ", ")
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

	_, err = file.Write(jsonData)
	if err != nil {
		os.Exit(1)
	}
}
