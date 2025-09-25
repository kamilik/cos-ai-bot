package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenRouterClient клиент для работы с OpenRouter API
type OpenRouterClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// OpenRouterRequest структура запроса к OpenRouter
type OpenRouterRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// Message структура сообщения для OpenRouter
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterResponse структура ответа от OpenRouter
type OpenRouterResponse struct {
	Choices []Choice `json:"choices"`
	Error   *Error   `json:"error,omitempty"`
}

// Choice структура выбора из ответа
type Choice struct {
	Message Message `json:"message"`
}

// Error структура ошибки
type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// NewOpenRouterClient создает новый клиент OpenRouter
func NewOpenRouterClient(apiKey string) *OpenRouterClient {
	return &OpenRouterClient{
		apiKey:  apiKey,
		baseURL: "https://openrouter.ai/api/v1",
		client: &http.Client{
			Timeout: 120 * time.Second, // Увеличиваем таймаут до 2 минут
		},
	}
}

// GetRecommendation получает рекомендацию от нейросети
func (c *OpenRouterClient) GetRecommendation(prompt string) (string, error) {
	req := OpenRouterRequest{
		Model: "deepseek/deepseek-r1",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   4000,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("ошибка маршалинга запроса: %v", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://cos-ai-bot.com")
	httpReq.Header.Set("X-Title", "Cos AI Bot")

	// Логируем запрос (без API ключа для безопасности)
	fmt.Printf("OpenRouter запрос: POST %s с API ключом длиной %d символов\n", c.baseURL+"/chat/completions", len(c.apiKey))

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("ошибка OpenRouter API: статус %d, тело: %s", resp.StatusCode, string(body))
	}

	var response OpenRouterResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("ошибка парсинга ответа: %v", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("ошибка OpenRouter: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("пустой ответ от нейросети")
	}

	return response.Choices[0].Message.Content, nil
}
