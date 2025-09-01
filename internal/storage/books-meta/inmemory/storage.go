package inmemory

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	booksmeta "github.com/mrkovshik/fortune_teller_bot/internal/storage/books-meta"
)

type BooksMetaStorage map[int64]*booksmeta.Book

var Storage = BooksMetaStorage{
	1: {
		ID:       1,
		Title:    "Трое в лодке, не считая собаки",
		Author:   "Дж.К.Джером",
		Language: config.Russian,
	},

	2: {
		ID:       2,
		Title:    "Портрет Дориана Грея",
		Author:   "Оскар Уайлд",
		Language: config.Russian,
	},

	3: {
		ID:       3,
		Title:    "Господа Головлёвы",
		Author:   "М.Е.Салтыков-Щедрин",
		Language: config.Russian,
	},
	4: {
		ID:       4,
		Title:    "Дети капитана Гранта",
		Author:   "Ж.Верн",
		Language: config.Russian,
	},
	5: {
		ID:       5,
		Title:    "Зов Ктулху",
		Author:   "Г.Лавкрафт",
		Language: config.Russian,
	},
	6: {
		ID:       6,
		Title:    "Избранное",
		Author:   "М.Зощенко",
		Language: config.Russian,
	},
	7: {
		ID:       7,
		Title:    "Moby Dick",
		Author:   "H.Melville",
		Language: config.English,
	},
	8: {
		ID:       8,
		Title:    "Frankenstein",
		Author:   "M.Shelley",
		Language: config.English,
	},
	9: {
		ID:       9,
		Title:    "Tremendous Trifles",
		Author:   "G.K.Chesterton",
		Language: config.English,
	},
	10: {
		ID:       10,
		Title:    "The Wisdom of Life",
		Author:   "A.Schopenhauer",
		Language: config.English,
	},
}

func (s BooksMetaStorage) GetBookByID(id int64) (*booksmeta.Book, error) {
	book, ok := s[id]
	if !ok {
		return nil, fmt.Errorf("book id %d not found", id)
	}
	return book, nil
}

func (s BooksMetaStorage) GetRandomBook(options ...booksmeta.ListOption) (*booksmeta.Book, error) {
	booksList, err := s.ListBooks(options...)
	if err != nil {
		return nil, err
	}
	rand.NewSource(time.Now().UnixNano())
	idx := rand.Intn(len(booksList))
	return booksList[idx], nil
}

func (s BooksMetaStorage) ListBooks(options ...booksmeta.ListOption) ([]*booksmeta.Book, error) {
	var opts booksmeta.ListOptions
	for _, opt := range options {
		opt(&opts)
	}
	result := make([]*booksmeta.Book, 0)
	for _, book := range s {
		match := true
		if opts.ID != nil && *opts.ID != book.ID {
			match = false
		}
		if opts.Author != nil && *opts.Author != book.Author {
			match = false
		}
		if opts.Language != nil && *opts.Language != book.Language {
			match = false
		}
		if opts.Title != nil && *opts.Title != book.Title {
			match = false
		}
		if match {
			result = append(result, book)
		}
	}
	if opts.Ordered {
		sort.Slice(result, func(i, j int) bool {
			return result[i].ID < result[j].ID
		})
	}
	return result, nil
}
