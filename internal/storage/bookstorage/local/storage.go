package local

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/mrkovshik/fortune_teller_bot/internal/embedded/books"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/bookstorage"
	"github.com/mrkovshik/fortune_teller_bot/internal/textparser/epub"
	"github.com/mrkovshik/fortune_teller_bot/internal/textparser/fb2"
	"github.com/mrkovshik/fortune_teller_bot/internal/updateprocessor"
	"go.uber.org/zap"
)

type Storage struct {
	fs     embed.FS
	logger *zap.SugaredLogger
}

func NewStorage(logger *zap.SugaredLogger) *Storage {
	return &Storage{
		fs:     books.GetBooksFS(),
		logger: logger,
	}
}

func (s *Storage) GetRandomSentenceFromBook(bookTitle string, seed int64) (*updateprocessor.Quote, error) {
	var parser bookstorage.TextParser
	fileName, exists := TitleToFileName[bookTitle]
	if !exists {
		return nil, fmt.Errorf("book title '%s' not exists", bookTitle)
	}
	data, err := s.fs.ReadFile("data/" + fileName)
	if err != nil {
		return nil, fmt.Errorf("can't read book: %w", err)
	}
	switch {
	case strings.HasSuffix(fileName, ".fb2"):
		parser = fb2.NewTextParser(s.logger)
	case strings.HasSuffix(fileName, ".epub"):
		parser = epub.NewTextParser(s.logger)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", fileName)
	}

	sentence, err := parser.ParseRandomSentence(data, seed)
	if err != nil {
		return nil, err
	}
	reply := &updateprocessor.Quote{
		Text:  sentence,
		Title: bookTitle,
	}
	return reply, nil
}

func (s *Storage) ListBooks() ([]string, error) {
	entries, err := fs.ReadDir(s.fs, "data")
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

func (s *Storage) GetRandomBookTitle() string {
	return GetRandomBookTitle()
}
