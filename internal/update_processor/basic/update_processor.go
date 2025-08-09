package basic

import (
	"fmt"
	"strings"
	"time"

	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/bookstorage/local"
	"github.com/mrkovshik/fortune_teller_bot/internal/update_processor"
	"go.uber.org/zap"
)

type UpdateProcessor struct {
	logger       *zap.SugaredLogger
	bookStorage  updateprocessor.BookStorage
	stateStorage updateprocessor.StateStorage
}

func NewUpdateProcessor(bookStorage updateprocessor.BookStorage, stateStorage updateprocessor.StateStorage, logger *zap.SugaredLogger) *UpdateProcessor {
	return &UpdateProcessor{
		logger:       logger,
		bookStorage:  bookStorage,
		stateStorage: stateStorage,
	}
}

func (cp *UpdateProcessor) ProcessMessage(message *model.Message) (map[string]interface{}, error) {
	chatID := message.Chat.ID
	command := message.Text
	payload := map[string]interface{}{
		"chat_id": chatID,
	}
	state, err := cp.stateStorage.Get(chatID)
	if err != nil {
		return nil, err
	}
	if state == nil || state.StepStack == nil {
		stepStack := model.NewStepStack()
		stepStack.Push(model.SelectStartCommand)
		cp.stateStorage.Update(chatID, &model.ChatState{
			StepStack: stepStack,
		})

		state, err = cp.stateStorage.Get(chatID)
		if err != nil {
			return nil, err
		}
	}
	currentStep, exist := state.StepStack.Peek()
	if !exist {
		cp.stateStorage.Clear(chatID)
		return nil, fmt.Errorf("can't find current step")
	}
	switch currentStep {
	case model.AskingQuestion:
		var seed int64
		for i := 0; i < len(command); i++ {
			seed += int64(command[i])
		}
		text, err := cp.bookStorage.GetRandomSentenceFromBook(local.GetRandomBookTitle(), seed)
		if err != nil {
			return nil, err
		}
		payload["text"] = text
		cp.stateStorage.Clear(chatID)
	default:
		state.StepStack = model.NewStepStack()
		state.StepStack.Push(model.SelectStartCommand)
		cp.stateStorage.Update(chatID, state)
		payload["text"] = "Что бы вы хотели сделать?"
		payload["reply_markup"] = startMenu
	}
	return payload, nil
}

func (cp *UpdateProcessor) ProcessCallback(callback *model.CallbackQuery) (map[string]interface{}, error) {
	chatID := callback.From.ID
	payload := map[string]interface{}{
		"chat_id": chatID,
	}
	state, err := cp.stateStorage.Get(chatID)
	if err != nil {
		return nil, err
	}
	if state == nil || state.StepStack == nil {
		cp.stateStorage.Update(chatID, &model.ChatState{
			StepStack: model.NewStepStack(),
		})
		state, err = cp.stateStorage.Get(chatID)
		if err != nil {
			return nil, err
		}
	}

	currentStep, exist := state.StepStack.Peek()
	if !exist {
		currentStep = model.SelectStartCommand
	}
	commandName := strings.TrimPrefix(callback.Data, string(currentStep))
	command := CallbackCommand(strings.TrimPrefix(commandName, ":"))
	switch currentStep {
	case model.SelectStartCommand:

		switch command {
		case GetRandomSentenceCommandName:
			payload["text"] = "Какую книгу вы хотите использовать для получения случайной цитаты?"
			payload["reply_markup"] = selectSourceMenu
			state.StepStack.Push(model.GetRandomSentenceMenu)
			cp.stateStorage.Update(chatID, state)

		case AskQuestionCommandName:
			payload["text"] = "Какую книгу вы хотите использовать для получения ответа на ваш вопрос?"
			payload["reply_markup"] = selectSourceMenu
			state.StepStack.Push(model.AskingQuestionMenu)
			cp.stateStorage.Update(chatID, state)

		case HelpCommandName:
			payload["text"] = "Есть такая народная забава - гадать на книгах. Человек мысленно или вслух задает вопрос, потом говорит случайную страницу и строку, и книга дает ему ответ. " +
				"Здесь все почти так же) Вы можете задать свой вопрос текстом - тогда бот использует этот текст для генерации случайных чисел, а можете просто получить случайную цитату из выбранной книги.\n\n" +
				"Что бы вы хотели сделать?"
			state.StepStack = model.NewStepStack()
			state.StepStack.Push(model.SelectStartCommand)
			cp.stateStorage.Update(chatID, state)
			payload["reply_markup"] = startMenu
		default:
			payload["text"] = "Так оно не работает. Попробуйте начать заново, нажав /start"
			cp.stateStorage.Clear(chatID)
		}
	case model.SelectBook:
		prevStep, exist := state.StepStack.PeekPrevious()
		if !exist {
			return nil, err
		}
		switch prevStep {
		case model.AskingQuestionMenu:
			payload["text"] = "Напишите вопрос, на который бы хотели получить ответ из книги, и мы используем его, как базу для поиска предсказания"
			state.StepStack.Push(model.AskingQuestion)
			cp.stateStorage.Update(chatID, state)
		case model.GetRandomSentenceMenu:
			fileName := strings.TrimPrefix(callback.Data, string(model.SelectBook))
			fileName = strings.TrimPrefix(fileName, ":")
			bookTitle, exist := local.FileNameToTitle[fileName]
			if !exist {
				payload["text"] = fmt.Sprintf("Книга с таким именем файла не найдена: %s", fileName)
				break
			}

			text, err := cp.bookStorage.GetRandomSentenceFromBook(bookTitle, time.Now().UnixNano())
			if err != nil {
				return nil, err
			}
			payload["text"] = text
			cp.stateStorage.Clear(chatID)
		default:
			payload["text"] = "Так оно не работает. Попробуйте начать заново, нажав /start"
			cp.stateStorage.Clear(chatID)
		}

	case model.AskingQuestionMenu:
		switch command {
		case ListBooksCommandName:
			payload["text"] = "Из каких книг вы хотите получить предсказание?"
			menu, err := cp.generateListBooksMenu()
			if err != nil {
				return nil, err
			}
			payload["reply_markup"] = menu
			state.StepStack.Push(model.SelectBook)
			cp.stateStorage.Update(chatID, state)
		case UseRandomBookCommandName:
			payload["text"] = "Напишите вопрос, на который бы хотели получить ответ из книги, и мы используем его, как базу для поиска предсказания"
			state.StepStack.Push(model.AskingQuestion)
			cp.stateStorage.Update(chatID, state)
		case GoBackCommandName:
			state.StepStack = model.NewStepStack()
			state.StepStack.Push(model.SelectStartCommand)
			cp.stateStorage.Update(chatID, state)
			payload["text"] = "Возвращаемся назад"
			payload["reply_markup"] = startMenu
		default:
			payload["text"] = "Так оно не работает. Попробуйте начать заново, нажав /start"
			cp.stateStorage.Clear(chatID)
		}
	case model.GetRandomSentenceMenu:
		switch command {
		case ListBooksCommandName:
			payload["text"] = "Из каких книг вы хотите получить предсказание?"
			menu, err := cp.generateListBooksMenu()
			if err != nil {
				return nil, err
			}
			payload["reply_markup"] = menu
			state.StepStack.Push(model.SelectBook)
			cp.stateStorage.Update(chatID, state)
		case UseRandomBookCommandName:
			text, err := cp.bookStorage.GetRandomSentenceFromBook(local.GetRandomBookTitle(), time.Now().UnixNano())
			if err != nil {
				return nil, err
			}
			if len(text) == 0 {
				text = "Извините, не получилось предсказать будущее"
			}
			payload["text"] = text
			cp.stateStorage.Clear(chatID)
		case GoBackCommandName:
			state.StepStack = model.NewStepStack()
			state.StepStack.Push(model.SelectStartCommand)
			cp.stateStorage.Update(chatID, state)
			payload["text"] = "Возвращаемся назад"
			payload["reply_markup"] = startMenu
		default:
			payload["text"] = "Так оно не работает. Попробуйте начать заново, нажав /start"
			cp.stateStorage.Clear(chatID)
		}
	}

	return payload, nil
}

func (cp *UpdateProcessor) generateListBooksMenu() (*InlineKeyboardMarkup, error) {
	books, err := cp.bookStorage.ListBooks()
	var keyboard [][]InlineKeyboardButton
	if err != nil {
		return nil, fmt.Errorf(`failed to list books: %w`, err)
	}
	for _, book := range books {
		button := InlineKeyboardButton{
			Text:         book,
			CallbackData: CallbackCommand(fmt.Sprintf("%s:%s", model.SelectBook, local.TitleToFileName[book])),
		}
		keyboard = append(keyboard, []InlineKeyboardButton{button})
	}
	return &InlineKeyboardMarkup{InlineKeyboard: keyboard}, nil
}
