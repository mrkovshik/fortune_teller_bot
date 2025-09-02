package basic

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	booksmeta "github.com/mrkovshik/fortune_teller_bot/internal/storage/books-meta"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/steps"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/userdata"
	"github.com/mrkovshik/fortune_teller_bot/internal/templates"
	"github.com/mrkovshik/fortune_teller_bot/internal/updateprocessor"
	"go.uber.org/zap"
)

type UpdateProcessor struct {
	logger          *zap.SugaredLogger
	cfg             *config.Config
	bookStorage     updateprocessor.BookStorage
	booksRepository updateprocessor.BookRepository
	stepStorage     updateprocessor.StepStorage
	userDataStorage updateprocessor.UserDataStorage
}

func NewUpdateProcessor(bookStorage updateprocessor.BookStorage,
	booksRepository updateprocessor.BookRepository,
	stepStack updateprocessor.StepStorage,
	userDataStorage updateprocessor.UserDataStorage,
	logger *zap.SugaredLogger,
	cfg *config.Config) *UpdateProcessor {
	return &UpdateProcessor{
		logger:          logger,
		bookStorage:     bookStorage,
		booksRepository: booksRepository,
		stepStorage:     stepStack,
		userDataStorage: userDataStorage,
		cfg:             cfg,
	}
}

func (cp *UpdateProcessor) ProcessMessage(ctx context.Context, message *model.Message) (map[string]interface{}, error) {
	chatID := message.Chat.ID
	userData, err := cp.userDataStorage.GetUserData(chatID)
	if err != nil {
		return nil, err
	}
	var userLang config.Language
	userLangRaw, ok := userData[userdata.LanguageKey]
	if !ok {
		userLang = cp.cfg.DefaultLanguage
	} else {
		userLang, ok = userLangRaw.(config.Language)
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
		var book *booksmeta.Book
		switch prevStep {
		case steps.SelectBook:
			bookIDRaw, ok := userData[userdata.BookIDKey]
			if !ok {
				return nil, fmt.Errorf("no id key found in userData for chatID %d", chatID)
			}
			bookID, ok := bookIDRaw.(int64)
			if !ok {
				return nil, fmt.Errorf("invalid book index for chatID %d", chatID)
			}
			book, err = cp.bookStorage.GetBookByID(ctx, bookID)
			if err != nil {
				return nil, err
			}
		case steps.AskingQuestionMenu:
			book, err = cp.bookStorage.GetRandomBook(ctx, booksmeta.WithLanguage(userLang))
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown sequence: %s -> %s", prevStep, currentStep)
		}
		bookData, err := cp.booksRepository.LoadBook(book)
		if err != nil {
			return nil, err
		}
		parser, err := book.GetParser()
		if err != nil {
			return nil, err
		}
		sentence, err := parser.ParseRandomSentence(bookData, seed)
		if err != nil {
			return nil, err
		}

		msg, err := templates.GenerateMessageWithData(templates.QuoteTemplateName, updateprocessor.Quote{Text: sentence, Book: book}, userLang)
		if err != nil {
			return nil, err
		}
		payload["text"] = msg
		payload["reply_markup"] = model.Menus[userLang][model.StartMenu]
		cp.stepStorage.Clear(chatID)
		if err := cp.stepStorage.Push(chatID, steps.SelectStartCommand); err != nil {
			return nil, err
		}

	case steps.SelectStartCommand:
		payload["text"] = templates.SimpleMessages[userLang][templates.StartTemplateName]
		payload["reply_markup"] = model.Menus[userLang][model.StartMenu]
	}
	return payload, nil
}

func (cp *UpdateProcessor) ProcessCallback(ctx context.Context, callback *model.CallbackQuery) (map[string]interface{}, error) {
	chatID := callback.From.ID
	userData, err := cp.userDataStorage.GetUserData(chatID)
	if err != nil {
		return nil, err
	}
	var userLang config.Language
	userLangRaw, ok := userData[userdata.LanguageKey]
	if !ok {
		userLang = cp.cfg.DefaultLanguage
	} else {
		userLang, ok = userLangRaw.(config.Language)
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

	command := model.CallbackCommand(callback.Data)

	switch currentStep {
	case steps.SelectStartCommand:
		switch command {
		case model.GetRandomSentenceCommandName:
			payload["text"] = templates.SimpleMessages[userLang][templates.SelectBookForRandomTemplateName]
			payload["reply_markup"] = model.Menus[userLang][model.SelectSourceMenu]
			if err := cp.stepStorage.Push(chatID, steps.GetRandomSentenceMenu); err != nil {
				return nil, err
			}

		case model.AskQuestionCommandName:
			payload["text"] = templates.SimpleMessages[userLang][templates.SelectBookForQuestionTemplateName]
			payload["reply_markup"] = model.Menus[userLang][model.SelectSourceMenu]
			if err := cp.stepStorage.Push(chatID, steps.AskingQuestionMenu); err != nil {
				return nil, err
			}

		case model.LanguageCommandName:
			payload["text"] = templates.SimpleMessages[userLang][templates.ListLanguagesTemplateName]
			var keyboard [][]model.InlineKeyboardButton
			for lang, langName := range config.SupportedLanguages {
				button := model.InlineKeyboardButton{
					Text:         langName,
					CallbackData: model.CallbackCommand(lang),
				}
				keyboard = append(keyboard, []model.InlineKeyboardButton{button})
			}
			payload["reply_markup"] = &model.InlineKeyboardMarkup{InlineKeyboard: keyboard}
			if err := cp.stepStorage.Push(chatID, steps.SelectLanguage); err != nil {
				return nil, err
			}
		case model.HelpCommandName:
			payload["text"] = templates.SimpleMessages[userLang][templates.HelpTemplateName]
			payload["reply_markup"] = model.Menus[userLang][model.StartMenu]
		default:
			payload["text"] = templates.SimpleMessages[userLang][templates.InvalidButtonTemplateName]
		}
	case steps.SelectBook:
		prevStep, err := cp.stepStorage.PeekPrevious(chatID)
		if err != nil {
			return nil, err
		}
		switch prevStep {
		case steps.AskingQuestionMenu:
			idx, err := strconv.Atoi(string(command))
			if err != nil {
				return nil, err
			}
			if err := cp.userDataStorage.SaveUserData(chatID, userdata.BookIDKey, int64(idx)); err != nil {
				return nil, err
			}
			payload["text"] = templates.SimpleMessages[userLang][templates.TypeQuestionTemplateName]
			if err := cp.stepStorage.Push(chatID, steps.AskingQuestion); err != nil {
				return nil, err
			}
		case steps.GetRandomSentenceMenu:
			bookID, err := strconv.Atoi(callback.Data)
			if err != nil {
				return nil, err
			}
			book, err := cp.bookStorage.GetBookByID(ctx, int64(bookID))
			if err != nil {
				return nil, err
			}
			bookData, err := cp.booksRepository.LoadBook(book)
			if err != nil {
				return nil, err
			}
			parser, err := book.GetParser()
			if err != nil {
				return nil, err
			}
			sentence, err := parser.ParseRandomSentence(bookData, time.Now().UnixNano())
			if err != nil {
				return nil, err
			}

			msg, err := templates.GenerateMessageWithData(templates.QuoteTemplateName, updateprocessor.Quote{Text: sentence, Book: book}, userLang)
			if err != nil {
				return nil, err
			}
			payload["text"] = msg
			payload["reply_markup"] = model.Menus[userLang][model.StartMenu]
			cp.stepStorage.Clear(chatID)
			if err := cp.stepStorage.Push(chatID, steps.SelectStartCommand); err != nil {
				return nil, err
			}
		default:
			payload["text"] = templates.SimpleMessages[userLang][templates.InvalidButtonTemplateName]
		}

	case steps.AskingQuestionMenu:
		switch command {
		case model.ListBooksCommandName:
			payload, err = cp.generateListBooksMenuPayload(ctx, chatID, userLang)
			if err != nil {
				return nil, err
			}
			if err := cp.stepStorage.Push(chatID, steps.SelectBook); err != nil {
				return nil, err
			}
		case model.UseRandomBookCommandName:
			payload["text"] = templates.SimpleMessages[userLang][templates.TypeQuestionTemplateName]
			if err := cp.stepStorage.Push(chatID, steps.AskingQuestion); err != nil {
				return nil, err
			}
		case model.GoBackCommandName:
			cp.stepStorage.Clear(chatID)
			payload["text"] = templates.SimpleMessages[userLang][templates.BackTemplateName]
			payload["reply_markup"] = model.Menus[userLang][model.StartMenu]
		default:
			payload["text"] = templates.SimpleMessages[userLang][templates.InvalidButtonTemplateName]
		}
	case steps.GetRandomSentenceMenu:
		switch command {
		case model.ListBooksCommandName:
			payload, err = cp.generateListBooksMenuPayload(ctx, chatID, userLang)
			if err != nil {
				return nil, err
			}
			if err := cp.stepStorage.Push(chatID, steps.SelectBook); err != nil {
				return nil, err
			}
		case model.UseRandomBookCommandName:
			book, err := cp.bookStorage.GetRandomBook(ctx, booksmeta.WithLanguage(userLang))
			if err != nil {
				return nil, err
			}
			bookData, err := cp.booksRepository.LoadBook(book)
			if err != nil {
				return nil, err
			}
			parser, err := book.GetParser()
			if err != nil {
				return nil, err
			}
			sentence, err := parser.ParseRandomSentence(bookData, time.Now().UnixNano())
			if err != nil {
				return nil, err
			}
			msg, err := templates.GenerateMessageWithData(templates.QuoteTemplateName, updateprocessor.Quote{Text: sentence, Book: book}, userLang)
			if err != nil {
				return nil, err
			}
			payload["text"] = msg
			payload["reply_markup"] = model.Menus[userLang][model.StartMenu]
			cp.stepStorage.Clear(chatID)
			if err := cp.stepStorage.Push(chatID, steps.SelectStartCommand); err != nil {
				return nil, err
			}
		case model.GoBackCommandName:
			cp.stepStorage.Clear(chatID)
			payload["text"] = templates.SimpleMessages[userLang][templates.BackTemplateName]
			payload["reply_markup"] = model.Menus[userLang][model.StartMenu]
		default:
			payload["text"] = templates.SimpleMessages[userLang][templates.InvalidButtonTemplateName]
		}
	case steps.SelectLanguage:
		userLang = config.Language(command)
		if err := cp.userDataStorage.SaveUserData(chatID, userdata.LanguageKey, userLang); err != nil {
			return nil, err
		}
		languageName, exist := config.SupportedLanguages[userLang]
		if !exist {
			return nil, fmt.Errorf(`language "%s" is not supported`, userLang)
		}

		msg, err := templates.GenerateMessageWithData(templates.ChangedLanguageTemplateName, languageName, userLang)
		if err != nil {
			return nil, err
		}
		payload["text"] = msg
		payload["reply_markup"] = model.Menus[userLang][model.StartMenu]
		cp.stepStorage.Clear(chatID)
		if err := cp.stepStorage.Push(chatID, steps.SelectStartCommand); err != nil {
			return nil, err
		}
	}

	return payload, nil
}

func (cp *UpdateProcessor) generateListBooksMenuPayload(ctx context.Context, chatID int64, lang config.Language) (map[string]interface{}, error) {
	var keyboard [][]model.InlineKeyboardButton
	books, err := cp.bookStorage.ListBooks(ctx, booksmeta.WithLanguage(lang), booksmeta.Ordered())
	if err != nil {
		return nil, fmt.Errorf(`failed to list books: %w`, err)
	}
	for i, book := range books {
		button := model.InlineKeyboardButton{
			Text:         fmt.Sprintf("%s - %s", book.Author, book.Title),
			CallbackData: model.CallbackCommand(strconv.Itoa(i)),
		}
		keyboard = append(keyboard, []model.InlineKeyboardButton{button})
	}
	payload := map[string]interface{}{
		"chat_id":      chatID,
		"text":         templates.SimpleMessages[lang][templates.ListBooksTemplateName],
		"reply_markup": &model.InlineKeyboardMarkup{InlineKeyboard: keyboard},
	}
	return payload, nil
}
