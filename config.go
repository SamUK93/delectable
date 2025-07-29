package main

import (
	"fmt"
	"os"
)

type Config struct {
	GeminiAPIKey string
}

func GetConfig() (Config, error) {
	geminiApiKey := os.Getenv("GEMINI_API_KEY")
	if geminiApiKey == "" {
		return Config{}, fmt.Errorf("failed to set required config GEMINI_API_KEY")
	}

	return Config{
		GeminiAPIKey: geminiApiKey,
	}, nil
}
