package updateprocessor

import (
	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/statestorage"
)

type UpdateProcessor interface {
	ProcessMessage(message *model.Message) (map[string]interface{}, error)
	ProcessCallback(callback *model.CallbackQuery) (map[string]interface{}, error)
}

type BookStorage interface {
	GetRandomSentenceFromBook(bookName string, seed int64) (string, error)
	ListBooks() ([]string, error)
}

type StateStorage interface {
	Update(chatID int64, state *statestorage.ChatState)
	Get(chatID int64) (*statestorage.ChatState, error)
	Add(chatID int64, state *statestorage.ChatState)
	Clear(chatID int64)
}
