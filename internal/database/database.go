package database

import (
	"log"
	"os"

	"cos-ai-bot/internal/api"
	"cos-ai-bot/internal/models"
)

var apiClient *api.Client

// Маппинг технических кодов в человекочитаемые значения
var valueMapping = map[string]string{
	// Типы кожи
	"skin_dry":       "Сухая",
	"skin_oily":      "Жирная",
	"skin_normal":    "Нормальная",
	"skin_sensitive": "Чувствительная",
	"skin_combined":  "Комбинированная",
	"skin_unknown":   "Не знаю",

	// Возраст
	"age_18_minus": "До 18 лет",
	"age_18_24":    "18-24 года",
	"age_25_34":    "25-34 года",
	"age_35_44":    "35-44 года",
	"age_45_plus":  "45+ лет",
	"age_ignore":   "Не учитывать",

	// Пол
	"gender_male":   "Мужчина",
	"gender_female": "Женщина",
	"gender_other":  "Другое",
	"gender_ignore": "Не учитывать",

	// Беременность/лактация
	"pregnancy":               "Беременность",
	"lactation":               "Лактация",
	"pregnancy_and_lactation": "Беременность и лактация",
	"none_of_above":           "Ничего из перечисленного",
	"pregnancy_ignore":        "Не учитывать",

	// Цели
	"goal_hydration":  "Увлажнение и питание",
	"goal_tone":       "Выравнивание тона",
	"goal_antiage":    "Антивозрастной уход",
	"goal_texture":    "Улучшение текстуры",
	"goal_refresh":    "Освежить и поддерживать",
	"goal_minimalism": "Минимализм, только базовый уход",
	"goal_other":      "Другое",

	// Бюджет
	"budget_low":     "До 1000₽",
	"budget_medium":  "1000-3000₽",
	"budget_high":    "3000-5000₽",
	"budget_premium": "5000₽+",
	"budget_ignore":  "Не важно",

	// Климат
	"climate_dry":       "Сухой",
	"climate_humid":     "Влажный",
	"climate_hot":       "Жаркий",
	"climate_cold":      "Холодный",
	"climate_temperate": "Переменный / умеренный",
	"climate_polluted":  "Загрязнённый (город, смог, пыль)",
	"climate_multiple":  "Живу в нескольких климатах (путешествую/переезды)",
	"climate_unknown":   "Не знаю",

	// Тип кожи по Фицпатрику
	"fitzpatrick_1":       "I – очень светлая, всегда обгорает",
	"fitzpatrick_2":       "II – светлая, обгорает, но может немного загорать",
	"fitzpatrick_3":       "III – светло-смуглая, легко загорает",
	"fitzpatrick_4":       "IV – смуглая, редко обгорает",
	"fitzpatrick_5":       "V – тёмная, почти не обгорает",
	"fitzpatrick_6":       "VI – очень тёмная, никогда не обгорает",
	"fitzpatrick_unknown": "Не знаю / Не хочу указывать",

	// Образ жизни
	"lifestyle_stress":   "Частые стрессы",
	"lifestyle_sleep":    "Недосып / сбитый режим",
	"lifestyle_screen":   "Много экранного времени",
	"lifestyle_sweat":    "Часто потею (спорт, жара и т.д.)",
	"lifestyle_computer": "Работаю за компьютером",
	"lifestyle_active":   "Активно двигаюсь в течение дня",
	"lifestyle_outdoor":  "Регулярно на улице",
	"lifestyle_passive":  "Пассивный / домашний образ жизни",
	"lifestyle_other":    "Другое",

	// Питание
	"diet_vegan":       "Веганство",
	"diet_vegetarian":  "Вегетарианство",
	"diet_halal":       "Халяль",
	"diet_keto":        "Кето / Палео / Низкоуглеводная",
	"diet_gluten_free": "Безглютеновая",
	"diet_no_alcohol":  "Я избегаю спирта в составе",
	"diet_no_animal":   "Я избегаю компонентов животного происхождения",
	"diet_none":        "Нет особых ограничений",
	"diet_other":       "Другое",

	// Аллергии
	"allergies_none":          "Нет аллергий",
	"allergies_nickel":        "Аллергия на никель",
	"allergies_lanolin":       "Аллергия на ланолин",
	"allergies_fragrance":     "Аллергия на отдушки",
	"allergies_preservatives": "Аллергия на консерванты",
	"allergies_other":         "Другое",
}

// convertToHumanReadable преобразует технический код в человекочитаемое значение
func convertToHumanReadable(value string) string {
	if humanReadable, exists := valueMapping[value]; exists {
		return humanReadable
	}
	return value // Если маппинг не найден, возвращаем исходное значение
}

// InitDB инициализирует API клиент
func InitDB() error {
	// Инициализируем API клиент
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080" // значение по умолчанию
	}
	apiClient = api.NewClient(apiURL)

	return nil
}

// CloseDB закрывает подключения
func CloseDB() {
	// API клиент не требует явного закрытия
}

// SaveUserState сохраняет состояние пользователя через API
func SaveUserState(userID int64, state *models.UserState) error {
	// Преобразуем UserState в APIUserProfileUpdate с человекочитаемыми значениями
	profileUpdate := &models.APIUserProfileUpdate{
		SkinType:    convertToHumanReadable(state.SkinType),
		Age:         convertToHumanReadable(state.Age),
		Gender:      convertToHumanReadable(state.Gender),
		Pregnancy:   convertToHumanReadable(state.Pregnancy),
		Concern:     state.Concerns, // уже человекочитаемое (текстовый ввод)
		Goal:        convertToHumanReadable(state.Goal),
		Climate:     convertToHumanReadable(state.Climate),
		Fitzpatrick: convertToHumanReadable(state.Fitzpatrick),
		Lifestyle:   convertToHumanReadable(state.Lifestyle),
		Diet:        convertToHumanReadable(state.Diet),
		Allergy:     convertToHumanReadable(state.Allergies),
	}

	log.Printf("Сохраняем профиль пользователя %d: SkinType='%s', Age='%s', Gender='%s', Pregnancy='%s', Concern='%s', Goal='%s', Climate='%s', Fitzpatrick='%s', Lifestyle='%s', Diet='%s', Allergy='%s'",
		userID, profileUpdate.SkinType, profileUpdate.Age, profileUpdate.Gender, profileUpdate.Pregnancy, profileUpdate.Concern, profileUpdate.Goal, profileUpdate.Climate, profileUpdate.Fitzpatrick, profileUpdate.Lifestyle, profileUpdate.Diet, profileUpdate.Allergy)

	return apiClient.UpdateUserProfile(userID, profileUpdate)
}

// GetUserState получает состояние пользователя через API
func GetUserState(userID int64) (*models.UserState, error) {
	profile, err := apiClient.GetUserProfile(userID)
	if err != nil {
		// Если профиль не найден, возвращаем пустое состояние
		return &models.UserState{Step: 0}, nil
	}

	// Преобразуем APIUserProfile в UserState
	// Для состояния анкеты мы не можем определить текущий шаг из API профиля
	// так как профиль содержит только заполненные данные, а не текущий прогресс
	// Поэтому возвращаем Step = 0, что означает "не в процессе заполнения анкеты"
	state := &models.UserState{
		Step:        0, // не можем определить из API профиля
		SkinType:    profile.SkinType,
		Age:         profile.Age,
		Gender:      profile.Gender,
		Pregnancy:   profile.Pregnancy,
		Concerns:    profile.Concern,
		Goal:        profile.Goal,
		Climate:     profile.Climate,
		Fitzpatrick: profile.Fitzpatrick,
		Lifestyle:   profile.Lifestyle,
		Diet:        profile.Diet,
		Allergies:   profile.Allergy,
	}

	return state, nil
}

// SearchProducts выполняет поиск продуктов через API
func SearchProducts(query string, limit, offset int, brandIDs, ingredientIDs, functionIDs, highlightIDs []int) ([]models.APIProduct, error) {
	return apiClient.SearchProducts(query, limit, offset, brandIDs, ingredientIDs, functionIDs, highlightIDs)
}

// GetProduct получает продукт по ID через API
func GetProduct(id int) (*models.APIProductDetail, error) {
	return apiClient.GetProduct(id)
}

// GetIngredient получает ингредиент по ID через API
func GetIngredient(id int) (*models.APIIngredient, error) {
	return apiClient.GetIngredient(id)
}

// GetUserProducts получает продукты пользователя через API
func GetUserProducts(userID int64) ([]models.APIUserProduct, error) {
	return apiClient.GetUserProducts(userID)
}

// AddUserProduct добавляет продукт пользователю через API
func AddUserProduct(userID int64, productID int) error {
	return apiClient.AddUserProduct(userID, productID)
}

// RemoveUserProduct удаляет продукт из коллекции пользователя через API
func RemoveUserProduct(userID int64, productID int) error {
	return apiClient.RemoveUserProduct(userID, productID)
}

// AddProduct добавляет новый продукт через API
func AddProduct(userID int64, product *models.APIProductCreate) error {
	return apiClient.AddProduct(userID, product)
}

// EmptyUserProfile очищает профиль пользователя через API
func EmptyUserProfile(userID int64) error {
	return apiClient.EmptyUserProfile(userID)
}

// GetUserProfile получает профиль пользователя через API
func GetUserProfile(userID int64) (*models.APIUserProfile, error) {
	log.Printf("Запрашиваем профиль пользователя %d через API", userID)
	profile, err := apiClient.GetUserProfile(userID)
	if err != nil {
		log.Printf("Ошибка API при получении профиля пользователя %d: %v", userID, err)
		return nil, err
	}
	log.Printf("API вернул профиль пользователя %d: %+v", userID, profile)
	return profile, nil
}
