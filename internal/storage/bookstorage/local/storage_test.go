package local_test

import (
	"time"

	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/bookstorage/local"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("Local storage functions test", func() {
	var (
		logger      *zap.Logger
		err         error
		testStorage *local.Storage
	)
	BeforeEach(func() {
		logger, err = zap.NewDevelopment()
		Expect(err).NotTo(HaveOccurred())
		testStorage = local.NewStorage(logger.Sugar())
	})

	It("Builds books list", func() {
		for lang := range config.SupportedLanguages {
			booksList, err := testStorage.ListBooks(lang)
			Expect(err).NotTo(HaveOccurred())
			Expect(booksList).NotTo(BeNil())
			Expect(len(booksList)).To(Equal(len(local.TitleToFileName[lang])))
		}
	})

	It("Takes random sentence from book", func() {
		for lang := range config.SupportedLanguages {
			for title := range local.TitleToFileName[lang] {
				quote, err := testStorage.GetRandomSentenceFromBook(title, lang, time.Now().UnixNano())
				Expect(err).NotTo(HaveOccurred())
				Expect(len(quote.Text)).To(BeNumerically(">", 20))
			}
		}
	})
})
