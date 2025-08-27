package templates

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed data/*
var templateFS embed.FS

var templates *template.Template

const (
	HelpTemplateName                  = "help"
	StartTemplateName                 = "start"
	InvalidButtonTemplateName         = "invalid_button"
	SelectBookForQuestionTemplateName = "select_book_for_question"
	SelectBookForRandomTemplateName   = "select_book_for_random"
	TypeQuestionTemplateName          = "type_question"
	ListBooksTemplateName             = "list_books"
	BackTemplateName                  = "back"
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
	}
	SimpleMessages     map[string]string
	SupportedLanguages = map[string]struct{}{
		"rus": {},
		"eng": {},
	}
)

func GenerateMessageWithData(templateName string, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	if err := templates.ExecuteTemplate(buf, templateName, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func generateSimpleMessage(templateName string) (string, error) {
	buf := new(bytes.Buffer)
	if err := templates.ExecuteTemplate(buf, templateName, nil); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func prepareSimpleMessages() (map[string]string, error) {
	var err error
	result := make(map[string]string)
	for _, message := range simpleMessagesList {
		result[message], err = generateSimpleMessage(message)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func InitTemplates(lang string) error {
	var err error
	dialogTemplate := template.New("dialogTemplate")
	path := fmt.Sprintf("data/dialog_templates_%s.tpl", lang)
	templates, err = dialogTemplate.ParseFS(templateFS, path)
	if err != nil {
		return fmt.Errorf("error parsing dialog templates: %w", err)
	}
	SimpleMessages, err = prepareSimpleMessages()
	if err != nil {
		return fmt.Errorf("error preparing simpleMessages: %w", err)
	}
	return nil
}
