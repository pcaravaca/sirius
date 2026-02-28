package config

import "os"

type Config struct {
	Port         string
	OpenWeatherAPI string
	ExchangeAPI  string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("SIRIUS_PORT", "8080"),
		OpenWeatherAPI: getEnv("OPENWEATHER_API_KEY", ""),
		ExchangeAPI:   getEnv("EXCHANGE_API_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
