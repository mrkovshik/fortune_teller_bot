package basic_test

import (
	"strconv"

	"github.com/golang/mock/gomock"
	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/userdata"

	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/steps"
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
	testBookIdx   = 1
)

var _ = Describe("Local bookStorage functions test", func() {
	var (
		logger            *zap.Logger
		ctrl              *gomock.Controller
		bookStorage       *mock.MockBookStorage
		stepStorage       *mock.MockStepStorage
		userDataStorage   *mock.MockUserDataStorage
		testProcessor     updateprocessor.UpdateProcessor
		testBookIdxString = strconv.Itoa(testBookIdx)
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		bookStorage = mock.NewMockBookStorage(ctrl)
		stepStorage = mock.NewMockStepStorage(ctrl)
		userDataStorage = mock.NewMockUserDataStorage(ctrl)
		cfg, err := config.GetConfig()
		Expect(err).NotTo(HaveOccurred())
		logger, err = zap.NewDevelopment()
		Expect(err).NotTo(HaveOccurred())
		testProcessor = basic.NewUpdateProcessor(bookStorage, stepStorage, userDataStorage, logger.Sugar(), cfg)
		Expect(err).NotTo(HaveOccurred())
		userDataStorage.EXPECT().GetUserData(testChatID).Return(userdata.UserData{
			userdata.LanguageKey:  config.Russian,
			userdata.BookTitleKey: testBookIdx,
		}, nil).AnyTimes()

	})
	AfterEach(func() {
		DeferCleanup(ctrl.Finish)
	})

	It("Get random sentence from specific book", func() {
		stepStorage.EXPECT().Push(testChatID, steps.SelectStartCommand).Return(nil)
		stepStorage.EXPECT().Peek(testChatID).Return(steps.SelectBook, nil)
		stepStorage.EXPECT().PeekPrevious(testChatID).Return(steps.GetRandomSentenceMenu, nil)
		bookStorage.EXPECT().
			GetRandomSentenceFromBook(testBookTitle, config.Russian, gomock.Any()).
			Return(&updateprocessor.Quote{
				Title: testBookTitle,
				Text:  testQuote,
			}, nil)
		stepStorage.EXPECT().Clear(testChatID)
		bookStorage.EXPECT().ListBooks(config.Russian).Return([]string{"some book 1", testBookTitle, "some book 3"}, nil)
		reply, err := testProcessor.ProcessCallback(&model.CallbackQuery{
			ID: "123",
			From: &model.User{
				ID: testChatID,
			},
			Data: testBookIdxString,
		})
		Expect(err).NotTo(HaveOccurred())
		sentence, ok := reply["text"].(string)
		Expect(ok).To(BeTrue())
		Expect(sentence).To(ContainSubstring(testQuote))
		Expect(sentence).To(ContainSubstring(testBookTitle))
	})

	It("Takes answer from specific book", func() {
		stepStorage.EXPECT().Push(testChatID, steps.SelectStartCommand).Return(nil)
		bookStorage.EXPECT().ListBooks(config.Russian).Return([]string{"some book 1", testBookTitle, "some book 3"}, nil)
		stepStorage.EXPECT().Peek(testChatID).Return(steps.AskingQuestion, nil)
		stepStorage.EXPECT().PeekPrevious(testChatID).Return(steps.SelectBook, nil)
		stepStorage.EXPECT().Clear(testChatID)
		bookStorage.EXPECT().
			GetRandomSentenceFromBook(testBookTitle, config.Russian, gomock.Any()).
			Return(&updateprocessor.Quote{
				Title: testBookTitle,
				Text:  testQuote,
			}, nil)
		reply, err := testProcessor.ProcessMessage(&model.Message{
			Chat: model.Chat{
				ID: testChatID,
			},
			Text: "Some random text",
		})
		Expect(err).NotTo(HaveOccurred())
		sentence, ok := reply["text"].(string)
		Expect(ok).To(BeTrue())
		Expect(sentence).To(ContainSubstring(testQuote))
		Expect(sentence).To(ContainSubstring(testBookTitle))
	})

	It("Takes random sentence from random book", func() {
		stepStorage.EXPECT().Peek(testChatID).Return(steps.GetRandomSentenceMenu, nil)
		stepStorage.EXPECT().Push(testChatID, gomock.Any()).Return(nil)
		bookStorage.EXPECT().ListBooks(config.Russian).Return([]string{"some book 1", "some book 2"}, nil)

		reply, err := testProcessor.ProcessCallback(&model.CallbackQuery{
			ID: "123",
			From: &model.User{
				ID: testChatID,
			},
			Data: string(model.ListBooksCommandName),
		})
		Expect(err).NotTo(HaveOccurred())
		sentence, ok := reply["text"].(string)
		Expect(ok).To(BeTrue())
		Expect(sentence).To(ContainSubstring("Из каких книг вы хотите получить предсказание?"))
		keyBoard, ok := reply["reply_markup"].(*model.InlineKeyboardMarkup)
		Expect(ok).To(BeTrue())
		Expect(keyBoard.InlineKeyboard[0][0].Text).To(Equal("some book 1"))
		Expect(len(keyBoard.InlineKeyboard)).To(BeNumerically(">", 1))
	})

	It("Get language menu", func() {
		stepStorage.EXPECT().Peek(testChatID).Return(steps.SelectStartCommand, nil)
		stepStorage.EXPECT().Push(testChatID, steps.SelectLanguage).Return(nil)

		reply, err := testProcessor.ProcessCallback(&model.CallbackQuery{
			ID: "123",
			From: &model.User{
				ID: testChatID,
			},
			Data: string(model.LanguageCommandName),
		})
		Expect(err).NotTo(HaveOccurred())
		sentence, ok := reply["text"].(string)
		Expect(ok).To(BeTrue())
		Expect(sentence).To(ContainSubstring("Вот языки, которые поддерживает наш бот:"))
		keyBoard, ok := reply["reply_markup"].(*model.InlineKeyboardMarkup)
		Expect(ok).To(BeTrue())
		for lang, langName := range config.SupportedLanguages {
			Expect(keyBoard.InlineKeyboard).To(ContainElement([]model.InlineKeyboardButton{{
				Text:         langName,
				CallbackData: model.CallbackCommand(lang),
			}}))
		}

		Expect(len(keyBoard.InlineKeyboard)).To(Equal(len(config.SupportedLanguages)))
	})

	It("Changing lang", func() {
		stepStorage.EXPECT().Push(testChatID, steps.SelectStartCommand).Return(nil)
		stepStorage.EXPECT().Peek(testChatID).Return(steps.SelectLanguage, nil)
		stepStorage.EXPECT().Clear(testChatID)
		userDataStorage.EXPECT().SaveUserData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		reply, err := testProcessor.ProcessCallback(&model.CallbackQuery{
			ID: "123",
			From: &model.User{
				ID: testChatID,
			},
			Data: string(config.Russian),
		})
		Expect(err).NotTo(HaveOccurred())
		sentence, ok := reply["text"].(string)
		Expect(ok).To(BeTrue())
		Expect(sentence).To(ContainSubstring("Язык бота был изменен на Русский"))
	})

	It("Get source menu", func() {
		stepStorage.EXPECT().Peek(testChatID).Return(steps.SelectStartCommand, nil)
		stepStorage.EXPECT().Push(testChatID, steps.AskingQuestionMenu).Return(nil)

		reply, err := testProcessor.ProcessCallback(&model.CallbackQuery{
			ID: "123",
			From: &model.User{
				ID: testChatID,
			},
			Data: string(model.AskQuestionCommandName),
		})
		Expect(err).NotTo(HaveOccurred())
		sentence, ok := reply["text"].(string)
		Expect(ok).To(BeTrue())
		Expect(sentence).To(ContainSubstring("Какую книгу вы хотите использовать для получения ответа на ваш вопрос?"))
		keyBoard, ok := reply["reply_markup"].(model.InlineKeyboardMarkup)
		Expect(ok).To(BeTrue())
		Expect(len(keyBoard.InlineKeyboard)).To(Equal(len(model.Menus["rus"][model.SelectSourceMenu].InlineKeyboard)))
	})
})
