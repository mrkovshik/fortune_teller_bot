package local

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/mrkovshik/fortune_teller_bot/internal/embedded"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/book_storage"
	"github.com/mrkovshik/fortune_teller_bot/internal/text_parser/epub"
	"github.com/mrkovshik/fortune_teller_bot/internal/text_parser/fb2"
	"go.uber.org/zap"
)

type Storage struct {
	fs     embed.FS
	logger *zap.SugaredLogger
}

func NewStorage(logger *zap.SugaredLogger) *Storage {
	return &Storage{
		fs:     embedded.GetBooksFS(),
		logger: logger,
	}
}

func (s *Storage) GetRandomSentenceFromBook(bookName string, seed int64) (string, error) {
	var parser bookstorage.TextParser
	fileName, exists := TitleToFileName[bookName]
	if !exists {
		return "К сожалению, пока такой книги у нас нет( Пожалуйста, выберите книгу из списка", nil
	}
	data, err := s.fs.ReadFile("books/" + fileName)
	if err != nil {
		return "", fmt.Errorf("can't read book: %w", err)
	}
	switch {
	case strings.HasSuffix(fileName, ".fb2"):
		parser = fb2.NewTextParser(s.logger)
	case strings.HasSuffix(fileName, ".epub"):
		parser = epub.NewTextParser(s.logger)
	default:
		return "", fmt.Errorf("unsupported file type: %s", fileName)
	}

	sentence, err := parser.ParseRandomSentence(data, seed)
	if err != nil {
		return "", err
	}
	return sentence, nil
}

func (s *Storage) ListBooks() ([]string, error) {
	entries, err := fs.ReadDir(s.fs, "books")
	if err != nil {
		return nil, err
	}

	var bookNames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			bookTitle, exist := FileNameToTitle[entry.Name()]
			if !exist {
				s.logger.Warnw("can't find book title for file. Please add it to 'FileNameToTitle' map or delete the file", "name", entry.Name())
				continue
			}
			bookNames = append(bookNames, bookTitle)
		}
	}
	return bookNames, nil
}
