package services

import (
	"fmt"
	"strings"

	"cos-ai-bot/internal/api"
	"cos-ai-bot/internal/database"
	"cos-ai-bot/internal/models"
)

// RecommendationService сервис для работы с рекомендациями
type RecommendationService struct {
	openRouterClient *api.OpenRouterClient
}

// NewRecommendationService создает новый сервис рекомендаций
func NewRecommendationService(openRouterAPIKey string) *RecommendationService {
	fmt.Printf("Инициализация сервиса рекомендаций с API ключом длиной %d символов\n", len(openRouterAPIKey))
	return &RecommendationService{
		openRouterClient: api.NewOpenRouterClient(openRouterAPIKey),
	}
}

// Рекомендации на основе анкеты
func (s *RecommendationService) GetAnketaRecommendations(userID int64) (string, error) {
	// Получаем профиль пользователя
	profile, err := database.GetUserProfile(userID)
	if err != nil {
		return "", fmt.Errorf("ошибка получения профиля: %v", err)
	}

	// Формируем анкету для промпта
	anketaText := s.formatAnketaForPrompt(profile)

	// Создаем промпт
	prompt := fmt.Sprintf(`Ты — профессиональный косметолог и дерматолог-консультант.

На основе анкеты, составь рекомендации:
- Распиши оптимальный план ухода (утро/вечер) с последовательностью применения (step-by-step).
- Учитывай тип кожи, возраст, пол, беременность, аллергию на ингредиенты, климат и цели ухода.
- Если найден потенциальный аллерген в составе — предупреди.
- Рекомендации должны быть понятными и аккуратными, как будто ты — дерматолог-консультант.
- Укажи, какие продукты можно использовать утром, какие вечером, какие через день, какие несовместимы между собой.
- В конце дай 2–3 рекомендации по продуктам, которых явно не хватает.

**Анкета пользователя:**
%s`, anketaText)

	return s.openRouterClient.GetRecommendation(prompt)
}

// GetProductsRecommendations получает рекомендации с учётом продуктов пользователя
func (s *RecommendationService) GetProductsRecommendations(userID int64) (string, error) {
	// Получаем профиль пользователя
	profile, err := database.GetUserProfile(userID)
	if err != nil {
		return "", fmt.Errorf("ошибка получения профиля: %v", err)
	}

	// Получаем продукты пользователя
	products, err := database.GetUserProducts(userID)
	if err != nil {
		return "", fmt.Errorf("ошибка получения продуктов: %v", err)
	}

	// Формируем анкету и продукты для промпта
	anketaText := s.formatAnketaForPrompt(profile)
	productsText := s.formatProductsForPrompt(products)

	// Создаем промпт
	prompt := fmt.Sprintf(`Ты — профессиональный косметолог и дерматолог-консультант.

На основе анкеты и списка косметических средств, составь рекомендации:
- Распиши оптимальный план ухода (утро/вечер) с последовательностью применения (step-by-step).
- Используй только средства, которые уже есть у пользователя, **но укажи, если какого-то этапа не хватает**.
- Учитывай тип кожи, возраст, пол, беременность, аллергию на ингредиенты, климат и цели ухода.
- Если найден потенциальный аллерген в составе — предупреди.
- Рекомендации должны быть понятными и аккуратными, как будто ты — дерматолог-консультант.
- Укажи, какие продукты можно использовать утром, какие вечером, какие через день, какие несовместимы между собой.
- В конце дай 2–3 рекомендации по продуктам, которых явно не хватает.

**Анкета пользователя:**
%s

**Продукты пользователя:**
%s`, anketaText, productsText)

	return s.openRouterClient.GetRecommendation(prompt)
}

// GetGeneralRecommendations получает общие рекомендации
func (s *RecommendationService) GetGeneralRecommendations(userID int64) (string, error) {
	// Получаем профиль пользователя
	profile, err := database.GetUserProfile(userID)
	if err != nil {
		return "", fmt.Errorf("ошибка получения профиля: %v", err)
	}

	// Получаем продукты пользователя
	products, err := database.GetUserProducts(userID)
	if err != nil {
		return "", fmt.Errorf("ошибка получения продуктов: %v", err)
	}

	// Формируем анкету и продукты для промпта
	anketaText := s.formatAnketaForPrompt(profile)
	productsText := s.formatProductsForPrompt(products)

	// Создаем промпт для общих рекомендаций
	prompt := fmt.Sprintf(`Ты — профессиональный косметолог и дерматолог-консультант.

На основе анкеты и списка косметических средств, составь общие рекомендации:
- Проанализируй текущий уход и дай общие советы по улучшению.
- Укажи, какие этапы ухода отсутствуют или недостаточно проработаны.
- Дай рекомендации по изменению образа жизни для улучшения состояния кожи.
- Предложи общие принципы ухода, которые подходят для данного типа кожи и возраста.
- Укажи сезонные особенности ухода.
- Дай советы по питанию и образу жизни для здоровья кожи.

**Анкета пользователя:**
%s

**Продукты пользователя:**
%s`, anketaText, productsText)

	return s.openRouterClient.GetRecommendation(prompt)
}

// formatAnketaForPrompt форматирует анкету для промпта
func (s *RecommendationService) formatAnketaForPrompt(profile *models.APIUserProfile) string {
	var parts []string

	if profile.SkinType != "" {
		parts = append(parts, fmt.Sprintf("Тип кожи: %s", profile.SkinType))
	}
	if profile.Age != "" {
		parts = append(parts, fmt.Sprintf("Возраст: %s", profile.Age))
	}
	if profile.Gender != "" {
		parts = append(parts, fmt.Sprintf("Пол: %s", profile.Gender))
	}
	if profile.Pregnancy != "" {
		parts = append(parts, fmt.Sprintf("Беременность/лактация: %s", profile.Pregnancy))
	}
	if profile.Concern != "" {
		parts = append(parts, fmt.Sprintf("Проблемы: %s", profile.Concern))
	}
	if profile.Goal != "" {
		parts = append(parts, fmt.Sprintf("Цель: %s", profile.Goal))
	}
	if profile.Climate != "" {
		parts = append(parts, fmt.Sprintf("Климат: %s", profile.Climate))
	}
	if profile.Fitzpatrick != "" {
		parts = append(parts, fmt.Sprintf("Тип кожи по Фитцпатрику: %s", profile.Fitzpatrick))
	}
	if profile.Lifestyle != "" {
		parts = append(parts, fmt.Sprintf("Образ жизни: %s", profile.Lifestyle))
	}
	if profile.Diet != "" {
		parts = append(parts, fmt.Sprintf("Питание: %s", profile.Diet))
	}
	if profile.Allergy != "" {
		parts = append(parts, fmt.Sprintf("Аллергии: %s", profile.Allergy))
	}

	return strings.Join(parts, "\n")
}

// formatProductsForPrompt форматирует продукты для промпта
func (s *RecommendationService) formatProductsForPrompt(products []models.APIUserProduct) string {
	if len(products) == 0 {
		return "У пользователя пока нет добавленных продуктов."
	}

	var parts []string
	for _, product := range products {
		parts = append(parts, fmt.Sprintf("- %s (%s)", product.Title, product.Brand))
		if product.Details != "" {
			parts = append(parts, fmt.Sprintf("  Описание: %s", product.Details))
		}
	}

	return strings.Join(parts, "\n")
}
