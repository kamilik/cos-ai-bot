package main

import (
	"log"

	"cos-ai-bot/internal/bot"
	"cos-ai-bot/internal/config"
	"cos-ai-bot/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Загружаем конфигурацию
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Логируем информацию о конфигурации
	log.Printf("OpenRouter API Key loaded: %d characters", len(cfg.OpenRouterAPIKey))

	// Инициализируем API клиент
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize API client: %v", err)
	}
	defer database.CloseDB()

	// Запускаем бота
	if err := bot.Run(cfg.BotToken, cfg.OpenRouterAPIKey); err != nil {
		log.Fatalf("Failed to run bot: %v", err)
	}
}
