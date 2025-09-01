package updateprocessor

import (
	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	booksmeta "github.com/mrkovshik/fortune_teller_bot/internal/storage/books-meta"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/steps"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/userdata"
)

type Quote struct {
	Text string
	Book *booksmeta.Book
}

type UpdateProcessor interface {
	ProcessMessage(message *model.Message) (map[string]interface{}, error)
	ProcessCallback(callback *model.CallbackQuery) (map[string]interface{}, error)
}

type BookStorage interface {
	GetBookByID(id int64) (*booksmeta.Book, error)
	ListBooks(options ...booksmeta.ListOption) ([]*booksmeta.Book, error)
	GetRandomBook(options ...booksmeta.ListOption) (*booksmeta.Book, error)
}

type BookRepository interface {
	LoadBook(book *booksmeta.Book) ([]byte, error)
}

type StepStorage interface {
	Peek(chatID int64) (steps.ChatStep, error)
	PeekPrevious(chatID int64) (steps.ChatStep, error)
	Push(chatID int64, step steps.ChatStep) error
	Clear(chatID int64)
}

type UserDataStorage interface {
	GetUserData(chatID int64) (userdata.UserData, error)
	SaveUserData(chatID int64, key string, value any) error
	ClearUserData(chatID int64) error
}
