package templates

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"text/template"
)

//go:embed data/*
var templateFS embed.FS

var templates = make(map[Language]*template.Template)

type Language string

const (
	HelpTemplateName                  = "help"
	StartTemplateName                 = "start"
	InvalidButtonTemplateName         = "invalid_button"
	SelectBookForQuestionTemplateName = "select_book_for_question"
	SelectBookForRandomTemplateName   = "select_book_for_random"
	TypeQuestionTemplateName          = "type_question"
	ListBooksTemplateName             = "list_books"
	BackTemplateName                  = "back"
	ErrorTemplateName                 = "error"
	QuoteTemplateName                 = "quote"
	ChangedLanguageTemplateName       = "changed_language"
	ListLanguagesTemplateName         = "list_languages"

	ListBooksButtonName         = "list_books_button"
	GetRandomSentenceButtonName = "get_random_sentence_button"
	UseRandomBookButtonName     = "use_random_book_button"
	AskQuestionButtonName       = "ask_question_button"
	GoBackButtonName            = "go_back_button"
	HelpButtonName              = "get_help_button"
	LanguageButtonName          = "change_language_button"

	English Language = "eng"
	Russian Language = "rus"
)

var (
	simpleMessagesList = []string{
		HelpTemplateName,
		StartTemplateName,
		InvalidButtonTemplateName,
		SelectBookForQuestionTemplateName,
		SelectBookForRandomTemplateName,
		TypeQuestionTemplateName,
		ListBooksTemplateName,
		BackTemplateName,
		ErrorTemplateName,
		ListLanguagesTemplateName,
	}

	buttonsList = []string{
		GetRandomSentenceButtonName,
		UseRandomBookButtonName,
		AskQuestionButtonName,
		GoBackButtonName,
		HelpButtonName,
		LanguageButtonName,
	}

	SimpleMessages     map[Language]map[string]string
	ButtonsTexts       map[Language]map[string]string
	SupportedLanguages = map[Language]string{
		English: "English",
		Russian: "Русский",
	}
)

func GenerateMessageWithData(templateName string, data interface{}, lang Language) (string, error) {
	tmpl, ok := templates[lang]
	if !ok || tmpl == nil {
		return "", fmt.Errorf("no templates for lang %q", lang)
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, templateName, data); err != nil {
		return "", fmt.Errorf("execute template %q: %w", templateName, err)
	}
	return buf.String(), nil
}

func generateSimpleMessage(templateName string, lang Language) (string, error) {
	buf := new(bytes.Buffer)
	if err := templates[lang].ExecuteTemplate(buf, templateName, nil); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func prepareSimpleMessages() (map[Language]map[string]string, error) {
	var err error
	result := make(map[Language]map[string]string)
	for language := range SupportedLanguages {
		result[language] = make(map[string]string)
		for _, message := range simpleMessagesList {
			result[language][message], err = generateSimpleMessage(message, language)
			if err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func prepareButtonsTexts() (map[Language]map[string]string, error) {
	var err error
	result := make(map[Language]map[string]string)
	for language := range SupportedLanguages {
		result[language] = make(map[string]string)
		for _, buttonName := range buttonsList {
			result[language][buttonName], err = generateSimpleMessage(buttonName, language)
			if err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func init() {
	var err error

	for language := range SupportedLanguages {
		path := fmt.Sprintf("data/dialog_templates_%s.tpl", language)
		dialogTemplate := template.New("dialogTemplate_" + string(language))
		templates[language], err = dialogTemplate.ParseFS(templateFS, path)
		if err != nil {
			log.Fatal(fmt.Errorf("error parsing dialog templates: %w", err))
		}
	}

	SimpleMessages, err = prepareSimpleMessages()
	if err != nil {
		log.Fatal("error preparing simpleMessages: %w", err)
	}
	ButtonsTexts, err = prepareButtonsTexts()
	if err != nil {
		log.Fatal("error preparing ButtonsTexts: %w", err)
	}
}
