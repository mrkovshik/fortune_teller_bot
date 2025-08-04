package basic

import (
	"fmt"
	"strings"

	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/book_storage/local"
	"github.com/mrkovshik/fortune_teller_bot/internal/update_processor"
	"go.uber.org/zap"
)

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type BookStorage interface {
	GetRandomSentenceFromBook(bookName string) (string, error)
	ListBooks() ([]string, error)
}

type StateStorage interface {
	Update(chatID int64, state model.ChatState)
	Get(chatID int64) (model.ChatState, error)
	Add(chatID int64, state model.ChatState)
	Clear(chatID int64)
}

type UpdateProcessor struct {
	logger       *zap.SugaredLogger
	bookStorage  BookStorage
	stateStorage StateStorage
}

func NewUpdateProcessor(bookStorage BookStorage, stateStorage StateStorage, logger *zap.SugaredLogger) *UpdateProcessor {
	return &UpdateProcessor{
		logger:       logger,
		bookStorage:  bookStorage,
		stateStorage: stateStorage,
	}
}

func (cp *UpdateProcessor) ProcessMessage(message *model.Message) (map[string]interface{}, error) {
	chatID := message.Chat.ID
	command := message.Text
	payload := map[string]interface{}{
		"chat_id": chatID,
	}

	state, err := cp.stateStorage.Get(chatID)
	if err != nil {
		return nil, err
	}
	switch command {
	case update_processor.ListBooksCommandName:
		books, err := cp.bookStorage.ListBooks()
		if err != nil {
			return nil, fmt.Errorf(`failed to list books: %w`, err)
		}
		var keyboard [][]InlineKeyboardButton
		for _, book := range books {
			button := InlineKeyboardButton{
				Text:         book,
				CallbackData: fmt.Sprintf("%s:%s", model.SelectBook, local.TitleToFileName[book]),
			}
			// одна кнопка в строке
			keyboard = append(keyboard, []InlineKeyboardButton{button})
		}
		payload["text"] = "Выберите книгу, по которой будем предсказывать будущее:"
		payload["reply_markup"] = InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		}
		cp.logger.Info(keyboard)
		state.CurrentStep = model.SelectBook
		cp.stateStorage.Update(chatID, state)

	case update_processor.GetMagicCommandName:
		cp.stateStorage.Clear(chatID)
		text, err := cp.bookStorage.GetRandomSentenceFromBook(local.GetRandomBookTitle())
		if err != nil {
			return nil, err
		}
		payload["text"] = text

	case update_processor.StartCommandName:
		cp.stateStorage.Clear(chatID)
		payload["text"] = fmt.Sprintf("Чтобы посмотреть перечень доступных книг, выберите команду %s, а чтобы узнать ответ на ваш вопрос, выберите %s", update_processor.ListBooksCommandName, update_processor.GetMagicCommandName)

	default:
		payload["text"] = fmt.Sprintf("Чтобы посмотреть перечень доступных книг, выберите команду %s, а чтобы узнать ответ на ваш вопрос, выберите %s", update_processor.ListBooksCommandName, update_processor.GetMagicCommandName)
	}

	return payload, nil
}

func (cp *UpdateProcessor) ProcessCallback(callback *model.CallbackQuery) (map[string]interface{}, error) {
	chatID := callback.From.ID
	payload := map[string]interface{}{
		"chat_id": chatID,
	}
	fileName := strings.TrimPrefix(callback.Data, string(model.SelectBook))
	fileName = strings.TrimPrefix(fileName, ":")
	bookTitle, exist := local.FileNameToTitle[fileName]
	if !exist {
		payload["text"] = fmt.Sprintf("Книга с таким именем файла не найдена: %s", fileName)
		return payload, nil
	}
	text, err := cp.bookStorage.GetRandomSentenceFromBook(bookTitle)
	if err != nil {
		return nil, err
	}
	payload["text"] = text

	return payload, nil
}
