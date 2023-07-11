package config

import (
	"os"
	"testing"
)

func TestReadConfigFromFile(t *testing.T) {
	// Create a sample configuration struct
	type Config struct {
		Foo string
		Bar int
	}
	var cfg Config

	// Create a sample TOML configuration file
	const configStr = `
Foo = "hello"
Bar = 123
`

	// Create a temporary file for the configuration

	file, err := os.CreateTemp("", "config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Write the configuration to the file
	if _, err := file.Write([]byte(configStr)); err != nil {
		t.Fatal(err)
	}

	// Read the configuration from the file
	if err := ReadConfigFromFile(file.Name(), &cfg); err != nil {
		t.Fatal(err)
	}

	// Ensure the configuration was parsed correctly
	if cfg.Foo != "hello" {
		t.Errorf("expected Foo to be 'hello', but got %q", cfg.Foo)
	}
}
