package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mirno/petshop/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandPrintsPetsWithDetailsFlag(t *testing.T) {
	cmd, err := command()
	require.NoError(t, err)

	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--details"})

	err = cmd.Execute()
	require.NoError(t, err)

	assert.Contains(t, output.String(), "Rex (dog): 80.00 [Region: Europe]")
	assert.Contains(t, output.String(), "Slytherin (snake): 400.00 [Region: North America]")
}

func TestCommandReadsDetailsFromEnvironment(t *testing.T) {
	t.Setenv("PETSHOP_DETAILS", "true")

	cmd, err := command()
	require.NoError(t, err)

	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&bytes.Buffer{})

	err = cmd.Execute()
	require.NoError(t, err)

	assert.Contains(t, output.String(), "[Region: Europe]")
	assert.Contains(t, output.String(), "[Region: North America]")
}

func TestCommandWritesJSONOutput(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pets.json")
	cmd, err := command()
	require.NoError(t, err)

	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--json-output", path})

	err = cmd.Execute()
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var pets []entities.Pet
	err = json.Unmarshal(data, &pets)
	require.NoError(t, err)

	require.Len(t, pets, 2)
	assert.Equal(t, "Rex", pets[0].Name)
	assert.Equal(t, "Slytherin", pets[1].Name)
}
