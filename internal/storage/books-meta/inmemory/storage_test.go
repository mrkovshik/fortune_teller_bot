package inmemory_test

import (
	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	booksmeta "github.com/mrkovshik/fortune_teller_bot/internal/storage/books-meta"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/books-meta/inmemory"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testID = 1

var _ = Describe("Local storage functions test", func() {
	var (
		err         error
		testStorage = inmemory.Storage
	)
	BeforeEach(func() {
		Expect(err).NotTo(HaveOccurred())
	})

	It("Builds books list", func() {
		res := 0
		for lang := range config.SupportedLanguages {
			booksList, err := testStorage.ListBooks(booksmeta.WithLanguage(lang))
			Expect(err).NotTo(HaveOccurred())
			Expect(booksList).NotTo(BeNil())
			Expect(len(booksList)).To(BeNumerically(">", 0))
			res += len(booksList)
		}
		Expect(res).To(Equal(len(inmemory.Storage)))
	})

	It("Takes random book from storage", func() {
		book, err := testStorage.GetRandomBook()
		Expect(err).NotTo(HaveOccurred())
		Expect(book).NotTo(BeNil())
	})

	It("Takes specific book from storage", func() {
		book, err := testStorage.GetBookByID(testID)
		Expect(err).NotTo(HaveOccurred())
		Expect(book).NotTo(BeNil())
		Expect(book.Title).To(Equal("Трое в лодке, не считая собаки"))
	})
})
