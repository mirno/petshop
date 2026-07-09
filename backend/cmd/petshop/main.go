package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mirno/petshop/internal/adapters"
	"github.com/mirno/petshop/internal/drivers/printer"
	"github.com/mirno/petshop/internal/testdata"
	"github.com/mirno/petshop/internal/usecases"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	detailsConfigKey    = "details"
	jsonOutputConfigKey = "json-output"
)

func main() {
	if err := configure(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	petFixtures := testdata.PetFixtures

	petFixtures[0].SetRegion("Europe")
	petFixtures[1].SetRegion("North America")

	petshop := usecases.NewPetshop(petFixtures...)
	petshopPrinter := adapters.PetshopPrinter{
		Petshop: petshop,
		Printer: configuredPrinter(),
	}

	petshopPrinter.PrintPets()
}

func configure() error {
	pflag.Bool(detailsConfigKey, false, "display pet metadata")
	pflag.String(jsonOutputConfigKey, "", "write pets to a JSON file")
	pflag.Parse()

	viper.SetEnvPrefix("PETSHOP")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	viper.SetDefault(detailsConfigKey, false)
	viper.SetDefault(jsonOutputConfigKey, "")

	if err := viper.BindPFlag(detailsConfigKey, pflag.Lookup(detailsConfigKey)); err != nil {
		return fmt.Errorf("bind %s flag: %w", detailsConfigKey, err)
	}
	if err := viper.BindPFlag(jsonOutputConfigKey, pflag.Lookup(jsonOutputConfigKey)); err != nil {
		return fmt.Errorf("bind %s flag: %w", jsonOutputConfigKey, err)
	}

	return nil
}

func configuredPrinter() usecases.Printer {
	options := make([]printer.Option, 0, 1)
	if viper.GetBool(detailsConfigKey) {
		options = append(options, printer.WithMetadata())
	}

	printers := []usecases.Printer{
		printer.NewConsolePrinter(options...),
	}

	if path := viper.GetString(jsonOutputConfigKey); path != "" {
		printers = append(printers, &printer.JSONPrinter{Path: path})
	}

	return adapters.NewPrinterChain(printers...)
}
