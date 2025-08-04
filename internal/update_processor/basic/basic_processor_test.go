package basic_test

import (
	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/book_storage/local"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/state_storage/in_memory"
	"github.com/mrkovshik/fortune_teller_bot/internal/update_processor"
	"github.com/mrkovshik/fortune_teller_bot/internal/update_processor/basic"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("Local bookStorage functions test", func() {
	var (
		logger        *zap.Logger
		err           error
		bookStorage   *local.Storage
		stateStorage  in_memory.StateStorage
		testProcessor *basic.UpdateProcessor
	)
	BeforeEach(func() {
		logger, err = zap.NewDevelopment()
		Expect(err).NotTo(HaveOccurred())
		stateStorage = in_memory.NewStateStorage()
		bookStorage = local.NewStorage(logger.Sugar()) // TODO: use mock
		testProcessor = basic.NewUpdateProcessor(bookStorage, stateStorage, logger.Sugar())
	})

	It("", func() {
		reply, err := testProcessor.ProcessCallback(&model.CallbackQuery{
			ID: "123",
			From: model.Chat{
				ID: 132,
			},
			Data: "2",
		})
		Expect(err).NotTo(HaveOccurred())
		sentence, ok := reply["text"].(string)
		Expect(ok).To(BeTrue())
		Expect(len(sentence)).To(BeNumerically(">", 20))
	})

	It("Takes random sentence from book", func() {
		reply, err := testProcessor.ProcessMessage(&model.Message{
			Chat: model.Chat{
				ID: 123,
			},
			Text: update_processor.GetMagicCommandName,
		})
		Expect(err).NotTo(HaveOccurred())
		sentence, ok := reply["text"].(string)
		Expect(ok).To(BeTrue())
		Expect(len(sentence)).To(BeNumerically(">", 20))
	})

	It("Takes random sentence from book", func() {
		reply, err := testProcessor.ProcessMessage(&model.Message{
			Chat: model.Chat{
				ID: 123,
			},
			Text: update_processor.ListBooksCommandName,
		})
		Expect(err).NotTo(HaveOccurred())
		sentence, ok := reply["text"].(string)
		Expect(ok).To(BeTrue())
		Expect(sentence).To(ContainSubstring("Выберите книгу, по которой будем предсказывать будущее:"))
		keyBoard, ok := reply["reply_markup"].(basic.InlineKeyboardMarkup)
		Expect(ok).To(BeTrue())
		Expect(len(keyBoard.InlineKeyboard)).To(BeNumerically(">", 1))
	})
})
