package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/atoscerebro/eviden-petshop/internal/adapters"
	"github.com/atoscerebro/eviden-petshop/internal/drivers/printer"
	"github.com/atoscerebro/eviden-petshop/internal/entities"
	"github.com/atoscerebro/eviden-petshop/internal/usecases"
	"github.com/google/uuid"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	detailsConfigKey    = "details"
	jsonOutputConfigKey = "json-output"
)

var (
	idRex       = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	idSlytherin = uuid.MustParse("00000000-0000-0000-0000-000000000002")

	petFixtures = []entities.Pet{
		{ID: idRex, Name: "Rex", Type: entities.Dog, Price: 80},
		{ID: idSlytherin, Name: "Slytherin", Type: entities.Snake, Price: 400},
	}
)

func main() {
	if err := configure(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	petFixtures[0].SetRegion("Europe")
	petFixtures[1].SetRegion("North America")

	petshop := usecases.NewPetshop(configuredPrinter(), petFixtures...)
	petshop.PrintPets()
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
