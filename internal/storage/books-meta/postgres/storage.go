package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/fortune_teller_bot/db"
	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	booksmeta "github.com/mrkovshik/fortune_teller_bot/internal/storage/books-meta"
	"github.com/mrkovshik/fortune_teller_bot/internal/updateprocessor"
)

type BooksMetaStorage struct {
	db *sqlx.DB
}

func NewStorage(cfg *config.Config) (updateprocessor.BookStorage, error) {
	dataBase, err := sqlx.Connect("postgres", cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(dataBase.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	src, err := iofs.New(db.MigrationsFS, "migrations")
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return &BooksMetaStorage{db: dataBase}, nil
}

func (s BooksMetaStorage) GetBookByID(ctx context.Context, id int64) (book *booksmeta.Book, err error) {
	book = &booksmeta.Book{}
	err = s.db.GetContext(ctx, book, "SELECT * FROM books WHERE id=$1", id)
	return
}

func (s BooksMetaStorage) AddBook(ctx context.Context, book *booksmeta.Book) error {
	if _, err := s.db.ExecContext(ctx, "INSERT INTO books (title, author, format, lang) VALUES ($1, $2, $3, $4)", book.Title, book.Author, book.Format, book.Language); err != nil {
		return err
	}
	return nil
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
	sb.WriteString(`SELECT id, title, author, lang, format FROM books`)

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
		conds = append(conds, fmt.Sprintf(`lang = $%d`, n))
		args = append(args, *opts.Language)
	}

	if opts.Format != nil {
		conds = append(conds, fmt.Sprintf(`format = $%d`, n))
		args = append(args, *opts.Format)
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
