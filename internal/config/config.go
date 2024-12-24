package config

import (
	"encoding/json"
	"os"
)

func read() (Config, error) {
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

