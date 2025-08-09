package basic_test

import (
	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/book_storage/local"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/state_storage/in_memory"
	"github.com/mrkovshik/fortune_teller_bot/internal/update_processor/basic"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

const testChatID = 111

var _ = Describe("Local bookStorage functions test", func() {
	var (
		logger        *zap.Logger
		err           error
		bookStorage   *local.Storage
		stateStorage  inmemory.StateStorage
		testProcessor *basic.UpdateProcessor
		stepStack     = model.NewStepStack()
	)
	BeforeEach(func() {
		logger, err = zap.NewDevelopment()
		Expect(err).NotTo(HaveOccurred())
		stateStorage = inmemory.NewStateStorage() // TODO: use mock
		stepStack.Push(model.SelectStartCommand)
		stateStorage.Update(testChatID, &model.ChatState{
			StepStack: stepStack,
		})
		bookStorage = local.NewStorage(logger.Sugar()) // TODO: use mock
		testProcessor = basic.NewUpdateProcessor(bookStorage, stateStorage, logger.Sugar())
	})

	It("Get random sentence from specific book", func() {
		stepStack.Push(model.SelectBook)
		reply, err := testProcessor.ProcessCallback(&model.CallbackQuery{
			ID: "123",
			From: model.Chat{
				ID: testChatID,
			},
			Data: "3",
		})
		Expect(err).NotTo(HaveOccurred())
		sentence, ok := reply["text"].(string)
		Expect(ok).To(BeTrue())
		Expect(len(sentence)).To(BeNumerically(">", 20))
	})

	It("Takes answer from specific book", func() {
		stepStack.Push(model.SelectBook)
		stepStack.Push(model.AskingQuestion)
		reply, err := testProcessor.ProcessMessage(&model.Message{
			Chat: model.Chat{
				ID: testChatID,
			},
			Text: "Some random text",
		})
		Expect(err).NotTo(HaveOccurred())
		sentence, ok := reply["text"].(string)
		Expect(ok).To(BeTrue())
		Expect(len(sentence)).To(BeNumerically(">", 20))
	})

	It("Takes random sentence from random book", func() {
		reply, err := testProcessor.ProcessCallback(&model.CallbackQuery{
			ID: "123",
			From: model.Chat{
				ID: testChatID,
			},
			Data: string(basic.ListBooksCommandName),
		})
		Expect(err).NotTo(HaveOccurred())
		sentence, ok := reply["text"].(string)
		Expect(ok).To(BeTrue())
		Expect(sentence).To(ContainSubstring("Из каких книг вы хотите получить предсказание?"))
		keyBoard, ok := reply["reply_markup"].(*basic.InlineKeyboardMarkup)
		Expect(ok).To(BeTrue())
		Expect(len(keyBoard.InlineKeyboard)).To(BeNumerically(">", 1))
	})
})
