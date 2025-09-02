package updateprocessor

import (
	"context"

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
	ProcessMessage(ctx context.Context, message *model.Message) (map[string]interface{}, error)
	ProcessCallback(ctx context.Context, callback *model.CallbackQuery) (map[string]interface{}, error)
}

type BookStorage interface {
	GetBookByID(ctx context.Context, id int64) (*booksmeta.Book, error)
	ListBooks(ctx context.Context, options ...booksmeta.ListOption) ([]*booksmeta.Book, error)
	GetRandomBook(ctx context.Context, options ...booksmeta.ListOption) (*booksmeta.Book, error)
	AddBook(ctx context.Context, book *booksmeta.Book) error
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
