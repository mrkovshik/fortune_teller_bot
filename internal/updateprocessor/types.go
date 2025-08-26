package updateprocessor

import (
	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/steps"
)

type UpdateProcessor interface {
	ProcessMessage(message *model.Message) (map[string]interface{}, error)
	ProcessCallback(callback *model.CallbackQuery) (map[string]interface{}, error)
}

type BookStorage interface {
	GetRandomSentenceFromBook(bookName string, seed int64) (string, error)
	ListBooks() ([]string, error)
	GetRandomBookTitle() string
}

type StepStorage interface {
	Peek(chatID int64) (steps.ChatStep, error)
	PeekPrevious(chatID int64) (steps.ChatStep, error)
	Push(chatID int64, step steps.ChatStep) error
	Clear(chatID int64)
}

type UserDataStorage interface {
	GetUserData(chatID int64, key string) (string, error)
	AddUserData(chatID int64, key string, value string) error
	ClearUserData(chatID int64) error
}
