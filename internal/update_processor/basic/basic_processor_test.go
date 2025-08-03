package basic_test

import (
	"github.com/mrkovshik/fortune_teller_bot/internal/command_processor"
	"github.com/mrkovshik/fortune_teller_bot/internal/command_processor/basic"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/book_storage/local"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("Local bookStorage functions test", func() {
	var (
		logger        *zap.Logger
		err           error
		testStorage   *local.Storage
		testProcessor *basic.UpdateProcessor
	)
	BeforeEach(func() {
		logger, err = zap.NewDevelopment()
		Expect(err).NotTo(HaveOccurred())
		testStorage = local.NewStorage(logger.Sugar()) // TODO: use mock
		testProcessor = basic.NewUpdateProcessor(logger.Sugar(), testStorage)
	})

	It("", func() {
		list, err := testProcessor.ProcessUpdate(update_processor.ListBooksCommandName)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(list)).To(BeNumerically(">", 0))
	})

	It("Takes random sentence from book", func() {
		sentence, err := testProcessor.ProcessUpdate(update_processor.GetMagicCommandName)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(sentence)).To(BeNumerically(">", 20))
	})
})
