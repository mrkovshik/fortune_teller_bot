package inmemory

import (
	"context"
	"errors"

	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/userdata"
)

type (
	UserDataStorage map[string]any
)

func NewUserDataStorage() UserDataStorage {
	return make(UserDataStorage)
}

func (s UserDataStorage) GetBookID(_ context.Context, chatID int64) (int64, error) {
	bookID, ok := s[userdata.GenerateKey(chatID, userdata.BookIDKey)]
	if !ok {
		return 0, errors.New("book_id not found")
	}
	bookIDInt, ok := bookID.(int64)
	if !ok {
		return 0, errors.New("book_id is not int64")
	}
	return bookIDInt, nil
}

func (s UserDataStorage) GetLanguage(_ context.Context, chatID int64) (config.Language, error) {
	languageRaw, ok := s[userdata.GenerateKey(chatID, userdata.LanguageKey)]
	if !ok {
		return "", errors.New("language not found")
	}
	language, ok := languageRaw.(config.Language)
	if !ok {
		return "", errors.New("language is not Language type")
	}
	return language, nil
}

func (s UserDataStorage) SaveBookID(_ context.Context, chatID, bookID int64) error {
	s[userdata.GenerateKey(chatID, userdata.BookIDKey)] = bookID
	return nil
}

func (s UserDataStorage) SaveLanguage(_ context.Context, chatID int64, language config.Language) error {
	s[userdata.GenerateKey(chatID, userdata.LanguageKey)] = language
	return nil
}
