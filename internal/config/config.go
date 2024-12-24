package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Read() (Config, error) {
	configPath, err := configFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	config := Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, nil
	}

	return config, nil
}

func configFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(dir, ConfigFile)

	return configPath, nil
}
