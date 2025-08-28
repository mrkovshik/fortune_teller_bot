package templates

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed data/*
var templateFS embed.FS

var templates = make(map[string]*template.Template)

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
	ChangeLanguageTemplateName        = "change_language"

	English = "eng"
	Russian = "rus"
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
		ChangeLanguageTemplateName,
	}
	SimpleMessages     map[string]map[string]string
	SupportedLanguages = map[string]string{
		English: "English",
		Russian: "Русский",
	}
)

func GenerateMessageWithData(templateName string, data interface{}, lang string) (string, error) {
	buf := new(bytes.Buffer)
	if err := templates[lang].ExecuteTemplate(buf, templateName, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func generateSimpleMessage(templateName string, lang string) (string, error) {
	buf := new(bytes.Buffer)
	if err := templates[lang].ExecuteTemplate(buf, templateName, nil); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func prepareSimpleMessages() (map[string]map[string]string, error) {
	var err error
	result := make(map[string]map[string]string)
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

func InitTemplates() error {
	var err error

	for language := range SupportedLanguages {
		path := fmt.Sprintf("data/dialog_templates_%s.tpl", language)
		dialogTemplate := template.New("dialogTemplate_" + language)
		templates[language], err = dialogTemplate.ParseFS(templateFS, path)
		if err != nil {
			return fmt.Errorf("error parsing dialog templates: %w", err)
		}
	}

	SimpleMessages, err = prepareSimpleMessages()
	if err != nil {
		return fmt.Errorf("error preparing simpleMessages: %w", err)
	}
	return nil
}
