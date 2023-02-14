package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

func ReadConfigFromFile(filename string, cfg interface{}) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(content, cfg)
	if err != nil {
		return err
	}

	return nil
}

func ReadConfigFromStdin(cfg interface{}) interface{} { return cfg }

func ReadConfigFromArgs(cfg interface{}) interface{} { return cfg }
