package local_test

import (
	"time"

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
		booksList, err := testStorage.ListBooks()
		Expect(err).NotTo(HaveOccurred())
		Expect(booksList).NotTo(BeNil())
		Expect(len(booksList)).To(Equal(len(local.TitleToFileName)))
	})

	It("Takes random sentence from book", func() {
		for title := range local.TitleToFileName {
			sentence, err := testStorage.GetRandomSentenceFromBook(title, time.Now().UnixNano())
			Expect(err).NotTo(HaveOccurred())
			Expect(len(sentence)).To(BeNumerically(">", 20))
		}
	})
})
