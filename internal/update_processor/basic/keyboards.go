package basic

type InlineKeyboardButton struct {
	Text         string          `json:"text"`
	CallbackData CallbackCommand `json:"callback_data"`
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

type CallbackCommand string

const (
	ListBooksCommandName         CallbackCommand = "list_books"
	GetRandomSentenceCommandName CallbackCommand = "get_random_sentence"
	UseRandomBookCommandName     CallbackCommand = "use_random_book"
	AskQuestionCommandName       CallbackCommand = "ask_question"
	GoBackCommandName            CallbackCommand = "go_back"
	HelpCommandName              CallbackCommand = "help"
)

var (
	getRandomQuoteButton = InlineKeyboardButton{
		Text:         "Получить случайную цитату",
		CallbackData: GetRandomSentenceCommandName,
	}

	listBooksButton = InlineKeyboardButton{
		Text:         "Выбрать конкретную книгу",
		CallbackData: ListBooksCommandName,
	}
	getQuoteFromRandomBookButton = InlineKeyboardButton{
		Text:         "Использовать случайную книгу",
		CallbackData: UseRandomBookCommandName,
	}

	askQuestionButton = InlineKeyboardButton{
		Text:         "Гадать на конкретный вопрос",
		CallbackData: AskQuestionCommandName,
	}

	helpButton = InlineKeyboardButton{
		Text:         "Что это за бот?",
		CallbackData: HelpCommandName,
	}

	startMenu = InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{getRandomQuoteButton},
			{askQuestionButton},
			{helpButton},
		},
	}
)

var (
	selectSourceMenu = InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{listBooksButton},
			{getQuoteFromRandomBookButton},
			{goBackButton},
		},
	}
)
