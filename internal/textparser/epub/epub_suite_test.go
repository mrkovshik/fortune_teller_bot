package epub_test

import (
	"context"
	"testing"

	"github.com/mrkovshik/fortune_teller_bot/internal/books-repository/embedded"
	booksmeta "github.com/mrkovshik/fortune_teller_bot/internal/storage/books-meta"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/books-meta/inmemory"
	"github.com/mrkovshik/fortune_teller_bot/internal/textparser"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var TestBooks map[int64][]byte

func TestEmbedded(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Embedded suite")
}

var _ = BeforeSuite(func() {
	ctx := context.Background()
	booksRepo := embedded.NewRepository()
	booksMeta := inmemory.Storage

	books, err := booksMeta.ListBooks(ctx, booksmeta.WithFormat(textparser.Epub))
	Expect(err).NotTo(HaveOccurred())
	TestBooks = make(map[int64][]byte)
	for _, book := range books {
		TestBooks[book.ID], err = booksRepo.LoadBook(book)
		Expect(err).NotTo(HaveOccurred())
		Expect(TestBooks[book.ID]).NotTo(BeNil())
	}
})
