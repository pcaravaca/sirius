package config

import "os"

type Config struct {
	Port           string
	OpenWeatherAPI string
	ExchangeAPI    string
	TasksFile      string
	OpenClawURL    string
	OllamaURL      string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("SIRIUS_PORT", "8080"),
		OpenWeatherAPI: getEnv("OPENWEATHER_API_KEY", ""),
		ExchangeAPI:    getEnv("EXCHANGE_API_KEY", ""),
		TasksFile:      getEnv("SIRIUS_TASKS_FILE", "./data/tasks.json"),
		OpenClawURL:    getEnv("OPENCLAW_URL", "http://localhost:18789"),
		OllamaURL:      getEnv("OLLAMA_URL", "http://192.168.0.100:11434"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
