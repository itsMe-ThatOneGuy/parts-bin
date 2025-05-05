package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func Read() (Config, error) {
	configDirPath, err := configFilePath()
	if err != nil {
		return Config{}, err
	}

	err = os.MkdirAll(configDirPath, 0755)
	if err != nil {
		return Config{}, err
	}

	configFilePath := filepath.Join(configDirPath, ConfigFile)
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		fmt.Println("Creating default config file")
		err = os.WriteFile(configFilePath, []byte(DefaultConfig), 0644)
		if err != nil {
			return Config{}, err
		}

		return Config{}, fmt.Errorf("update the config at '/home/user/.confg/parts-bin'")
	}

	file, err := os.Open(configFilePath)
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

	configPath := filepath.Join(dir, ".config", "parts-bin")

	return configPath, nil
}
