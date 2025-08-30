package model

import (
	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/fortune_teller_bot/internal/embedded/templates"
)

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text              string          `json:"text"`
	CallbackData      CallbackCommand `json:"callback_data,omitempty"`
	URL               string          `json:"url,omitempty"`
	SwitchInlineQuery string          `json:"switch_inline_query,omitempty"`
}

var Menus = make(map[config.Language]map[string]InlineKeyboardMarkup)

type CallbackCommand string

const (
	ListBooksCommandName         CallbackCommand = "list_books"
	GetRandomSentenceCommandName CallbackCommand = "get_random_sentence"
	UseRandomBookCommandName     CallbackCommand = "use_random_book"
	AskQuestionCommandName       CallbackCommand = "ask_question"
	GoBackCommandName            CallbackCommand = "go_back"
	HelpCommandName              CallbackCommand = "get_help"
	LanguageCommandName          CallbackCommand = "change_language"

	StartMenu        = "start_menu"
	SelectSourceMenu = "select_source_menu"
)

var (
	ButtonsToCommands = map[string]CallbackCommand{
		templates.ListBooksButtonName:         ListBooksCommandName,
		templates.GetRandomSentenceButtonName: GetRandomSentenceCommandName,
		templates.UseRandomBookButtonName:     UseRandomBookCommandName,
		templates.AskQuestionButtonName:       AskQuestionCommandName,
		templates.GoBackButtonName:            GoBackCommandName,
		templates.HelpButtonName:              HelpCommandName,
		templates.LanguageButtonName:          LanguageCommandName,
	}
	menusToButtons = map[string][]string{
		StartMenu:        {templates.GetRandomSentenceButtonName, templates.AskQuestionButtonName, templates.LanguageButtonName, templates.HelpButtonName},
		SelectSourceMenu: {templates.ListBooksButtonName, templates.UseRandomBookButtonName, templates.GoBackButtonName},
	}
)

func init() {
	for lang := range config.SupportedLanguages {
		Menus[lang] = make(map[string]InlineKeyboardMarkup)
		for menuName, buttonsNames := range menusToButtons {
			var keyboard [][]InlineKeyboardButton
			for _, buttonName := range buttonsNames {
				keyboard = append(keyboard, []InlineKeyboardButton{{
					Text:         templates.ButtonsTexts[lang][buttonName],
					CallbackData: ButtonsToCommands[buttonName],
				},
				})
			}
			Menus[lang][menuName] = InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			}
		}
	}
}
