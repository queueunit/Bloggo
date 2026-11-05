package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Db_url      string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error retreiving home directory: %v\n", err)
		return Config{}
	}

	configFile := homeDir + `/` + `.gatorconfig.json`
	fileBytes, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println("Could not read config file")
		return Config{}
	}

	var config Config
	if err := json.Unmarshal(fileBytes, &config); err != nil {
		fmt.Println("Could not unmarshal config")
	}

	return config
}

func (config *Config) SetUser(currentUser string) error {
	if currentUser == "" {
		return errors.New("Need 1 arguement (username)")
	} else {
		config.CurrentUser = currentUser
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error retreiving home directory: %v\n", err)
		return err
	}

	configFile := homeDir + `/` + configFileName
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling data:", err)
		return err
	}

	err = os.WriteFile(configFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}
	return errors.New("")
}
