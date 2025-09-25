package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"cos-ai-bot/internal/models"
)

// SendChecklist отправляет чеклист через Telegram API
func SendChecklist(botToken string, chatID int64, checklist models.Checklist) error {
	payload := models.ChecklistPayload{
		ChatID:    chatID,
		Checklist: checklist,
		ParseMode: "HTML",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга JSON: %v", err)
	}

	// Отправляем запрос к Telegram API
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendChecklist", botToken)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка HTTP запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка API: %s", string(body))
	}

	return nil
}

// SendMessageWithChecklist отправляет сообщение с чеклистом
func SendMessageWithChecklist(botToken string, chatID int64, text string, checklist models.Checklist) error {
	payload := models.MessagePayload{
		ChatID: chatID,
		Text:   text,
		ReplyMarkup: models.ReplyMarkup{
			ChecklistOptions: models.ChecklistOptions{
				MinSelected: 1,
				MaxSelected: 5,
				SubmitButton: models.SubmitButton{
					Text:         "Подтвердить выбор",
					CallbackData: "checklist_submit",
				},
			},
		},
		ParseMode: "HTML",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга JSON: %v", err)
	}

	// Отправляем запрос к Telegram API
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка HTTP запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка API: %s", string(body))
	}

	return nil
}
