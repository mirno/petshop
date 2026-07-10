package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mirno/petshop/internal/adapters"
	"github.com/mirno/petshop/internal/drivers/inmemorykvstore"
	"github.com/mirno/petshop/internal/drivers/printer"
	"github.com/mirno/petshop/internal/entities"
	"github.com/mirno/petshop/internal/testdata"
	"github.com/mirno/petshop/internal/usecases"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	detailsConfigKey    = "details"
	jsonOutputConfigKey = "json-output"
	metadataFieldName   = "Region"
	envPrefix           = "PETSHOP"
)

type Config struct {
	Details    bool
	JSONOutput string
}

type commandRunner struct {
	settings *viper.Viper
	config   Config
}

func main() {
	cmd, err := command()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func command() (*cobra.Command, error) {
	runner := &commandRunner{}

	cmd := &cobra.Command{
		Use:          "petshop",
		Short:        "Print the petshop inventory",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		PreRunE:      runner.preRun,
		RunE:         runner.run,
		PostRunE:     runner.postRun,
	}

	setupFlags(cmd)

	settings, err := setupViper(cmd)
	if err != nil {
		return nil, err
	}
	runner.settings = settings

	return cmd, nil
}

func (runner *commandRunner) preRun(_ *cobra.Command, _ []string) error {
	runner.config = getConfig(runner.settings)
	return nil
}

func (runner *commandRunner) run(cmd *cobra.Command, _ []string) error {
	store := inmemorykvstore.NewInMemoryKVStore[entities.Pet]()
	petshop := usecases.NewPetshop(store, petFixtures()...)
	petshopPrinter := adapters.PetshopPrinter{
		Petshop: petshop,
		Printer: configuredPrinter(runner.config, cmd),
	}

	petshopPrinter.PrintPets()

	return nil
}

func (runner *commandRunner) postRun(_ *cobra.Command, _ []string) error {
	return nil
}

func setupFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(detailsConfigKey, false, "display pet metadata")
	cmd.Flags().String(jsonOutputConfigKey, "", "write pets to a JSON file")
}

func setupViper(cmd *cobra.Command) (*viper.Viper, error) {
	config := viper.New()
	config.SetEnvPrefix(envPrefix)
	config.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	config.AutomaticEnv()
	config.SetDefault(detailsConfigKey, false)
	config.SetDefault(jsonOutputConfigKey, "")

	if err := config.BindPFlag(detailsConfigKey, cmd.Flags().Lookup(detailsConfigKey)); err != nil {
		return nil, fmt.Errorf("bind %s flag: %w", detailsConfigKey, err)
	}
	if err := config.BindPFlag(jsonOutputConfigKey, cmd.Flags().Lookup(jsonOutputConfigKey)); err != nil {
		return nil, fmt.Errorf("bind %s flag: %w", jsonOutputConfigKey, err)
	}

	return config, nil
}

func getConfig(config *viper.Viper) Config {
	return Config{
		Details:    config.GetBool(detailsConfigKey),
		JSONOutput: config.GetString(jsonOutputConfigKey),
	}
}

func configuredPrinter(config Config, cmd *cobra.Command) usecases.Printer {
	options := []printer.Option{
		printer.WithWriter(cmd.OutOrStdout()),
	}
	if config.Details {
		options = append(options, printer.WithMetadata())
	}

	printers := []usecases.Printer{
		printer.NewConsolePrinter(options...),
	}

	if config.JSONOutput != "" {
		printers = append(printers, &printer.JSONPrinter{Path: config.JSONOutput})
	}

	return adapters.NewPrinterChain(printers...)
}

func petFixtures() []entities.Pet {
	pets := make([]entities.Pet, len(testdata.PetFixtures))
	copy(pets, testdata.PetFixtures)

	pets[0].Metadata.Set(metadataFieldName, "Europe")
	pets[1].Metadata.Set(metadataFieldName, "North America")

	return pets
}
