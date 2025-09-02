package booksmeta

import (
	"fmt"

	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/fortune_teller_bot/internal/textparser"
	"github.com/mrkovshik/fortune_teller_bot/internal/textparser/epub"
	"github.com/mrkovshik/fortune_teller_bot/internal/textparser/fb2"
)

type TextParser interface {
	ParseRandomSentence(data []byte, seed int64) (string, error)
}

type Book struct {
	ID       int64             `json:"id" db:"id"`
	Title    string            `json:"title" db:"title"`
	Author   string            `json:"author" db:"author"`
	Language config.Language   `json:"language" db:"lang"`
	Format   textparser.Format `json:"format" db:"format"`
}

func (b Book) GetParser() (TextParser, error) {
	switch b.Format {
	case textparser.Epub:
		return epub.NewTextParser(), nil
	case textparser.Fb2:
		return fb2.NewTextParser(), nil
	default:
		return nil, fmt.Errorf("unknown format: %s", b.Format)
	}
}

type ListOptions struct {
	ID       *int64
	Title    *string
	Author   *string
	Language *config.Language
	Ordered  bool
}

type ListOption func(o *ListOptions)

func WithID(id int64) ListOption {
	return func(o *ListOptions) {
		o.ID = &id
	}
}
func WithTitle(title string) ListOption {
	return func(o *ListOptions) {
		o.Title = &title
	}
}
func WithAuthor(author string) ListOption {
	return func(o *ListOptions) {
		o.Author = &author
	}
}
func WithLanguage(lang config.Language) ListOption {
	return func(o *ListOptions) {
		o.Language = &lang
	}
}

func Ordered() ListOption {
	return func(o *ListOptions) {
		o.Ordered = true
	}
}
