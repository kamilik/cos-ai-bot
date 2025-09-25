package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"cos-ai-bot/internal/database"
	"cos-ai-bot/internal/models"
	"cos-ai-bot/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var userStates = make(map[int64]*models.UserState)        // userID -> состояние анкеты (fallback)
var recommendationService *services.RecommendationService // сервис рекомендаций

// deleteMessage удаляет сообщение
func deleteMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	bot.Send(deleteMsg)
}

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

	// Аллергии
	"allergies_none":          "Нет аллергий",
	"allergies_nickel":        "Аллергия на никель",
	"allergies_lanolin":       "Аллергия на ланолин",
	"allergies_fragrance":     "Аллергия на отдушки",
	"allergies_preservatives": "Аллергия на консерванты",
	"allergies_other":         "Другое",

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
}

// convertToHumanReadable преобразует технический код в человекочитаемое значение
func convertToHumanReadable(value string) string {
	if humanReadable, exists := valueMapping[value]; exists {
		return humanReadable
	}
	return value // Если маппинг не найден, возвращаем исходное значение
}

// formatRecommendationForTelegram форматирует текст рекомендации для красивого отображения в Telegram
func formatRecommendationForTelegram(text string) string {
	// Заменяем Markdown разметку на Telegram форматирование
	formatted := text

	// Сначала заменяем жирный текст, чтобы избежать конфликтов
	formatted = strings.ReplaceAll(formatted, "**", "<b>")

	// Заменяем курсив
	formatted = strings.ReplaceAll(formatted, "*", "<i>")

	// Заменяем заголовки
	formatted = strings.ReplaceAll(formatted, "### ", "🔹 <b>")
	formatted = strings.ReplaceAll(formatted, "#### ", "🔸 <b>")
	formatted = strings.ReplaceAll(formatted, "## ", "🔹 <b>")
	formatted = strings.ReplaceAll(formatted, "# ", "🔹 <b>")

	// Заменяем горизонтальные линии
	formatted = strings.ReplaceAll(formatted, "---", "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	formatted = strings.ReplaceAll(formatted, "--", "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Добавляем эмодзи для списков
	formatted = strings.ReplaceAll(formatted, "1. ", "1️⃣ ")
	formatted = strings.ReplaceAll(formatted, "2. ", "2️⃣ ")
	formatted = strings.ReplaceAll(formatted, "3. ", "3️⃣ ")
	formatted = strings.ReplaceAll(formatted, "4. ", "4️⃣ ")
	formatted = strings.ReplaceAll(formatted, "5. ", "5️⃣ ")
	formatted = strings.ReplaceAll(formatted, "6. ", "6️⃣ ")
	formatted = strings.ReplaceAll(formatted, "7. ", "7️⃣ ")
	formatted = strings.ReplaceAll(formatted, "8. ", "8️⃣ ")
	formatted = strings.ReplaceAll(formatted, "9. ", "9️⃣ ")

	// Добавляем эмодзи для подпунктов
	formatted = strings.ReplaceAll(formatted, "   * ", "   • ")
	formatted = strings.ReplaceAll(formatted, "   - ", "   • ")

	// Добавляем эмодзи для важных слов
	formatted = strings.ReplaceAll(formatted, "Зачем:", "💡 <b>Зачем:</b>")
	formatted = strings.ReplaceAll(formatted, "Средство:", "🧴 <b>Средство:</b>")
	formatted = strings.ReplaceAll(formatted, "Тип кожи:", "👤 <b>Тип кожи:</b>")
	formatted = strings.ReplaceAll(formatted, "Цель:", "🎯 <b>Цель:</b>")

	// Добавляем эмодзи для этапов ухода
	formatted = strings.ReplaceAll(formatted, "Очищение", "🧼 <b>Очищение</b>")
	formatted = strings.ReplaceAll(formatted, "Тоник", "💧 <b>Тоник</b>")
	formatted = strings.ReplaceAll(formatted, "Активный уход", "⚡ <b>Активный уход</b>")
	formatted = strings.ReplaceAll(formatted, "Увлажнение", "💧 <b>Увлажнение</b>")
	formatted = strings.ReplaceAll(formatted, "Солнцезащита", "☀️ <b>Солнцезащита</b>")
	formatted = strings.ReplaceAll(formatted, "Двойное очищение", "🔄 <b>Двойное очищение</b>")

	// Добавляем эмодзи для времени
	formatted = strings.ReplaceAll(formatted, "Утренний уход", "🌅 <b>Утренний уход</b>")
	formatted = strings.ReplaceAll(formatted, "Вечерний уход", "🌙 <b>Вечерний уход</b>")

	// Добавляем эмодзи для ингредиентов
	formatted = strings.ReplaceAll(formatted, "салициловая кислота", "<b>салициловая кислота</b>")
	formatted = strings.ReplaceAll(formatted, "ниацинамид", "<b>ниацинамид</b>")
	formatted = strings.ReplaceAll(formatted, "гиалуроновая кислота", "<b>гиалуроновая кислота</b>")
	formatted = strings.ReplaceAll(formatted, "азулаиновая кислота", "<b>азулаиновая кислота</b>")
	formatted = strings.ReplaceAll(formatted, "оксид цинка", "<b>оксид цинка</b>")

	// Исправляем двойные теги
	formatted = strings.ReplaceAll(formatted, "<b><b>", "<b>")
	formatted = strings.ReplaceAll(formatted, "</b></b>", "</b>")
	formatted = strings.ReplaceAll(formatted, "<i><i>", "<i>")
	formatted = strings.ReplaceAll(formatted, "</i></i>", "</i>")

	// Валидация и исправление HTML тегов
	formatted = validateAndFixHTMLTags(formatted)

	return formatted
}

// validateAndFixHTMLTags проверяет и исправляет HTML теги
func validateAndFixHTMLTags(text string) string {
	// Подсчитываем открывающие и закрывающие теги
	bOpenCount := strings.Count(text, "<b>")
	bCloseCount := strings.Count(text, "</b>")
	iOpenCount := strings.Count(text, "<i>")
	iCloseCount := strings.Count(text, "</i>")

	// Добавляем недостающие закрывающие теги
	if bOpenCount > bCloseCount {
		text += strings.Repeat("</b>", bOpenCount-bCloseCount)
	}
	if iOpenCount > iCloseCount {
		text += strings.Repeat("</i>", iOpenCount-iCloseCount)
	}

	return text
}

// Run запускает бота
func Run(token, openRouterAPIKey string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("ошибка создания бота: %v", err)
	}

	bot.Debug = true
	log.Printf("Бот запущен: %s", bot.Self.UserName)

	// Инициализируем сервис рекомендаций
	recommendationService = services.NewRecommendationService(openRouterAPIKey)
	log.Printf("Сервис рекомендаций инициализирован")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil && update.InlineQuery == nil {
			continue
		}

		if update.Message != nil {
			handleMessage(bot, update.Message)
		}

		if update.CallbackQuery != nil {
			handleCallbackQuery(bot, update.CallbackQuery)
		}

		if update.InlineQuery != nil {
			handleInlineQuery(bot, update.InlineQuery)
		}
	}

	return nil
}

// handleMessage обрабатывает входящие сообщения
func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	text := message.Text

	log.Printf("[%s] %s", message.From.UserName, text)

	// Обработка команд
	if message.IsCommand() {
		handleCommand(bot, message)
		return
	}

	// Обработка парсинга URL
	if strings.Contains(text, "incidecoder.com") {
		handleIncidecoderURL(bot, message)
		return
	}

	// Обработка формы
	handleFormInput(bot, message)
}

// handleCommand обрабатывает команды
func handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	command := message.Command()

	switch command {
	case "start":
		// Отправляем фото с приветственным сообщением
		photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("images/01.png"))
		photo.Caption = `✨ Я — твой умный бьюти-бот, созданный, чтобы наконец навести порядок в косметичке. Этот бот - часть проекта Cos AI, созданного для того, чтобы помочь тебе собрать персонализированный уход за кожей.
Хочешь попробовать? Давай начнем с небольшой анкеты 💬👇`
		photo.ParseMode = "HTML"

		// Создаем клавиатуру с кнопками
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📋 Анкета", "start_form"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🤖 Рекомендации", "recommendations"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🧴 Мои продукты", "my_products"),
			),
		)
		photo.ReplyMarkup = keyboard
		bot.Send(photo)

	case "help":
		helpText := `Доступные команды:
/start - Начать работу с ботом
/help - Показать эту справку
/form - Заполнить форму подбора ухода
/myproducts - Показать мои продукты

🔍 Для поиска продуктов используйте inline режим:
@cosmetics_lab_ai_bot add [название продукта]`
		msg := tgbotapi.NewMessage(chatID, helpText)
		bot.Send(msg)

	case "form":
		// Инициализируем новое состояние для анкеты
		newState := &models.UserState{Step: 1}
		saveUserState(chatID, newState)
		ShowSkincareFormStep(bot, chatID, 1)

	case "myproducts":
		// Показываем продукты пользователя
		handleMyProductsCommand(bot, message)

	default:
		msg := tgbotapi.NewMessage(chatID, "Неизвестная команда. Используйте /help для справки.")
		bot.Send(msg)
	}
}

// handleCallbackQuery обрабатывает нажатия на inline кнопки
func handleCallbackQuery(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	log.Printf("[CALLBACK] %s: %s", callback.From.UserName, data)

	// Удаляем предыдущее сообщение с кнопками
	deleteMessage(bot, chatID, callback.Message.MessageID)

	// Отвечаем на callback
	callbackAnswer := tgbotapi.NewCallback(callback.ID, "")
	bot.Request(callbackAnswer)

	switch {
	case data == "start_form":
		handleAnketa(bot, callback)

	case strings.HasPrefix(data, "skin_") || strings.HasPrefix(data, "age_") ||
		strings.HasPrefix(data, "gender_") || strings.HasPrefix(data, "pregnancy_") ||
		strings.HasPrefix(data, "goal_") || strings.HasPrefix(data, "climate_") ||
		strings.HasPrefix(data, "fitzpatrick_") || strings.HasPrefix(data, "lifestyle_") ||
		strings.HasPrefix(data, "diet_") || strings.HasPrefix(data, "allergies_"):
		handleFormCallback(bot, callback)

	case strings.HasPrefix(data, "product_"):
		handleProductSelection(bot, callback)

	case strings.HasPrefix(data, "add_product_"):
		handleAddProductToCollection(bot, callback)

	case strings.HasPrefix(data, "remove_product_"):
		handleRemoveProductFromCollection(bot, callback)

	case data == "recommendations":
		handleRecommendations(bot, callback)

	case data == "recommendations_anketa":
		handleRecommendationsAnketa(bot, callback)

	case data == "recommendations_products":
		handleRecommendationsProducts(bot, callback)

	case data == "recommendations_general":
		handleRecommendationsGeneral(bot, callback)

	case data == "my_products":
		handleMyProducts(bot, callback)

	case data == "delete_products":
		handleDeleteProducts(bot, callback)

	case data == "delete_anketa":
		handleDeleteAnketa(bot, callback)

	case data == "retake_anketa":
		// Очищаем локальное состояние при начале анкеты заново
		delete(userStates, chatID)
		// Инициализируем новое состояние для анкеты
		newState := &models.UserState{Step: 1}
		saveUserState(chatID, newState)
		ShowSkincareFormStep(bot, chatID, 1)

	case data == "start_form_new":
		// Инициализируем новое состояние для анкеты
		newState := &models.UserState{Step: 1}
		saveUserState(chatID, newState)
		ShowSkincareFormStep(bot, chatID, 1)

	case data == "back_to_start":
		handleBackToStart(bot, callback)

	default:
		log.Printf("Неизвестный callback: %s", data)
	}
}

// handleFormCallback обрабатывает ответы на форму
func handleFormCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	log.Printf("Обработка callback для пользователя %d: %s", chatID, data)

	// Получаем текущее состояние пользователя (с fallback)
	state := getUserState(chatID)
	log.Printf("Текущее состояние пользователя %d: шаг %d", chatID, state.Step)

	// Обновляем состояние в зависимости от ответа
	switch {
	case strings.HasPrefix(data, "skin_"):
		state.SkinType = data
		state.Step = 2
		log.Printf("Пользователь %d выбрал тип кожи: %s, переходим к шагу 2", chatID, data)
	case strings.HasPrefix(data, "age_"):
		state.Age = data
		state.Step = 3
		log.Printf("Пользователь %d выбрал возраст: %s, переходим к шагу 3", chatID, data)
	case strings.HasPrefix(data, "gender_"):
		state.Gender = data
		// Если выбран мужской пол, автоматически пропускаем вопрос о беременности
		if data == "gender_male" {
			state.Pregnancy = "none_of_above"
			state.Step = 5 // переходим сразу к вопросу о проблемах
			log.Printf("Пользователь %d выбрал пол: %s, автоматически пропускаем беременность, переходим к шагу 5", chatID, data)
		} else {
			state.Step = 4 // переходим к вопросу о беременности
			log.Printf("Пользователь %d выбрал пол: %s, переходим к шагу 4", chatID, data)
		}
	case strings.HasPrefix(data, "pregnancy_"):
		state.Pregnancy = data
		state.Step = 5
		log.Printf("Пользователь %d выбрал беременность: %s, переходим к шагу 5", chatID, data)
	case strings.HasPrefix(data, "goal_"):
		state.Goal = data
		state.Step = 7
		log.Printf("Пользователь %d выбрал цель: %s, переходим к шагу 7", chatID, data)
	case strings.HasPrefix(data, "climate_"):
		state.Climate = data
		state.Step = 8
		log.Printf("Пользователь %d выбрал климат: %s, переходим к шагу 8", chatID, data)
	case strings.HasPrefix(data, "fitzpatrick_"):
		state.Fitzpatrick = data
		state.Step = 9
		log.Printf("Пользователь %d выбрал тип кожи по Фицпатрику: %s, переходим к шагу 9", chatID, data)
	case strings.HasPrefix(data, "lifestyle_"):
		state.Lifestyle = data
		state.Step = 10
		log.Printf("Пользователь %d выбрал образ жизни: %s, переходим к шагу 10", chatID, data)
	case strings.HasPrefix(data, "diet_"):
		state.Diet = data
		state.Step = 11
		log.Printf("Пользователь %d выбрал питание: %s, переходим к шагу 11", chatID, data)
	case strings.HasPrefix(data, "allergies_"):
		state.Allergies = data
		state.Step = 12
		log.Printf("Пользователь %d выбрал аллергии: %s, завершаем форму", chatID, data)
		// Форма завершена, показываем результаты
		showFormResults(bot, callback.Message, state)
		return
	default:
		log.Printf("Неизвестный callback data: %s", data)
		return
	}

	// Сохраняем состояние (с fallback)
	saveUserState(chatID, state)
	log.Printf("Состояние пользователя %d сохранено: шаг %d", chatID, state.Step)

	// Показываем следующий шаг
	ShowSkincareFormStep(bot, chatID, state.Step)
}

// handleFormInput обрабатывает текстовые ответы на форму
func handleFormInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := message.Text

	log.Printf("Обработка текстового ввода для пользователя %d: '%s'", chatID, text)

	// Получаем текущее состояние пользователя (с fallback)
	state := getUserState(chatID)
	log.Printf("Текущее состояние пользователя %d: шаг %d, Concerns: '%s'", chatID, state.Step, state.Concerns)

	// Проверяем, что мы действительно в процессе заполнения формы
	if state.Step == 0 {
		// Не в процессе заполнения формы
		log.Printf("Пользователь %d не в процессе заполнения формы (шаг 0)", chatID)
		return
	}

	// Обрабатываем текстовые ответы
	switch state.Step {
	case 5:
		state.Concerns = text
		state.Step = 6
		log.Printf("Пользователь %d ввел проблемы: %s, переходим к шагу 6", chatID, text)
	default:
		// Не ожидаем текстовый ввод на этом шаге
		log.Printf("Пользователь %d отправил текст на шаге %d, но это не ожидается", chatID, state.Step)
		return
	}

	// Сохраняем состояние (с fallback)
	saveUserState(chatID, state)
	log.Printf("Состояние пользователя %d сохранено: шаг %d", chatID, state.Step)

	// Показываем следующий шаг
	log.Printf("Показываем шаг %d для пользователя %d", state.Step, chatID)
	ShowSkincareFormStep(bot, chatID, state.Step)
}

// showFormResults показывает результаты заполнения формы
func showFormResults(bot *tgbotapi.BotAPI, message *tgbotapi.Message, state *models.UserState) {
	chatID := message.Chat.ID

	resultText := fmt.Sprintf(`✅ Форма заполнена! Вот ваши данные:

👤 Тип кожи: %s
📅 Возраст: %s
🚻 Пол: %s
🤱 Беременность/лактация: %s
💭 Проблемы: %s
🎯 Цель: %s
🌍 Климат: %s
☀️ Тип кожи по Фицпатрику: %s
🏃 Образ жизни: %s
🥗 Питание: %s
⚠️ Аллергии: %s

Теперь я могу подобрать для вас подходящие средства!`,
		convertToHumanReadable(state.SkinType),
		convertToHumanReadable(state.Age),
		convertToHumanReadable(state.Gender),
		convertToHumanReadable(state.Pregnancy),
		state.Concerns, // уже человекочитаемое (текстовый ввод)
		convertToHumanReadable(state.Goal),
		convertToHumanReadable(state.Climate),
		convertToHumanReadable(state.Fitzpatrick),
		convertToHumanReadable(state.Lifestyle),
		convertToHumanReadable(state.Diet),
		convertToHumanReadable(state.Allergies))

	// Отправляем фото с результатами анкеты
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("images/12.png"))
	photo.Caption = resultText
	photo.ParseMode = "HTML"

	// Создаем клавиатуру с кнопками
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить анкету", "delete_anketa"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Пройти заново", "retake_anketa"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_to_start"),
		),
	)
	photo.ReplyMarkup = keyboard
	bot.Send(photo)

	// Сохраняем финальное состояние анкеты в API
	log.Printf("Сохраняем финальное состояние анкеты пользователя %d", chatID)
	saveUserState(chatID, state)
}

// handleIncidecoderURL обрабатывает URL с Incidecoder
func handleIncidecoderURL(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID

	msg := tgbotapi.NewMessage(chatID, "Парсинг продукта с Incidecoder...")
	bot.Send(msg)

	// TODO: Реализовать парсинг через API или создать отдельный сервис
	// Пока что просто сообщаем, что функция в разработке
	infoMsg := tgbotapi.NewMessage(chatID, "Функция парсинга продуктов с Incidecoder будет реализована в следующих версиях.")
	bot.Send(infoMsg)
}

// handleProductSelection обрабатывает выбор продукта
func handleProductSelection(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	// Извлекаем ID продукта
	productIDStr := strings.TrimPrefix(data, "product_")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		errorMsg := tgbotapi.NewMessage(chatID, "Ошибка: неверный ID продукта")
		bot.Send(errorMsg)
		return
	}

	// Получаем детальную информацию о продукте через API
	product, err := database.GetProduct(productID)
	if err != nil {
		errorMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ошибка получения продукта: %v", err))
		bot.Send(errorMsg)
		return
	}

	// Формируем сообщение с информацией о продукте
	var productText strings.Builder
	productText.WriteString(fmt.Sprintf("🧴 <b>%s %s</b>\n\n", product.Brand, product.Title))

	if product.Details != "" {
		productText.WriteString(fmt.Sprintf("📝 <b>Описание:</b>\n%s\n\n", product.Details))
	}

	if len(product.Ingredients) > 0 {
		productText.WriteString("🧪 <b>Ингредиенты:</b>\n")
		for i, ingredient := range product.Ingredients {
			if i >= 10 { // Ограничиваем количество ингредиентов
				productText.WriteString(fmt.Sprintf("... и еще %d ингредиентов", len(product.Ingredients)-10))
				break
			}
			productText.WriteString(fmt.Sprintf("• %s\n", ingredient.Name))
		}
	}

	// Создаем клавиатуру с действиями
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Добавить в коллекцию", fmt.Sprintf("add_product_%d", product.ID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, productText.String())
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

// handleAddProductToCollection обрабатывает добавление продукта в коллекцию пользователя
func handleAddProductToCollection(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	// Извлекаем ID продукта
	productIDStr := strings.TrimPrefix(data, "add_product_")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		errorMsg := tgbotapi.NewMessage(chatID, "Ошибка: неверный ID продукта")
		bot.Send(errorMsg)
		return
	}

	// Добавляем продукт в коллекцию пользователя через API
	err = database.AddUserProduct(chatID, productID)
	if err != nil {
		errorMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ошибка добавления продукта: %v", err))
		bot.Send(errorMsg)
		return
	}

	successMsg := tgbotapi.NewMessage(chatID, "✅ Продукт успешно добавлен в вашу коллекцию!")
	bot.Send(successMsg)
}

// handleRemoveProductFromCollection обрабатывает удаление продукта из коллекции пользователя
func handleRemoveProductFromCollection(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	// Извлекаем ID продукта
	productIDStr := strings.TrimPrefix(data, "remove_product_")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		errorMsg := tgbotapi.NewMessage(chatID, "❌ Ошибка: неверный ID продукта")
		bot.Send(errorMsg)
		return
	}

	// Удаляем продукт из коллекции пользователя через API
	err = database.RemoveUserProduct(chatID, productID)
	if err != nil {
		errorMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Ошибка удаления продукта: %v", err))
		bot.Send(errorMsg)
		return
	}

	successMsg := tgbotapi.NewMessage(chatID, "✅ Продукт успешно удален из вашей коллекции!")
	bot.Send(successMsg)

	// Показываем обновленный список продуктов
	handleMyProducts(bot, callback)
}

// handleRecommendations обрабатывает запрос рекомендаций
func handleRecommendations(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	// Отправляем фото с подписью
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("images/04.png"))
	photo.Caption = `🤖 <b>Рекомендации</b>

Мы можем предложить тебе советы по уходу на основе анкеты, твоих текущих продуктов, или подсказать, чего не хватает в твоем уходе.

<i>Данные рекомендации только для ознакомления и не заменяют консультацию дерматолога и не ставят точные диагнозы.</i>`
	photo.ParseMode = "HTML"

	// Создаем клавиатуру с кнопками рекомендаций
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Рекомендации на основе анкеты", "recommendations_anketa"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧴 Рекомендации с учётом моих продуктов", "recommendations_products"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧩 Общие рекомендации", "recommendations_general"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_to_start"),
		),
	)
	photo.ReplyMarkup = keyboard
	bot.Send(photo)
}

// handleRecommendationsAnketa обрабатывает рекомендации на основе анкеты
func handleRecommendationsAnketa(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	// Отправляем сообщение о загрузке
	loadingMsg := tgbotapi.NewMessage(chatID, "🤖 Генерирую рекомендации на основе вашей анкеты...\n\n⏳ Это может занять до 2 минут. Пожалуйста, подождите...")
	bot.Send(loadingMsg)

	// Получаем рекомендации
	recommendations, err := recommendationService.GetAnketaRecommendations(chatID)
	if err != nil {
		var errorText string
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline exceeded") {
			errorText = "⏰ Время ожидания истекло. Нейросеть работает медленно. Попробуйте еще раз через несколько минут."
		} else if strings.Contains(err.Error(), "401") {
			errorText = "🔑 Ошибка аутентификации. Проверьте настройки API."
		} else {
			errorText = fmt.Sprintf("❌ Ошибка получения рекомендаций: %v", err)
		}
		errorMsg := tgbotapi.NewMessage(chatID, errorText)

		// Добавляем кнопку "Попробовать снова" для ошибок таймаута
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline exceeded") {
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🔄 Попробовать снова", "recommendations_anketa"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к рекомендациям", "recommendations"),
				),
			)
			errorMsg.ReplyMarkup = keyboard
		} else {
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к рекомендациям", "recommendations"),
				),
			)
			errorMsg.ReplyMarkup = keyboard
		}

		bot.Send(errorMsg)
		return
	}

	// Отправляем рекомендации с красивым форматированием
	formattedRecommendations := formatRecommendationForTelegram(recommendations)
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("📊 <b>Рекомендации на основе анкеты</b>\n\n%s", formattedRecommendations))
	msg.ParseMode = "HTML"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к рекомендациям", "recommendations"),
		),
	)
	msg.ReplyMarkup = keyboard

	// Если отправка с HTML не удалась, отправляем без форматирования
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки с HTML форматированием: %v", err)
		// Отправляем без форматирования
		plainMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("📊 Рекомендации на основе анкеты\n\n%s", recommendations))
		plainMsg.ReplyMarkup = keyboard
		bot.Send(plainMsg)
	}
}

// handleRecommendationsProducts обрабатывает рекомендации с учётом продуктов
func handleRecommendationsProducts(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	// Отправляем сообщение о загрузке
	loadingMsg := tgbotapi.NewMessage(chatID, "🤖 Генерирую рекомендации с учётом ваших продуктов...\n\n⏳ Это может занять до 2 минут. Пожалуйста, подождите...")
	bot.Send(loadingMsg)

	// Получаем рекомендации
	recommendations, err := recommendationService.GetProductsRecommendations(chatID)
	if err != nil {
		errorMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Ошибка получения рекомендаций: %v", err))
		bot.Send(errorMsg)
		return
	}

	// Отправляем рекомендации с красивым форматированием
	formattedRecommendations := formatRecommendationForTelegram(recommendations)
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("🧴 <b>Рекомендации с учётом моих продуктов</b>\n\n%s", formattedRecommendations))
	msg.ParseMode = "HTML"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к рекомендациям", "recommendations"),
		),
	)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

// handleRecommendationsGeneral обрабатывает общие рекомендации
func handleRecommendationsGeneral(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	// Отправляем сообщение о загрузке
	loadingMsg := tgbotapi.NewMessage(chatID, "🤖 Генерирую общие рекомендации...\n\n⏳ Это может занять до 2 минут. Пожалуйста, подождите...")
	bot.Send(loadingMsg)

	// Получаем рекомендации
	recommendations, err := recommendationService.GetGeneralRecommendations(chatID)
	if err != nil {
		errorMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Ошибка получения рекомендаций: %v", err))
		bot.Send(errorMsg)
		return
	}

	// Отправляем рекомендации с красивым форматированием
	formattedRecommendations := formatRecommendationForTelegram(recommendations)
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("🧩 <b>Общие рекомендации</b>\n\n%s", formattedRecommendations))
	msg.ParseMode = "HTML"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к рекомендациям", "recommendations"),
		),
	)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

// handleMyProducts обрабатывает запрос на просмотр продуктов пользователя
func handleMyProducts(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	// Отправляем сообщение о загрузке
	loadingMsg := tgbotapi.NewMessage(chatID, "🔄 Загружаю ваши продукты...")
	bot.Send(loadingMsg)

	// Получаем продукты пользователя через API
	products, err := database.GetUserProducts(chatID)
	if err != nil {
		errorMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Ошибка получения ваших продуктов: %v", err))
		bot.Send(errorMsg)
		return
	}

	if len(products) == 0 {
		// Отправляем фото с сообщением об отсутствии продуктов
		photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("images/08.png"))
		photo.Caption = `🧴 <b>Ваша коллекция пуста</b>

У вас пока нет добавленных продуктов в коллекцию.

<b>Для поиска продуктов введите:</b>
@cosmetics_lab_ai_bot add [продукт который хотите найти]

<b>Пример:</b>
@cosmetics_lab_ai_bot add Repair Sunscreen SPF 50`
		photo.ParseMode = "HTML"

		// Добавляем только кнопку "Назад"
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_to_start"),
			),
		)
		photo.ReplyMarkup = keyboard
		bot.Send(photo)
		return
	}

	// Формируем красивое сообщение со списком продуктов
	var productsText strings.Builder
	productsText.WriteString(fmt.Sprintf("🧴 <b>Ваша коллекция (%d продуктов)</b>\n\n", len(products)))
	productsText.WriteString("💡 <b>Для добавления новых продуктов введите:</b>\n@cosmetics_lab_ai_bot add [название продукта]\n\n")

	// Показываем первые 10 продуктов с подробной информацией
	for i, product := range products {
		if i >= 10 {
			productsText.WriteString(fmt.Sprintf("... и еще %d продуктов\n", len(products)-10))
			break
		}

		productsText.WriteString(fmt.Sprintf("🔸 <b>%s %s</b>\n", product.Brand, product.Title))

		// Добавляем описание, если есть
		if product.Details != "" {
			// Обрезаем описание если оно слишком длинное
			details := product.Details
			if len(details) > 100 {
				details = details[:97] + "..."
			}
			productsText.WriteString(fmt.Sprintf("   📝 %s\n", details))
		}

		// Добавляем дату добавления
		if product.AddedAt != "" {
			productsText.WriteString(fmt.Sprintf("   📅 Добавлено: %s\n", product.AddedAt))
		}

		productsText.WriteString("\n")
	}

	// Отправляем фото с подписью
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("images/02.png"))
	photo.Caption = productsText.String()
	photo.ParseMode = "HTML"

	// Создаем клавиатуру с действиями
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить продукты", "delete_products"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_to_start"),
		),
	)
	photo.ReplyMarkup = keyboard
	bot.Send(photo)
}

// handleMyProductsCommand обрабатывает команду /myproducts
func handleMyProductsCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID

	// Отправляем сообщение о загрузке
	loadingMsg := tgbotapi.NewMessage(chatID, "🔄 Загружаю ваши продукты...")
	bot.Send(loadingMsg)

	// Получаем продукты пользователя через API
	products, err := database.GetUserProducts(chatID)
	if err != nil {
		errorMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Ошибка получения ваших продуктов: %v", err))
		bot.Send(errorMsg)
		return
	}

	if len(products) == 0 {
		// Отправляем фото с сообщением об отсутствии продуктов
		photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("images/08.png"))
		photo.Caption = `🧴 <b>Ваша коллекция пуста</b>

У вас пока нет добавленных продуктов в коллекцию.

<b>Для поиска продуктов введите:</b>
@cosmetics_lab_ai_bot add [продукт который хотите найти]

<b>Пример:</b>
@cosmetics_lab_ai_bot add Repair Sunscreen SPF 50`
		photo.ParseMode = "HTML"
		bot.Send(photo)
		return
	}

	// Формируем красивое сообщение со списком продуктов
	var productsText strings.Builder
	productsText.WriteString(fmt.Sprintf("🧴 <b>Ваша коллекция (%d продуктов)</b>\n\n", len(products)))
	productsText.WriteString("💡 <b>Для добавления новых продуктов введите:</b>\n@cosmetics_lab_ai_bot add [название продукта]\n\n")

	// Показываем первые 10 продуктов с подробной информацией
	for i, product := range products {
		if i >= 10 {
			productsText.WriteString(fmt.Sprintf("... и еще %d продуктов\n", len(products)-10))
			break
		}

		productsText.WriteString(fmt.Sprintf("🔸 <b>%s %s</b>\n", product.Brand, product.Title))

		// Добавляем описание, если есть
		if product.Details != "" {
			// Обрезаем описание если оно слишком длинное
			details := product.Details
			if len(details) > 100 {
				details = details[:97] + "..."
			}
			productsText.WriteString(fmt.Sprintf("   📝 %s\n", details))
		}

		// Добавляем дату добавления
		if product.AddedAt != "" {
			productsText.WriteString(fmt.Sprintf("   📅 Добавлено: %s\n", product.AddedAt))
		}

		productsText.WriteString("\n")
	}

	// Отправляем фото с подписью
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("images/02.png"))
	photo.Caption = productsText.String()
	photo.ParseMode = "HTML"

	// Создаем клавиатуру с действиями
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить продукты", "delete_products"),
		),
	)
	photo.ReplyMarkup = keyboard
	bot.Send(photo)
}

// handleDeleteProducts обрабатывает удаление продуктов из коллекции
func handleDeleteProducts(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	// Получаем продукты пользователя через API
	products, err := database.GetUserProducts(chatID)
	if err != nil {
		errorMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Ошибка получения ваших продуктов: %v", err))
		bot.Send(errorMsg)
		return
	}

	if len(products) == 0 {
		msg := tgbotapi.NewMessage(chatID, "🧴 У вас нет продуктов для удаления.")
		bot.Send(msg)
		return
	}

	// Формируем сообщение со списком продуктов для удаления
	var productsText strings.Builder
	productsText.WriteString(fmt.Sprintf("🗑️ <b>Выберите продукты для удаления (%d продуктов):</b>\n\n", len(products)))

	// Показываем первые 10 продуктов
	for i, product := range products {
		if i >= 10 {
			productsText.WriteString(fmt.Sprintf("... и еще %d продуктов\n", len(products)-10))
			break
		}

		productsText.WriteString(fmt.Sprintf("🔸 <b>%s %s</b>\n", product.Brand, product.Title))
	}

	msg := tgbotapi.NewMessage(chatID, productsText.String())
	msg.ParseMode = "HTML"

	// Создаем клавиатуру с кнопками для удаления продуктов
	var keyboard [][]tgbotapi.InlineKeyboardButton
	for i, product := range products {
		if i >= 10 { // Ограничиваем количество кнопок
			break
		}
		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("🗑️ %s %s", product.Brand, product.Title),
			fmt.Sprintf("remove_product_%d", product.ProductID),
		)
		keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(button))
	}

	// Добавляем кнопку "Назад"
	keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к продуктам", "my_products"),
	))

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	bot.Send(msg)
}

// handleAnketa обрабатывает кнопку "Анкета"
func handleAnketa(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	log.Printf("Проверяем анкету пользователя %d через API", chatID)

	// Получаем профиль пользователя через API
	profile, err := database.GetUserProfile(chatID)
	if err != nil {
		log.Printf("Ошибка получения профиля пользователя %d: %v", chatID, err)
		// Если профиль не найден, показываем кнопки для прохождения анкеты
		msg := tgbotapi.NewMessage(chatID, "📋 У вас пока нет заполненной анкеты.\n\nЗаполните анкету, чтобы получить персонализированные рекомендации по уходу за кожей!")

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📝 Пройти анкету", "start_form_new"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_to_start"),
			),
		)
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
		return
	}

	// Логируем полученные данные профиля
	log.Printf("Получен профиль пользователя %d: SkinType='%s', Age='%s', Gender='%s', Pregnancy='%s', Concern='%s', Goal='%s', Climate='%s', Fitzpatrick='%s', Lifestyle='%s', Diet='%s', Allergy='%s'",
		chatID, profile.SkinType, profile.Age, profile.Gender, profile.Pregnancy, profile.Concern, profile.Goal, profile.Climate, profile.Fitzpatrick, profile.Lifestyle, profile.Diet, profile.Allergy)

	// Проверяем, заполнена ли анкета (есть ли хотя бы одно поле)
	if profile.SkinType == "" && profile.Age == "" && profile.Gender == "" &&
		profile.Pregnancy == "" && profile.Concern == "" && profile.Goal == "" &&
		profile.Climate == "" && profile.Fitzpatrick == "" && profile.Lifestyle == "" &&
		profile.Diet == "" && profile.Allergy == "" {
		// Анкета пустая
		log.Printf("Анкета пользователя %d пустая", chatID)
		msg := tgbotapi.NewMessage(chatID, "📋 У вас пока нет заполненной анкеты.\n\nЗаполните анкету, чтобы получить персонализированные рекомендации по уходу за кожей!")

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📝 Пройти анкету", "start_form_new"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_to_start"),
			),
		)
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
		return
	}

	// Анкета заполнена, показываем её содержимое
	log.Printf("Анкета пользователя %d заполнена, отображаем содержимое", chatID)
	var anketaText strings.Builder
	anketaText.WriteString("📋 <b>Ваша анкета:</b>\n\n")

	if profile.SkinType != "" {
		anketaText.WriteString(fmt.Sprintf("👤 <b>Тип кожи:</b> %s\n", profile.SkinType))
	}
	if profile.Age != "" {
		anketaText.WriteString(fmt.Sprintf("📅 <b>Возраст:</b> %s\n", profile.Age))
	}
	if profile.Gender != "" {
		anketaText.WriteString(fmt.Sprintf("🚻 <b>Пол:</b> %s\n", profile.Gender))
	}
	if profile.Pregnancy != "" {
		anketaText.WriteString(fmt.Sprintf("🤱 <b>Беременность/лактация:</b> %s\n", profile.Pregnancy))
	}
	if profile.Concern != "" {
		anketaText.WriteString(fmt.Sprintf("💭 <b>Проблемы:</b> %s\n", profile.Concern))
	}
	if profile.Goal != "" {
		anketaText.WriteString(fmt.Sprintf("🎯 <b>Цель:</b> %s\n", profile.Goal))
	}
	if profile.Climate != "" {
		anketaText.WriteString(fmt.Sprintf("🌍 <b>Климат:</b> %s\n", profile.Climate))
	}
	if profile.Fitzpatrick != "" {
		anketaText.WriteString(fmt.Sprintf("☀️ <b>Тип кожи по Фитцпатрику:</b> %s\n", profile.Fitzpatrick))
	}
	if profile.Lifestyle != "" {
		anketaText.WriteString(fmt.Sprintf("🏃 <b>Образ жизни:</b> %s\n", profile.Lifestyle))
	}
	if profile.Diet != "" {
		anketaText.WriteString(fmt.Sprintf("🥗 <b>Питание:</b> %s\n", profile.Diet))
	}
	if profile.Allergy != "" {
		anketaText.WriteString(fmt.Sprintf("⚠️ <b>Аллергии:</b> %s\n", profile.Allergy))
	}

	// Отправляем фото с подписью вместо текстового сообщения
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("images/12.png"))
	photo.Caption = anketaText.String()
	photo.ParseMode = "HTML"

	// Создаем клавиатуру с действиями
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить анкету", "delete_anketa"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Пройти анкету заново", "retake_anketa"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_to_start"),
		),
	)
	photo.ReplyMarkup = keyboard
	bot.Send(photo)
}

// handleDeleteAnketa обрабатывает удаление анкеты
func handleDeleteAnketa(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	// Очищаем локальное состояние
	delete(userStates, chatID)

	// Очищаем профиль пользователя через API
	err := database.EmptyUserProfile(chatID)
	if err != nil {
		errorMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ошибка удаления анкеты: %v", err))
		bot.Send(errorMsg)
		return
	}

	// Отправляем сообщение об успешном удалении с кнопкой "Назад"
	msg := tgbotapi.NewMessage(chatID, "✅ Ваша анкета удалена!")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_to_start"),
		),
	)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

// handleBackToStart возвращает к главному меню
func handleBackToStart(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	// Отправляем фото с приветственным сообщением (как в /start)
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("images/01.png"))
	photo.Caption = `✨ Я — твой умный бьюти-бот, созданный, чтобы наконец навести порядок в косметичке. Этот бот - часть проекта Cos AI, созданного для того, чтобы помочь тебе собрать персонализированный уход за кожей.
Хочешь попробовать? Давай начнем с небольшой анкеты 💬👇`
	photo.ParseMode = "HTML"

	// Создаем клавиатуру с кнопками
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Анкета", "start_form"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🤖 Рекомендации", "recommendations"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧴 Мои продукты", "my_products"),
		),
	)
	photo.ReplyMarkup = keyboard
	bot.Send(photo)
}

// getUserState получает состояние пользователя с fallback на локальное хранение
func getUserState(chatID int64) *models.UserState {
	// Сначала проверяем локальное состояние (для активной анкеты)
	if localState, exists := userStates[chatID]; exists {
		log.Printf("Используем локальное состояние для пользователя %d: шаг %d", chatID, localState.Step)
		return localState
	}

	// Если локального состояния нет, пытаемся получить из API
	state, err := database.GetUserState(chatID)
	if err != nil {
		log.Printf("Ошибка получения состояния из API, создаем новое: %v", err)
		// Если API недоступен, создаем новое состояние
		return &models.UserState{Step: 0}
	}
	log.Printf("Получено состояние из API для пользователя %d: шаг %d", chatID, state.Step)
	return state
}

// saveUserState сохраняет состояние пользователя с fallback на локальное хранение
func saveUserState(chatID int64, state *models.UserState) {
	// Всегда сохраняем локально для активной анкеты
	userStates[chatID] = state
	log.Printf("Состояние пользователя %d сохранено локально: шаг %d", chatID, state.Step)

	// Также пытаемся сохранить в API (для персистентности)
	if err := database.SaveUserState(chatID, state); err != nil {
		log.Printf("Ошибка сохранения в API: %v", err)
	} else {
		log.Printf("Состояние пользователя %d также сохранено в API: шаг %d", chatID, state.Step)
	}
}

// handleInlineQuery обрабатывает inline запросы
func handleInlineQuery(bot *tgbotapi.BotAPI, inlineQuery *tgbotapi.InlineQuery) {
	query := inlineQuery.Query
	userID := inlineQuery.From.ID

	log.Printf("[INLINE] Получен inline запрос от пользователя %d: '%s'", userID, query)
	log.Printf("[INLINE] Детали запроса: ID=%s, Offset=%s", inlineQuery.ID, inlineQuery.Offset)

	// Оставляем запрос как есть, включая "add"
	log.Printf("[INLINE] Используем полный запрос: '%s'", query)

	// Проверяем, что запрос содержит минимум 3 символа
	if len(query) < 3 {
		log.Printf("[INLINE] Запрос слишком короткий (%d символов), показываем сообщение", len(query))
		// Создаем результат с сообщением о коротком запросе
		result := tgbotapi.NewInlineQueryResultArticle(
			"too_short",
			"⚠️ Запрос слишком короткий",
			"Введите минимум 3 символа для поиска продуктов.",
		)
		result.Description = "Минимум 3 символа для поиска"

		answerInlineQuery := tgbotapi.InlineConfig{
			InlineQueryID: inlineQuery.ID,
			Results:       []interface{}{result},
		}
		bot.Request(answerInlineQuery)
		return
	}

	// Выполняем поиск продуктов через API
	log.Printf("[INLINE] Выполняем поиск продуктов для запроса: '%s'", query)
	products, err := database.SearchProducts(query, 20, 0, nil, nil, nil, nil)
	if err != nil {
		log.Printf("[INLINE] Ошибка поиска продуктов для inline запроса: %v", err)
		// Создаем результат с сообщением об ошибке
		result := tgbotapi.NewInlineQueryResultArticle(
			"error",
			"❌ Ошибка поиска",
			"Произошла ошибка при поиске продуктов. Попробуйте позже.",
		)
		result.Description = "Ошибка соединения с сервером"

		answerInlineQuery := tgbotapi.InlineConfig{
			InlineQueryID: inlineQuery.ID,
			Results:       []interface{}{result},
		}
		bot.Request(answerInlineQuery)
		return
	}

	log.Printf("[INLINE] Найдено %d продуктов для запроса '%s'", len(products), query)

	if len(products) == 0 {
		// Создаем результат "Продукт не найден"
		log.Printf("[INLINE] Продукты не найдены, показываем сообщение 'Продукт не найден'")

		// Создаем inline результат с сообщением
		result := tgbotapi.NewInlineQueryResultArticle(
			"not_found",
			"❌ Продукт не найден",
			"По вашему запросу ничего не найдено. Попробуйте другой поисковый запрос.",
		)
		result.Description = fmt.Sprintf("По запросу '%s' ничего не найдено", query)

		answerInlineQuery := tgbotapi.InlineConfig{
			InlineQueryID: inlineQuery.ID,
			Results:       []interface{}{result},
		}
		bot.Request(answerInlineQuery)
		return
	}

	// Создаем inline результаты
	log.Printf("[INLINE] Создаем inline результаты для %d продуктов", len(products))
	var results []interface{}
	for i, product := range products {
		log.Printf("[INLINE] Обрабатываем продукт %d: %s %s", i+1, product.Brand, product.Title)

		// Создаем описание продукта (только детали, без названия)
		var description string
		if product.Details != "" {
			// Обрезаем описание если оно слишком длинное (Telegram ограничивает до 512 символов)
			details := product.Details
			if len(details) > 200 {
				details = details[:197] + "..."
			}
			description = details
		}

		// Создаем результат как статью с картинкой
		result := tgbotapi.NewInlineQueryResultArticle(
			fmt.Sprintf("product_%d", product.ID),
			fmt.Sprintf("%s %s", product.Brand, product.Title),
			description,
		)
		result.Description = description
		result.ThumbURL = product.Image

		// Добавляем кнопку для добавления в коллекцию
		result.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
				{
					tgbotapi.NewInlineKeyboardButtonData("➕ Добавить в коллекцию", fmt.Sprintf("add_product_%d", product.ID)),
				},
			},
		}

		results = append(results, result)

		// Если у нас много продуктов, ограничиваем количество для тестирования
		if i >= 9 { // Показываем максимум 10 продуктов
			break
		}
	}

	// Отправляем ответ на inline запрос
	log.Printf("[INLINE] Отправляем ответ с %d результатами", len(results))
	answerInlineQuery := tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.ID,
		Results:       results,
		CacheTime:     300, // Кешируем результаты на 5 минут
	}

	response, err := bot.Request(answerInlineQuery)
	if err != nil {
		log.Printf("[INLINE] Ошибка отправки inline ответа: %v", err)
	} else {
		log.Printf("[INLINE] Inline ответ успешно отправлен: %+v", response)
	}
}
