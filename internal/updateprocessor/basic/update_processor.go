package basic

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/fortune_teller_bot/internal/embedded/templates"
	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/steps"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/userdata"
	"github.com/mrkovshik/fortune_teller_bot/internal/updateprocessor"
	"go.uber.org/zap"
)

type UpdateProcessor struct {
	logger          *zap.SugaredLogger
	cfg             *config.Config
	bookStorage     updateprocessor.BookStorage
	stepStorage     updateprocessor.StepStorage
	userDataStorage updateprocessor.UserDataStorage
}

func NewUpdateProcessor(bookStorage updateprocessor.BookStorage,
	stepStack updateprocessor.StepStorage,
	userDataStorage updateprocessor.UserDataStorage,
	logger *zap.SugaredLogger,
	cfg *config.Config) *UpdateProcessor {
	return &UpdateProcessor{
		logger:          logger,
		bookStorage:     bookStorage,
		stepStorage:     stepStack,
		userDataStorage: userDataStorage,
		cfg:             cfg,
	}
}

func (cp *UpdateProcessor) ProcessMessage(message *model.Message) (map[string]interface{}, error) {
	chatID := message.Chat.ID
	userData, err := cp.userDataStorage.GetUserData(chatID)
	if err != nil {
		return nil, err
	}
	var userLang string
	userLangRaw, ok := userData[userdata.LanguageKey]
	if !ok {
		userLang = cp.cfg.DefaultLanguage
	} else {
		userLang, ok = userLangRaw.(string)
		if !ok {
			return nil, fmt.Errorf("user language value is not a string")
		}
	}

	command := message.Text
	payload := map[string]interface{}{
		"chat_id": chatID,
	}
	currentStep, err := cp.stepStorage.Peek(chatID)
	if err != nil {
		if !errors.Is(err, steps.ErrChatNotFound) {
			return nil, err
		}
		if err := cp.stepStorage.Push(chatID, steps.SelectStartCommand); err != nil {
			return nil, err
		}
		currentStep = steps.SelectStartCommand
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
			idx, ok := userData[userdata.BookTitleKey]
			if !ok {
				return nil, fmt.Errorf("no title key found in userData for chatID %d", chatID)
			}
			books, err := cp.bookStorage.ListBooks()
			if err != nil {
				return nil, fmt.Errorf(`failed to list books: %w`, err)
			}
			idxInt, ok := idx.(int)
			if !ok {
				return nil, fmt.Errorf("invalid book index for chatID %d", chatID)
			}
			title = books[idxInt]
		case steps.AskingQuestionMenu:
			title = cp.bookStorage.GetRandomBookTitle()
		}

		quote, err := cp.bookStorage.GetRandomSentenceFromBook(title, seed)
		if err != nil {
			return nil, err
		}
		msg, err := templates.GenerateMessageWithData(templates.QuoteTemplateName, quote, userLang)
		if err != nil {
			return nil, err
		}
		payload["text"] = msg
		cp.stepStorage.Clear(chatID)

	case steps.SelectStartCommand:
		payload["text"] = templates.SimpleMessages[templates.StartTemplateName]
		payload["reply_markup"] = StartMenu
	}
	return payload, nil
}

func (cp *UpdateProcessor) ProcessCallback(callback *model.CallbackQuery) (map[string]interface{}, error) {
	chatID := callback.From.ID
	userData, err := cp.userDataStorage.GetUserData(chatID)
	if err != nil {
		return nil, err
	}
	var userLang string
	userLangRaw, ok := userData[userdata.LanguageKey]
	if !ok {
		userLang = cp.cfg.DefaultLanguage
	} else {
		userLang, ok = userLangRaw.(string)
		if !ok {
			return nil, fmt.Errorf("user language value is not a string")
		}
	}

	payload := map[string]interface{}{
		"chat_id": chatID,
	}
	currentStep, err := cp.stepStorage.Peek(chatID)
	if err != nil {
		if !errors.Is(err, steps.ErrChatNotFound) {
			return nil, err
		}
		if err := cp.stepStorage.Push(chatID, steps.SelectStartCommand); err != nil {
			return nil, err
		}
		currentStep = steps.SelectStartCommand
	}

	commandName := strings.TrimPrefix(callback.Data, string(currentStep))
	command := CallbackCommand(strings.TrimPrefix(commandName, ":"))
	switch currentStep {
	case steps.SelectStartCommand:

		switch command {
		case GetRandomSentenceCommandName:
			payload["text"] = templates.SimpleMessages[templates.SelectBookForRandomTemplateName]
			payload["reply_markup"] = selectSourceMenu
			if err := cp.stepStorage.Push(chatID, steps.GetRandomSentenceMenu); err != nil {
				return nil, err
			}

		case AskQuestionCommandName:
			payload["text"] = templates.SimpleMessages[templates.SelectBookForQuestionTemplateName]
			payload["reply_markup"] = selectSourceMenu
			if err := cp.stepStorage.Push(chatID, steps.AskingQuestionMenu); err != nil {
				return nil, err
			}

		case HelpCommandName:
			payload["text"] = templates.SimpleMessages[templates.HelpTemplateName]
			payload["reply_markup"] = StartMenu
		default:
			payload["text"] = templates.SimpleMessages[templates.InvalidButtonTemplateName]
		}
	case steps.SelectBook:
		prevStep, err := cp.stepStorage.PeekPrevious(chatID)
		if err != nil {
			return nil, err
		}
		switch prevStep {
		case steps.AskingQuestionMenu:
			if err := cp.userDataStorage.SaveUserData(chatID, userdata.BookTitleKey, string(command)); err != nil {
				return nil, err
			}
			payload["text"] = templates.SimpleMessages[templates.TypeQuestionTemplateName]
			if err := cp.stepStorage.Push(chatID, steps.AskingQuestion); err != nil {
				return nil, err
			}
		case steps.GetRandomSentenceMenu:
			books, err := cp.bookStorage.ListBooks()
			if err != nil {
				return nil, fmt.Errorf(`failed to list books: %w`, err)
			}
			idx, err := strconv.Atoi(callback.Data)
			if err != nil {
				return nil, err
			}
			quote, err := cp.bookStorage.GetRandomSentenceFromBook(books[idx], time.Now().UnixNano())
			if err != nil {
				return nil, err
			}
			msg, err := templates.GenerateMessageWithData(templates.QuoteTemplateName, quote, userLang)
			if err != nil {
				return nil, err
			}
			payload["text"] = msg
			cp.stepStorage.Clear(chatID)
		default:
			payload["text"] = templates.SimpleMessages[templates.InvalidButtonTemplateName]
		}

	case steps.AskingQuestionMenu:
		switch command {
		case ListBooksCommandName:
			payload, err = cp.generateListBooksMenuPayload(chatID, userLang)
			if err != nil {
				return nil, err
			}
			if err := cp.stepStorage.Push(chatID, steps.SelectBook); err != nil {
				return nil, err
			}
		case UseRandomBookCommandName:
			payload["text"] = templates.SimpleMessages[templates.TypeQuestionTemplateName]
			if err := cp.stepStorage.Push(chatID, steps.AskingQuestion); err != nil {
				return nil, err
			}
		case GoBackCommandName:
			cp.stepStorage.Clear(chatID)
			payload["text"] = "Возвращаемся назад"
			payload["reply_markup"] = StartMenu
		default:
			payload["text"] = templates.SimpleMessages[templates.InvalidButtonTemplateName]
		}
	case steps.GetRandomSentenceMenu:
		switch command {
		case ListBooksCommandName:
			payload, err = cp.generateListBooksMenuPayload(chatID, userLang)
			if err != nil {
				return nil, err
			}
			if err := cp.stepStorage.Push(chatID, steps.SelectBook); err != nil {
				return nil, err
			}
		case UseRandomBookCommandName:
			quote, err := cp.bookStorage.GetRandomSentenceFromBook(cp.bookStorage.GetRandomBookTitle(), time.Now().UnixNano())
			if err != nil {
				return nil, err
			}
			msg, err := templates.GenerateMessageWithData(templates.QuoteTemplateName, quote, userLang)
			if err != nil {
				return nil, err
			}
			payload["text"] = msg
			cp.stepStorage.Clear(chatID)
		case GoBackCommandName:
			cp.stepStorage.Clear(chatID)
			payload["text"] = "Возвращаемся назад"
			payload["reply_markup"] = StartMenu
		default:
			payload["text"] = templates.SimpleMessages[templates.InvalidButtonTemplateName]
		}
	}

	return payload, nil
}

func (cp *UpdateProcessor) generateListBooksMenuPayload(chatID int64, lang string) (map[string]interface{}, error) {
	var keyboard [][]InlineKeyboardButton
	books, err := cp.bookStorage.ListBooks()
	if err != nil {
		return nil, fmt.Errorf(`failed to list books: %w`, err)
	}
	for i, book := range books {
		button := InlineKeyboardButton{
			Text:         book,
			CallbackData: CallbackCommand(strconv.Itoa(i)),
		}
		keyboard = append(keyboard, []InlineKeyboardButton{button})
	}
	payload := map[string]interface{}{
		"chat_id":      chatID,
		"text":         templates.SimpleMessages[lang][templates.ListBooksTemplateName],
		"reply_markup": &InlineKeyboardMarkup{InlineKeyboard: keyboard},
	}
	return payload, nil
}
