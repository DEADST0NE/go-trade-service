package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetConfig() Config {
	var config Config

	// Загрузка конфигурации из файла JSON
	if err := loadConfigFile("config.json", &config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: No .env file found - %v\n", err)
		os.Exit(1)
	}

	// Обновление конфигурации переменными окружения
	if err := updateConfigFromEnv(&config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return config
}

func loadConfigFile(filePath string, config *Config) error {
	configFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening config file: %w", err)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(config); err != nil {
		return fmt.Errorf("error parsing config json: %w", err)
	}

	return nil
}

func updateConfigFromEnv(config *Config) error {
	rabbitUrl, exists := os.LookupEnv("APP_RABBIT_URL")
	if !exists {
		return fmt.Errorf("APP_RABBIT_URL not set in the environment")
	}
	config.Rabbit.Url = rabbitUrl

	return nil
}
