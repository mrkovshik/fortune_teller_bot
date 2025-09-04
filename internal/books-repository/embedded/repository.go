package embedded

import (
	"embed"
	"fmt"

	booksmeta "github.com/mrkovshik/fortune_teller_bot/internal/storage/books-meta"
)

type Repository struct {
	embed.FS
}

//go:embed data/*
var booksFS embed.FS

func NewRepository() *Repository {
	return &Repository{
		FS: booksFS,
	}
}

func (r *Repository) LoadBook(book *booksmeta.Book) ([]byte, error) {
	filePath := fmt.Sprintf("data/%d.%s", book.ID, book.Format)
	data, err := r.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't read book: %w", err)
	}
	return data, nil
}
