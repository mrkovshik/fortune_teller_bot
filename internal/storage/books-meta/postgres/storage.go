package postgres

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	booksmeta "github.com/mrkovshik/fortune_teller_bot/internal/storage/books-meta"
	"github.com/mrkovshik/fortune_teller_bot/internal/updateprocessor"
)

type BooksMetaStorage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) updateprocessor.BookStorage {
	return &BooksMetaStorage{db: db}
}

func (s BooksMetaStorage) GetBookByID(ctx context.Context, id int64) (book *booksmeta.Book, err error) {
	err = s.db.GetContext(ctx, &book, "SELECT * FROM books WHERE id=$1", id)
	return
}

func (s BooksMetaStorage) GetRandomBook(ctx context.Context, options ...booksmeta.ListOption) (book *booksmeta.Book, err error) {
	rand.NewSource(time.Now().UnixNano())
	books, err := s.ListBooks(ctx, options...)
	if err != nil {
		return nil, err
	}
	randomBookID := books[rand.Intn(len(books))].ID
	return s.GetBookByID(ctx, randomBookID)
}

func (s BooksMetaStorage) ListBooks(ctx context.Context, options ...booksmeta.ListOption) (books []*booksmeta.Book, err error) {
	query, args := prepareListQuery(options...)
	err = s.db.SelectContext(ctx, &books, query, args...)
	return
}

func prepareListQuery(options ...booksmeta.ListOption) (string, []any) {
	var opts booksmeta.ListOptions
	for _, opt := range options {
		opt(&opts)
	}

	sb := strings.Builder{}
	sb.WriteString(`SELECT id, title, author, language, format FROM books`)

	var (
		conds []string
		args  []any
		n     = 1
	)

	if opts.ID != nil {
		conds = append(conds, fmt.Sprintf("id = $%d", n))
		args = append(args, *opts.ID)
		n++
	}
	if opts.Title != nil && *opts.Title != "" {
		conds = append(conds, fmt.Sprintf("title ILIKE $%d", n))
		args = append(args, "%"+*opts.Title+"%")
		n++
	}
	if opts.Author != nil && *opts.Author != "" {
		conds = append(conds, fmt.Sprintf("author ILIKE $%d", n))
		args = append(args, "%"+*opts.Author+"%")
		n++
	}
	if opts.Language != nil {
		conds = append(conds, fmt.Sprintf(`language = $%d`, n))
		args = append(args, *opts.Language)
	}

	if len(conds) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(conds, " AND "))
	}

	if opts.Ordered {
		sb.WriteString(" ORDER BY author, title")
	}

	return sb.String(), args
}
