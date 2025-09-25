package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ShowSkincareFormStep показывает шаг формы по уходу за кожей
func ShowSkincareFormStep(bot *tgbotapi.BotAPI, chatID int64, step int) {
	log.Printf("ShowSkincareFormStep: показываем шаг %d для пользователя %d", step, chatID)
	photoUrl := "https://images.unsplash.com/photo-1464983953574-0892a716854b"
	var caption string
	var keyboard tgbotapi.InlineKeyboardMarkup

	switch step {
	case 1:
		caption = "Какой ваш тип кожи?\n\nТип кожи влияет на выбор текстур и активных ингредиентов — от этого зависит, как хорошо средство будет работать"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Сухая", "skin_dry"),
				tgbotapi.NewInlineKeyboardButtonData("Жирная", "skin_oily"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Нормальная", "skin_normal"),
				tgbotapi.NewInlineKeyboardButtonData("Чувствительная", "skin_sensitive"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Комбинированная", "skin_combined"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Я не знаю какой у меня тип", "skin_unknown"),
			),
		)
	case 2:
		caption = "Какой ваш возраст?\n\nВ 20, 30 и 50 лет коже нужны разные вещи. Уточним возраст, чтобы подобрать то, что подходит именно вам"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("<18", "age_18_minus"),
				tgbotapi.NewInlineKeyboardButtonData("18–24", "age_18_24"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("25–34", "age_25_34"),
				tgbotapi.NewInlineKeyboardButtonData("35–44", "age_35_44"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("45+", "age_45_plus"),
				tgbotapi.NewInlineKeyboardButtonData("Не учитывать", "age_ignore"),
			),
		)
	case 3:
		caption = "Укажите ваш пол\n\nМужская и женская кожа отличаются по структуре и гормональному фону — это помогает нам точнее подобрать уход"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Мужчина", "gender_male"),
				tgbotapi.NewInlineKeyboardButtonData("Женщина", "gender_female"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Другое", "gender_other"),
				tgbotapi.NewInlineKeyboardButtonData("Не учитывать", "gender_ignore"),
			),
		)
	case 4:
		caption = "Находитесь ли вы сейчас в периоде беременности или кормления?\n\nНекоторые ингредиенты не рекомендуются в этот период. Мы подберём безопасные альтернативы."
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Беременность", "pregnancy"),
				tgbotapi.NewInlineKeyboardButtonData("Лактация", "lactation"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("И то, и другое", "pregnancy_and_lactation"),
				tgbotapi.NewInlineKeyboardButtonData("Ничего из перечисленного", "none_of_above"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Не учитывать", "pregnancy_ignore"),
			),
		)
	case 5:
		caption = "Что беспокоит вас больше всего?\n\nЧто вы ждёте от ухода: убрать проблему, предотвратить, освежить внешний вид? Ответ в свободной форме, например:\n\nХочу исправить повышенную чувствительность у моей кожи, а так же меня беспокоит акне и чёрные точки"
		// Свободный ввод, без кнопок
	case 6:
		caption = "Какой результат вы хотите получить?\n\nВаша цель = наша стратегия. Разберёмся, куда стремиться. Если ни один из вариантов не подходит, вы можете написать ответ в свободной форме"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Увлажнение и питание", "goal_hydration"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Выравнивание тона", "goal_tone"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Антивозрастной уход", "goal_antiage"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Улучшение текстуры", "goal_texture"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Освежить и поддерживать", "goal_refresh"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Минимализм, только базовый уход", "goal_minimalism"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Другое (напишу сам)", "goal_other"),
			),
		)
	case 7:
		caption = "Какой у вас климат?\n\nКлимат влияет на потребности кожи в увлажнении и защите"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Сухой", "climate_dry"),
				tgbotapi.NewInlineKeyboardButtonData("Влажный", "climate_humid"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Жаркий", "climate_hot"),
				tgbotapi.NewInlineKeyboardButtonData("Холодный", "climate_cold"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Переменный / умеренный", "climate_temperate"),
				tgbotapi.NewInlineKeyboardButtonData("Загрязнённый (город, смог, пыль)", "climate_polluted"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Живу в нескольких климатах (путешествую/переезды)", "climate_multiple"),
				tgbotapi.NewInlineKeyboardButtonData("Не знаю", "climate_unknown"),
			),
		)
	case 8:
		caption = "Как бы вы описали свою кожу по реакции на солнце?\n\nЭто поможет подобрать правильную защиту от солнца"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("I – очень светлая, всегда обгорает", "fitzpatrick_1"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("II – светлая, обгорает, но может немного загорать", "fitzpatrick_2"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("III – светло-смуглая, легко загорает", "fitzpatrick_3"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("IV – смуглая, редко обгорает", "fitzpatrick_4"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("V – тёмная, почти не обгорает", "fitzpatrick_5"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("VI – очень тёмная, никогда не обгорает", "fitzpatrick_6"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Не знаю / Не хочу указывать", "fitzpatrick_unknown"),
			),
		)
	case 9:
		caption = "Какой у вас ритм жизни?\n\nОбраз жизни влияет на выбор средств и режим ухода"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Частые стрессы", "lifestyle_stress"),
				tgbotapi.NewInlineKeyboardButtonData("Недосып / сбитый режим", "lifestyle_sleep"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Много экранного времени", "lifestyle_screen"),
				tgbotapi.NewInlineKeyboardButtonData("Часто потею (спорт, жара и т.д.)", "lifestyle_sweat"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Работаю за компьютером", "lifestyle_computer"),
				tgbotapi.NewInlineKeyboardButtonData("Активно двигаюсь в течение дня", "lifestyle_active"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Регулярно на улице", "lifestyle_outdoor"),
				tgbotapi.NewInlineKeyboardButtonData("Пассивный / домашний образ жизни", "lifestyle_passive"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Другое", "lifestyle_other"),
			),
		)
	case 10:
		caption = "Есть ли у вас особенности в питании или убеждения, которые важно учесть?\n\nЭто поможет подобрать подходящие ингредиенты"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Веганство", "diet_vegan"),
				tgbotapi.NewInlineKeyboardButtonData("Вегетарианство", "diet_vegetarian"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Халяль", "diet_halal"),
				tgbotapi.NewInlineKeyboardButtonData("Кето / Палео / Низкоуглеводная", "diet_keto"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Безглютеновая", "diet_gluten_free"),
				tgbotapi.NewInlineKeyboardButtonData("Я избегаю спирта в составе", "diet_no_alcohol"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Я избегаю компонентов животного происхождения", "diet_no_animal"),
				tgbotapi.NewInlineKeyboardButtonData("Нет особых ограничений", "diet_none"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Другое", "diet_other"),
			),
		)
	case 11:
		caption = "Есть ли у вас аллергии или непереносимость?\n\nВажно знать, чтобы исключить проблемные ингредиенты"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Нет аллергий", "allergies_none"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Аллергия на никель", "allergies_nickel"),
				tgbotapi.NewInlineKeyboardButtonData("Аллергия на ланолин", "allergies_lanolin"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Аллергия на отдушки", "allergies_fragrance"),
				tgbotapi.NewInlineKeyboardButtonData("Аллергия на консерванты", "allergies_preservatives"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Другое (напишу сам)", "allergies_other"),
			),
		)
	default:
		log.Printf("Неизвестный шаг формы: %d", step)
		return
	}

	// Отправляем фото с подписью и клавиатурой
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(photoUrl))
	photo.Caption = caption
	photo.ParseMode = "HTML"

	if step == 5 {
		// Для свободного ввода не добавляем клавиатуру, но и не убираем существующую
		// Пользователь может ввести текст напрямую
	} else {
		photo.ReplyMarkup = keyboard
	}

	if _, err := bot.Send(photo); err != nil {
		log.Printf("Ошибка отправки фото: %v", err)
	} else {
		log.Printf("ShowSkincareFormStep: успешно отправлен шаг %d для пользователя %d", step, chatID)
	}
}
