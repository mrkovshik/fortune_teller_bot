package postgres

import (
	"context"

	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	booksmeta "github.com/mrkovshik/fortune_teller_bot/internal/storage/books-meta"
	"github.com/mrkovshik/fortune_teller_bot/internal/textparser"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testID = 1

var (
	ctx       = context.Background()
	testBooks = []*booksmeta.Book{
		{
			Title:    "Test book 1",
			Author:   "Test book author 1",
			Language: config.English,
			Format:   textparser.Epub,
		},
		{
			Title:    "Test book 2",
			Author:   "Test book author 2",
			Language: config.English,
			Format:   textparser.Fb2,
		},
		{
			Title:    "Test book 3",
			Author:   "Test book author 3",
			Language: config.Russian,
			Format:   textparser.Epub,
		},
	}
)

var _ = Describe("Posgres DB storage", func() {
	It("Prepares list query", func() {
		query, args := prepareListQuery(booksmeta.WithLanguage(config.Russian), booksmeta.WithAuthor("test author"))
		Expect(query).To(Equal("SELECT id, title, author, lang, format FROM books WHERE author ILIKE $1 AND lang = $2"))
		Expect(args).To(ContainElements("%test author%", config.Russian))
	})

	It("Init DB", func() {
		bookStorage, err := NewStorage(&config.Config{DatabaseURI: DSN})
		Expect(err).NotTo(HaveOccurred())
		Expect(bookStorage).NotTo(BeNil())

		for _, book := range testBooks {
			err = bookStorage.AddBook(ctx, book)
			Expect(err).NotTo(HaveOccurred())
		}

		books, err := bookStorage.ListBooks(ctx, booksmeta.WithLanguage(config.Russian))
		Expect(err).NotTo(HaveOccurred())
		Expect(books).To(HaveLen(1))

		books, err = bookStorage.ListBooks(ctx, booksmeta.WithLanguage(config.English))
		Expect(err).NotTo(HaveOccurred())
		Expect(books).To(HaveLen(2))

		book, err := bookStorage.GetBookByID(ctx, books[0].ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(book.Title).To(Equal(books[0].Title))

		randBook, err := bookStorage.GetRandomBook(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(randBook).NotTo(BeNil())
		Expect(book.Title).NotTo(Equal(""))
	})

})
