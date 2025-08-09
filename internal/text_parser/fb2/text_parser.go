package fb2

import (
	"encoding/xml"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/mrkovshik/fortune_teller_bot/internal/text_parser/helpers"
	"go.uber.org/zap"
)

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
type TextParcer struct {
	logger *zap.SugaredLogger
}

func NewTextParser(logger *zap.SugaredLogger) *TextParcer {
	return &TextParcer{
		logger: logger,
	}
}

func (tp *TextParcer) ParseRandomSentence(data []byte, seed int64) (string, error) {
	if len(data) == 0 {
		return "", errors.New("empty data")
	}
	var book FictionBook
	if err := xml.Unmarshal(data, &book); err != nil {
		return "", fmt.Errorf("can't parse XML: %w", err)
	}

	var sentences []string
	for _, section := range book.Body.Sections {
		for _, p := range section.Paragraphs {
			text := strings.TrimSpace(p.Text)
			cleanText := helpers.CleanHTMLContent(text)
			if len(cleanText) > 20 && !strings.HasSuffix(text, ":") {
				sentences = append(sentences, text)
			}
		}
	}

	if len(sentences) == 0 {
		return "", fmt.Errorf("no usable paragraphs found")
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return sentences[rand.Intn(len(sentences))], nil
}
