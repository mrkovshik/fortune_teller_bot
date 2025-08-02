package local

import (
	"embed"
	"encoding/xml"
	"fmt"
	"io/fs"
	"math/rand"
	"strings"
	"time"

	"github.com/mrkovshik/fortune_teller_bot/internal/embedded"
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

type Paragraph struct {
	Text string `xml:",chardata"`
}

type Section struct {
	Paragraphs []Paragraph `xml:"p"`
}

type Body struct {
	Sections []Section `xml:"section"`
}

type FictionBook struct {
	Body Body `xml:"body"`
}

func (s *Storage) GetRandomSentenceFromBook(bookName string) (string, error) {
	sentences, err := s.LoadBookText(titleToFileName[bookName])
	if err != nil {
		return "", err
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return sentences[rand.Intn(len(sentences))], nil
}

func (s *Storage) LoadBookText(bookName string) ([]string, error) {
	bookPath := "books/" + bookName
	data, err := s.fs.ReadFile(bookPath)
	if err != nil {
		return nil, fmt.Errorf("can't read book: %w", err)
	}

	var book FictionBook
	if err := xml.Unmarshal(data, &book); err != nil {
		return nil, fmt.Errorf("can't parse XML: %w", err)
	}

	var sentences []string
	for _, section := range book.Body.Sections {
		for _, p := range section.Paragraphs {
			text := strings.TrimSpace(p.Text)
			if len(text) > 20 {
				sentences = append(sentences, text)
			}
		}
	}

	if len(sentences) == 0 {
		return nil, fmt.Errorf("no usable paragraphs found")
	}

	return sentences, nil
}

func (s *Storage) ListBooks() ([]string, error) {
	entries, err := fs.ReadDir(s.fs, "books")
	if err != nil {
		return nil, err
	}

	var bookNames []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".fb2") {

			bookNames = append(bookNames, fileNameToTitle[entry.Name()])
		}
	}
	return bookNames, nil
}
