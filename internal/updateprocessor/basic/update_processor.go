package basic

import (
	"fmt"
	"strings"
	"time"

	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/bookstorage/local"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/steps"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/userdata"
	"github.com/mrkovshik/fortune_teller_bot/internal/updateprocessor"
	"go.uber.org/zap"
)

type UpdateProcessor struct {
	logger          *zap.SugaredLogger
	bookStorage     updateprocessor.BookStorage
	stepStorage     updateprocessor.StepStorage
	userDataStorage updateprocessor.UserDataStorage
}

func NewUpdateProcessor(bookStorage updateprocessor.BookStorage, stepStack updateprocessor.StepStorage, userDataStorage updateprocessor.UserDataStorage, logger *zap.SugaredLogger) *UpdateProcessor {
	return &UpdateProcessor{
		logger:          logger,
		bookStorage:     bookStorage,
		stepStorage:     stepStack,
		userDataStorage: userDataStorage,
	}
}

func (cp *UpdateProcessor) ProcessMessage(message *model.Message) (map[string]interface{}, error) {
	chatID := message.Chat.ID
	command := message.Text
	payload := map[string]interface{}{
		"chat_id": chatID,
	}
	currentStep, err := cp.stepStorage.Peek(chatID)
	if err != nil {
		return nil, err
	}
	if currentStep == "" {
		if err := cp.stepStorage.Push(chatID, steps.SelectStartCommand); err != nil {
			return nil, err
		}
	}
	switch currentStep {
	case steps.AskingQuestion:
		var seed int64
		for i := 0; i < len(command); i++ {
			seed += int64(command[i])
		}
		prevStep, err := cp.stepStorage.PeekPrevious(chatID)
		if err != nil {
			return nil, err
		}
		var title string
		switch prevStep {
		case steps.SelectBook:
			title, err = cp.userDataStorage.GetUserData(chatID, userdata.BookTitleKey)
			if err != nil {
				return nil, err
			}
		case steps.AskingQuestionMenu:
			title = cp.bookStorage.GetRandomBookTitle()
		}

		text, err := cp.bookStorage.GetRandomSentenceFromBook(title, seed)
		if err != nil {
			return nil, err
		}
		payload["text"] = text
		cp.stepStorage.Clear(chatID)

	case steps.SelectStartCommand:
		payload["text"] = "Что бы вы хотели сделать?"
		payload["reply_markup"] = StartMenu
	}
	return payload, nil
}

func (cp *UpdateProcessor) ProcessCallback(callback *model.CallbackQuery) (map[string]interface{}, error) {
	chatID := callback.From.ID
	payload := map[string]interface{}{
		"chat_id": chatID,
	}
	currentStep, err := cp.stepStorage.Peek(chatID)
	if err != nil {
		return nil, err
	}
	if currentStep == "" {
		if err := cp.stepStorage.Push(chatID, steps.SelectStartCommand); err != nil {
			return nil, err
		}
	}

	commandName := strings.TrimPrefix(callback.Data, string(currentStep))
	command := CallbackCommand(strings.TrimPrefix(commandName, ":"))
	switch currentStep {
	case steps.SelectStartCommand:

		switch command {
		case GetRandomSentenceCommandName:
			payload["text"] = "Какую книгу вы хотите использовать для получения случайной цитаты?"
			payload["reply_markup"] = selectSourceMenu
			if err := cp.stepStorage.Push(chatID, steps.GetRandomSentenceMenu); err != nil {
				return nil, err
			}

		case AskQuestionCommandName:
			payload["text"] = "Какую книгу вы хотите использовать для получения ответа на ваш вопрос?"
			payload["reply_markup"] = selectSourceMenu
			if err := cp.stepStorage.Push(chatID, steps.AskingQuestionMenu); err != nil {
				return nil, err
			}

		case HelpCommandName:
			payload["text"] = "Есть такая народная забава - гадать на книгах. Человек мысленно или вслух задает вопрос, потом говорит случайную страницу и строку, и книга дает ему ответ. " +
				"Здесь все почти так же) Вы можете задать свой вопрос текстом - тогда бот использует этот текст для генерации случайных чисел, а можете просто получить случайную цитату из выбранной книги.\n\n" +
				"Что бы вы хотели сделать?"
			payload["reply_markup"] = StartMenu
		default:
			payload["text"] = "Так оно не работает. Используйте последнее меню или начните заново, нажав /start"
		}
	case steps.SelectBook:
		prevStep, err := cp.stepStorage.PeekPrevious(chatID)
		if err != nil {
			return nil, err
		}
		switch prevStep {
		case steps.AskingQuestionMenu:
			payload["text"] = "Напишите вопрос, на который бы хотели получить ответ из книги, и мы используем его, как базу для поиска предсказания"
			if err := cp.stepStorage.Push(chatID, steps.AskingQuestion); err != nil {
				return nil, err
			}
		case steps.GetRandomSentenceMenu:

			text, err := cp.bookStorage.GetRandomSentenceFromBook(callback.Data, time.Now().UnixNano())
			if err != nil {
				return nil, err
			}
			payload["text"] = text
			cp.stepStorage.Clear(chatID)
		default:
			payload["text"] = "Так оно не работает. Используйте последнее меню или начните заново, нажав /start"
		}

	case steps.AskingQuestionMenu:
		switch command {
		case ListBooksCommandName:
			payload, err = cp.generateListBooksMenuPayload(chatID)
			if err != nil {
				return nil, err
			}
			if err := cp.stepStorage.Push(chatID, steps.SelectBook); err != nil {
				return nil, err
			}
		case UseRandomBookCommandName:
			payload["text"] = "Напишите вопрос, на который бы хотели получить ответ из книги, и мы используем его, как базу для поиска предсказания"
			if err := cp.stepStorage.Push(chatID, steps.AskingQuestion); err != nil {
				return nil, err
			}
		case GoBackCommandName:
			cp.stepStorage.Clear(chatID)
			payload["text"] = "Возвращаемся назад"
			payload["reply_markup"] = StartMenu
		default:
			payload["text"] = "Так оно не работает. Используйте последнее меню или начните заново, нажав /start"
		}
	case steps.GetRandomSentenceMenu:
		switch command {
		case ListBooksCommandName:
			payload, err = cp.generateListBooksMenuPayload(chatID)
			if err != nil {
				return nil, err
			}
			if err := cp.stepStorage.Push(chatID, steps.SelectBook); err != nil {
				return nil, err
			}
		case UseRandomBookCommandName:
			text, err := cp.bookStorage.GetRandomSentenceFromBook(local.GetRandomBookTitle(), time.Now().UnixNano())
			if err != nil {
				return nil, err
			}
			if len(text) == 0 {
				text = "Извините, не получилось предсказать будущее"
			}
			payload["text"] = text
			cp.stepStorage.Clear(chatID)
		case GoBackCommandName:
			cp.stepStorage.Clear(chatID)
			payload["text"] = "Возвращаемся назад"
			payload["reply_markup"] = StartMenu
		default:
			payload["text"] = "Так оно не работает. Используйте последнее меню или начните заново, нажав /start"
		}
	}

	return payload, nil
}

func (cp *UpdateProcessor) generateListBooksMenuPayload(chatID int64) (map[string]interface{}, error) {
	books, err := cp.bookStorage.ListBooks()
	var keyboard [][]InlineKeyboardButton
	if err != nil {
		return nil, fmt.Errorf(`failed to list books: %w`, err)
	}
	for _, book := range books {
		button := InlineKeyboardButton{
			Text:         book,
			CallbackData: CallbackCommand(local.TitleToFileName[book]),
		}
		keyboard = append(keyboard, []InlineKeyboardButton{button})
	}
	payload := map[string]interface{}{
		"chat_id":      chatID,
		"text":         "Из каких книг вы хотите получить предсказание?",
		"reply_markup": &InlineKeyboardMarkup{InlineKeyboard: keyboard},
	}
	return payload, nil
}
