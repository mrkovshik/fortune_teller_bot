package basic

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

var (
	goBackButton = InlineKeyboardButton{
		Text:         "Вернуться назад",
		CallbackData: "go_back",
	}
)

var (
	listBooksButton = InlineKeyboardButton{
		Text:         "Выбрать книгу для гадания",
		CallbackData: "list_books",
	}
	getQuoteFromRandomBookButton = InlineKeyboardButton{
		Text:         "Случайная фраза из случайной книги",
		CallbackData: "random_sentence",
	}

	askQuestionButton = InlineKeyboardButton{
		Text:         "Гадать на конкретный вопрос",
		CallbackData: "ask_question",
	}

	helpButton = InlineKeyboardButton{
		Text:         "Что это за бот?",
		CallbackData: "help",
	}

	startMenu = InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{listBooksButton},
			{getQuoteFromRandomBookButton},
			{askQuestionButton},
			{helpButton},
		},
	}
)

var (
	listBooksForAnswerButton = InlineKeyboardButton{
		Text:         "Выбрать книгу для получения ответа на вопрос",
		CallbackData: "list_books_for_answer",
	}
	getAnswerFromRandomBookButton = InlineKeyboardButton{
		Text:         "Получить ответ на вопрос из случайной книги",
		CallbackData: "list_books",
	}

	askQuestionMenu = InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{listBooksForAnswerButton},
			{getAnswerFromRandomBookButton},
			{goBackButton},
		},
	}
)
