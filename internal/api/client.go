package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cos-ai-bot/internal/models"
)

// Client представляет HTTP клиент для работы с API
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient создает новый API клиент
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchProducts выполняет поиск продуктов
func (c *Client) SearchProducts(query string, limit, offset int, brandIDs, ingredientIDs, functionIDs, highlightIDs []int) ([]models.APIProduct, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))

	if len(brandIDs) > 0 {
		params.Set("brand_ids", formatIntArray(brandIDs))
	}
	if len(ingredientIDs) > 0 {
		params.Set("ingredient_ids", formatIntArray(ingredientIDs))
	}
	if len(functionIDs) > 0 {
		params.Set("function_ids", formatIntArray(functionIDs))
	}
	if len(highlightIDs) > 0 {
		params.Set("highlight_ids", formatIntArray(highlightIDs))
	}

	req, err := http.NewRequest("GET", c.baseURL+"/products/search?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var products []models.APIProduct
	if err := c.doRequest(req, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// GetProduct получает продукт по ID
func (c *Client) GetProduct(id int) (*models.APIProductDetail, error) {
	params := url.Values{}
	params.Set("id", strconv.Itoa(id))

	req, err := http.NewRequest("GET", c.baseURL+"/api/products?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var product models.APIProductDetail
	if err := c.doRequest(req, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

// GetIngredient получает ингредиент по ID
func (c *Client) GetIngredient(id int) (*models.APIIngredient, error) {
	params := url.Values{}
	params.Set("id", strconv.Itoa(id))

	req, err := http.NewRequest("GET", c.baseURL+"/api/ingredients?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var ingredient models.APIIngredient
	if err := c.doRequest(req, &ingredient); err != nil {
		return nil, err
	}

	return &ingredient, nil
}

// GetUserProducts получает продукты пользователя
func (c *Client) GetUserProducts(userID int64) ([]models.APIUserProduct, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/user/products", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("tg-id", strconv.FormatInt(userID, 10))

	var products []models.APIUserProduct
	if err := c.doRequest(req, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// AddUserProduct добавляет продукт пользователю
func (c *Client) AddUserProduct(userID int64, productID int) error {
	params := url.Values{}
	params.Set("id", strconv.Itoa(productID))

	req, err := http.NewRequest("POST", c.baseURL+"/api/user/products?"+params.Encode(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("tg-id", strconv.FormatInt(userID, 10))

	var response map[string]string
	return c.doRequest(req, &response)
}

// RemoveUserProduct удаляет продукт из коллекции пользователя
func (c *Client) RemoveUserProduct(userID int64, productID int) error {
	params := url.Values{}
	params.Set("id", strconv.Itoa(productID))

	req, err := http.NewRequest("DELETE", c.baseURL+"/api/user/products?"+params.Encode(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("tg-id", strconv.FormatInt(userID, 10))

	var response map[string]string
	return c.doRequest(req, &response)
}

// GetUserProfile получает профиль пользователя
func (c *Client) GetUserProfile(userID int64) (*models.APIUserProfile, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/user/profile", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("tg-id", strconv.FormatInt(userID, 10))

	// Логируем запрос
	fmt.Printf("API запрос: GET %s с заголовком tg-id: %d\n", c.baseURL+"/user/profile", userID)

	var profile models.APIUserProfile
	if err := c.doRequest(req, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// UpdateUserProfile обновляет профиль пользователя
func (c *Client) UpdateUserProfile(userID int64, profile *models.APIUserProfileUpdate) error {
	jsonData, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", c.baseURL+"/user/profile", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("tg-id", strconv.FormatInt(userID, 10))
	req.Header.Set("Content-Type", "application/json")

	// Логируем запрос
	fmt.Printf("API запрос: PUT %s с заголовком tg-id: %d, данные: %s\n", c.baseURL+"/user/profile", userID, string(jsonData))

	var response map[string]string
	return c.doRequest(req, &response)
}

// EmptyUserProfile очищает профиль пользователя
func (c *Client) EmptyUserProfile(userID int64) error {
	req, err := http.NewRequest("POST", c.baseURL+"/user/profile/empty", nil)
	if err != nil {
		return err
	}

	req.Header.Set("tg-id", strconv.FormatInt(userID, 10))

	// Логируем запрос
	fmt.Printf("API запрос: POST %s с заголовком tg-id: %d\n", c.baseURL+"/user/profile/empty", userID)

	var response map[string]string
	return c.doRequest(req, &response)
}

// AddProduct добавляет новый продукт
func (c *Client) AddProduct(userID int64, product *models.APIProductCreate) error {
	jsonData, err := json.Marshal(product)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/products", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("tg-id", strconv.FormatInt(userID, 10))
	req.Header.Set("Content-Type", "application/json")

	var response map[string]string
	return c.doRequest(req, &response)
}

// doRequest выполняет HTTP запрос
func (c *Client) doRequest(req *http.Request, result interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("Ошибка HTTP запроса: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка чтения ответа: %v\n", err)
		return err
	}

	fmt.Printf("API ответ: статус %d, тело: %s\n", resp.StatusCode, string(body))

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	if result != nil {
		return json.Unmarshal(body, result)
	}

	return nil
}

// formatIntArray форматирует массив int в строку для API
func formatIntArray(arr []int) string {
	if len(arr) == 0 {
		return "[]"
	}

	strs := make([]string, len(arr))
	for i, v := range arr {
		strs[i] = strconv.Itoa(v)
	}

	return "[" + strings.Join(strs, ",") + "]"
}
