package config

import (
	"encoding/json"
	"fmt"
	"os"
	"errors"
)

type Config struct {
	Database string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (*Config, error) {
	fullpath, err := getConfigFilePath()
	if err != nil {
		return nil, errors.New("unable to get config file path")
	}
	data, err := os.ReadFile(fullpath)
	if err != nil{
		return nil, errors.New("unable to read config file")
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil{
		return nil, errors.New("unable to decode json")
	}
	return &config, nil
}

func (c *Config) SetUser(user string) error{
	c.CurrentUserName = user
	return write(*c)
}

func write(cfg Config) error {
	fullpath, err := getConfigFilePath()
	if err != nil{
		return fmt.Errorf("unable to get config file path: %w", err)
	}
	file, err := os.Create(fullpath)
	if err != nil {
		return fmt.Errorf("unable to create config file: %w", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("unable to encode config to json: %w", err)
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to fetch home directory: %w", err)
	}
	return homeDir + "/" + configFileName, nil
}