package local_test

import (
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/local"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

const dorianGreyTitle = "Оскар Уайлд - Портрет Дориана Грея"

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
		Expect(len(booksList)).To(Equal(1))
	})

	It("Takes random sentence from book", func() {
		sentence, err := testStorage.GetRandomSentenceFromBook(dorianGreyTitle)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(sentence)).To(BeNumerically(">", 20))
	})
})
