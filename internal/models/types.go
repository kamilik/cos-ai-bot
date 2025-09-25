package models

// ========== Структуры для чеклиста ==========
type ChecklistButton struct {
	Text      string `json:"text"`
	Checklist struct {
		ID string `json:"id"`
	} `json:"checklist"`
}

type SubmitButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type ChecklistOptions struct {
	MinSelected  int          `json:"min_selected"`
	MaxSelected  int          `json:"max_selected"`
	SubmitButton SubmitButton `json:"submit_button"`
}

type ReplyMarkup struct {
	InlineKeyboard   [][]ChecklistButton `json:"inline_keyboard"`
	ChecklistOptions ChecklistOptions    `json:"checklist_options"`
}

type MessagePayload struct {
	ChatID      int64       `json:"chat_id"`
	Text        string      `json:"text"`
	ReplyMarkup ReplyMarkup `json:"reply_markup"`
	ParseMode   string      `json:"parse_mode,omitempty"`
}

// ========== Структуры для поиска продуктов ==========
type SearchResult struct {
	ID    string
	Brand string
	Title string
	URL   string
}

// ========== Структуры для sendChecklist ==========
type ChecklistItem struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type ChecklistSubmitButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type Checklist struct {
	Title        string                `json:"title"`
	Tasks        []ChecklistItem       `json:"tasks"`
	MinSelected  int                   `json:"min_selected"`
	MaxSelected  int                   `json:"max_selected"`
	SubmitButton ChecklistSubmitButton `json:"submit_button"`
}

type ChecklistPayload struct {
	ChatID    int64     `json:"chat_id"`
	Checklist Checklist `json:"checklist"`
	ParseMode string    `json:"parse_mode,omitempty"`
}

// ========== Структуры для продуктов ==========
type Product struct {
	Brand        string
	ProductTitle string
	Image        string
	Ingredients  string
	Description  string
}

// ========== Структуры для пользователей ==========
type UserState struct {
	Step        int
	SkinType    string
	Age         string
	Gender      string
	Pregnancy   string
	Concerns    string
	Goal        string
	Climate     string
	Fitzpatrick string
	Lifestyle   string
	Diet        string
	Allergies   string
}

// ========== Структуры для API ==========

// APIProduct представляет продукт из API
type APIProduct struct {
	ID      int    `json:"id"`
	Brand   string `json:"brand"`
	Title   string `json:"title"`
	Details string `json:"details"`
	Image   string `json:"image"`
}

// APIProductDetail представляет детальную информацию о продукте из API
type APIProductDetail struct {
	ID          int                `json:"id"`
	Brand       string             `json:"brand"`
	Title       string             `json:"title"`
	Details     string             `json:"details"`
	Image       string             `json:"image"`
	Ingredients []APIIngredientRef `json:"ingredients"`
}

// APIIngredientRef представляет ссылку на ингредиент в продукте
type APIIngredientRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// APIIngredient представляет ингредиент из API
type APIIngredient struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	AltName     string   `json:"alt_name"`
	Description string   `json:"description"`
	Slug        string   `json:"slug"`
	Functions   []string `json:"functions"`
}

// APIUserProduct представляет продукт пользователя из API
type APIUserProduct struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Brand     string `json:"brand"`
	Title     string `json:"title"`
	Details   string `json:"details"`
	Image     string `json:"image"`
	AddedAt   string `json:"added_at"`
}

// APIUserProfile представляет профиль пользователя из API
type APIUserProfile struct {
	UserID      int64  `json:"user_id"`
	SkinType    string `json:"skin_type"`
	Age         string `json:"age"`
	Gender      string `json:"gender"`
	Pregnancy   string `json:"pregnancy"`
	Concern     string `json:"concern"`
	Goal        string `json:"goal"`
	Climate     string `json:"climate"`
	Fitzpatrick string `json:"fitzpatrick"`
	Lifestyle   string `json:"lifestyle"`
	Diet        string `json:"diet"`
	Allergy     string `json:"allergy"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// APIUserProfileUpdate представляет данные для обновления профиля пользователя
type APIUserProfileUpdate struct {
	SkinType    string `json:"skin_type"`
	Age         string `json:"age"`
	Gender      string `json:"gender"`
	Pregnancy   string `json:"pregnancy"`
	Concern     string `json:"concern"`
	Goal        string `json:"goal"`
	Climate     string `json:"climate"`
	Fitzpatrick string `json:"fitzpatrick"`
	Lifestyle   string `json:"lifestyle"`
	Diet        string `json:"diet"`
	Allergy     string `json:"allergy"`
}

// APIProductCreate представляет данные для создания нового продукта
type APIProductCreate struct {
	Brand       string             `json:"brand"`
	Details     string             `json:"details"`
	Image       string             `json:"image"`
	Ingredients []APIIngredientRef `json:"ingredients"`
	Title       string             `json:"title"`
}
