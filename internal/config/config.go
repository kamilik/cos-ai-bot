package config

import (
	"os"
	"strconv"
)

// Config содержит конфигурацию приложения
type Config struct {
	BotToken         string
	DatabaseURL      string
	APIURL           string
	OpenRouterAPIKey string
	Debug            bool
	Port             int
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8080
	}

	return &Config{
		BotToken:         os.Getenv("BOT_TOKEN"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		APIURL:           os.Getenv("API_URL"),
		OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"),
		Debug:            debug,
		Port:             port,
	}
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	if c.BotToken == "" {
		return ErrMissingBotToken
	}
	if c.APIURL == "" {
		return ErrMissingAPIURL
	}
	if c.OpenRouterAPIKey == "" {
		return ErrMissingOpenRouterAPIKey
	}
	return nil
}

// Ошибки конфигурации
var (
	ErrMissingBotToken         = &ConfigError{"BOT_TOKEN не установлен"}
	ErrMissingAPIURL           = &ConfigError{"API_URL не установлен"}
	ErrMissingOpenRouterAPIKey = &ConfigError{"OPENROUTER_API_KEY не установлен"}
)

// ConfigError представляет ошибку конфигурации
type ConfigError struct {
	message string
}

func (e *ConfigError) Error() string {
	return e.message
}
