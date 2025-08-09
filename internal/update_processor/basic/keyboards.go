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
	ListBooksCommandName          CallbackCommand = "list_books"
	GetRandomSentenceCommandName  CallbackCommand = "get_random_sentence"
	GetSentenceCommandName        CallbackCommand = "random_book_sentence"
	AskQuestionCommandName        CallbackCommand = "ask_question"
	ListBooksForAnswerCommandName CallbackCommand = "list_books"
	HelpCommandName               CallbackCommand = "help"
)

var (
	getRandomQuoteButton = InlineKeyboardButton{
		Text:         "Получить случайную цитату",
		CallbackData: "list_books",
	}

	listBooksButton = InlineKeyboardButton{
		Text:         "Выбрать конкретную книгу",
		CallbackData: ListBooksCommandName,
	}
	getQuoteFromRandomBookButton = InlineKeyboardButton{
		Text:         "Использовать случайную книгу",
		CallbackData: GetRandomSentenceCommandName,
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
	askQuestionMenu = InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{listBooksButton},
			{getQuoteFromRandomBookButton},
			{goBackButton},
		},
	}
)
