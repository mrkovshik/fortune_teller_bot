package basic

import (
	"fmt"
	"strings"

	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/book_storage/local"
	"github.com/mrkovshik/fortune_teller_bot/internal/update_processor"
	"go.uber.org/zap"
)

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

func (cp *UpdateProcessor) ProcessUpdate(update *model.Update) (map[string]interface{}, error) {
	chatID := update.Message.Chat.ID
	command := update.Message.Text
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
		var buttons [][]map[string]string
		for _, book := range books {
			button := map[string]string{
				"text":          book,
				"callback_data": string(model.SelectBook) + book,
			}
			buttons = append(buttons, []map[string]string{button})
		}
		payload["text"] = "Выберите книгу, по которой будем предсказывать будущее:"
		payload["reply_markup"] = map[string]interface{}{
			"inline_keyboard": buttons,
		}
		state.CurrentStep = model.SelectBook
		cp.stateStorage.Update(chatID, state)

	case update_processor.GetMagicCommandName:
		text, err := cp.bookStorage.GetRandomSentenceFromBook(local.GetRandomBookTitle())
		if err != nil {
			return nil, err
		}
		payload["text"] = text

	case update_processor.StartCommandName:
		cp.stateStorage.Clear(chatID)
		payload["text"] = fmt.Sprintf("Чтобы посмотреть перечень доступных книг, выберите команду %s, а чтобы узнать ответ на ваш вопрос, выберите %s", update_processor.ListBooksCommandName, update_processor.GetMagicCommandName)
	}

	switch state.CurrentStep {
	case model.SelectBook:
		if update.CallbackQuery == nil {
			return nil, fmt.Errorf("step %s must have a callback query", state.CurrentStep)
		}
		bookTitle := strings.TrimPrefix(update.CallbackQuery.Data, string(model.SelectBook))
		text, err := cp.bookStorage.GetRandomSentenceFromBook(bookTitle)
		if err != nil {
			return nil, err
		}
		payload["text"] = text
	default:
		payload["text"] = fmt.Sprintf("Чтобы посмотреть перечень доступных книг, выберите команду %s, а чтобы узнать ответ на ваш вопрос, выберите %s", update_processor.ListBooksCommandName, update_processor.GetMagicCommandName)
	}

	return payload, nil
}
