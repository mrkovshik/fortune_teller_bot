package basic_test

import (
	"github.com/golang/mock/gomock"

	"github.com/mrkovshik/fortune_teller_bot/internal/embedded/templates"
	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/steps"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/userdata"

	"github.com/mrkovshik/fortune_teller_bot/internal/updateprocessor"
	"github.com/mrkovshik/fortune_teller_bot/internal/updateprocessor/basic"
	mock "github.com/mrkovshik/fortune_teller_bot/mocks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

const (
	testQuote     = "Awesome test quote of destiny"
	testChatID    = int64(111)
	testBookTitle = "awesome book 2"
	testBookIdx   = "1"
)

var _ = Describe("Local bookStorage functions test", func() {
	var (
		logger          *zap.Logger
		ctrl            *gomock.Controller
		err             error
		bookStorage     *mock.MockBookStorage
		stepStorage     *mock.MockStepStorage
		userDataStorage *mock.MockUserDataStorage
		testProcessor   updateprocessor.UpdateProcessor
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		bookStorage = mock.NewMockBookStorage(ctrl)
		stepStorage = mock.NewMockStepStorage(ctrl)
		userDataStorage = mock.NewMockUserDataStorage(ctrl)
		logger, err = zap.NewDevelopment()
		testProcessor = basic.NewUpdateProcessor(bookStorage, stepStorage, userDataStorage, logger.Sugar())
		Expect(err).NotTo(HaveOccurred())
		Expect(templates.InitTemplates("rus")).To(Succeed())

	})
	AfterEach(func() {
		DeferCleanup(ctrl.Finish)
	})

	It("Get random sentence from specific book", func() {
		stepStorage.EXPECT().Peek(testChatID).Return(steps.SelectBook, nil)
		stepStorage.EXPECT().PeekPrevious(testChatID).Return(steps.GetRandomSentenceMenu, nil)
		bookStorage.EXPECT().GetRandomSentenceFromBook(testBookTitle, gomock.Any()).Return(testQuote, nil)
		stepStorage.EXPECT().Clear(testChatID)
		bookStorage.EXPECT().ListBooks().Return([]string{"some book 1", testBookTitle, "some book 3"}, nil)
		reply, err := testProcessor.ProcessCallback(&model.CallbackQuery{
			ID: "123",
			From: model.Chat{
				ID: testChatID,
			},
			Data: testBookIdx,
		})
		Expect(err).NotTo(HaveOccurred())
		sentence, ok := reply["text"].(string)
		Expect(ok).To(BeTrue())
		Expect(sentence).To(Equal(testQuote))
	})

	It("Takes answer from specific book", func() {
		bookStorage.EXPECT().ListBooks().Return([]string{"some book 1", testBookTitle, "some book 3"}, nil)
		stepStorage.EXPECT().Peek(testChatID).Return(steps.AskingQuestion, nil)
		stepStorage.EXPECT().PeekPrevious(testChatID).Return(steps.SelectBook, nil)
		stepStorage.EXPECT().Clear(testChatID)
		bookStorage.EXPECT().GetRandomSentenceFromBook(testBookTitle, gomock.Any()).Return(testQuote, nil)
		userDataStorage.EXPECT().GetUserData(testChatID, userdata.BookTitleKey).Return(testBookIdx, nil)
		reply, err := testProcessor.ProcessMessage(&model.Message{
			Chat: model.Chat{
				ID: testChatID,
			},
			Text: "Some random text",
		})
		Expect(err).NotTo(HaveOccurred())
		sentence, ok := reply["text"].(string)
		Expect(ok).To(BeTrue())
		Expect(sentence).To(Equal(testQuote))
	})

	It("Takes random sentence from random book", func() {
		stepStorage.EXPECT().Peek(testChatID).Return(steps.GetRandomSentenceMenu, nil)
		stepStorage.EXPECT().Push(testChatID, gomock.Any()).Return(nil)
		bookStorage.EXPECT().ListBooks().Return([]string{"some book 1", "some book 2"}, nil)

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
		Expect(keyBoard.InlineKeyboard[0][0].Text).To(Equal("some book 1"))
		Expect(len(keyBoard.InlineKeyboard)).To(BeNumerically(">", 1))
	})
})
