package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Db_url      string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

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

func (config *Config) SetUser(currentUser string) {
	if currentUser == "" {
		fmt.Println("Need a user name")
		return
	} else {
		config.CurrentUser = currentUser
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error retreiving home directory: %v\n", err)
		return
	}

	configFile := homeDir + `/` + `.gatorconfig.json`
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling data:", err)
		return
	}

	err = os.WriteFile(configFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	return
}
