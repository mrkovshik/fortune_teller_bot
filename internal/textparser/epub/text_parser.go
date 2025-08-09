package epub

import (
	"archive/zip"
	"bytes"
	"errors"
	"math/rand"
	"strings"

	"github.com/mrkovshik/fortune_teller_bot/internal/textparser/helpers"
	"go.uber.org/zap"
	"golang.org/x/net/html"
)

type TextParser struct {
	logger *zap.SugaredLogger
}

func NewTextParser(logger *zap.SugaredLogger) *TextParser {
	return &TextParser{
		logger: logger,
	}
}

func (tp *TextParser) ParseRandomSentence(data []byte, seed int64) (string, error) {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}

	var allText strings.Builder

	for _, file := range r.File {
		if !strings.HasSuffix(file.Name, ".xhtml") && !strings.HasSuffix(file.Name, ".html") {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			continue
		}

		node, err := html.Parse(rc)
		rc.Close()
		if err != nil {
			continue
		}

		extractText(node, &allText)
	}

	if allText.Len() == 0 {
		return "", errors.New("no text found in EPUB")
	}
	cleanText := helpers.CleanHTMLContent(allText.String())
	sentences := splitIntoSentences(cleanText)
	if len(sentences) == 0 {
		return "", errors.New("no sentences found")
	}

	rand.New(rand.NewSource(seed))
	return sentences[rand.Intn(len(sentences))], nil
}

func extractText(n *html.Node, b *strings.Builder) {
	if n.Type == html.TextNode {
		b.WriteString(n.Data)
		b.WriteString(" ")
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, b)
	}
}

func splitIntoSentences(text string) []string {

	text = strings.ReplaceAll(text, "\n", " ")
	parts := strings.Split(text, ".")
	var sentences []string
	for _, s := range parts {
		s = strings.TrimSpace(s)
		if len(s) > 20 {
			sentences = append(sentences, s+".")
		}
	}
	return sentences
}
