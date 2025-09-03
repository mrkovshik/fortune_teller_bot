package userdata

import (
	"errors"
	"strconv"
	"strings"
)

const (
	BookIDKey   = "book_id"
	LanguageKey = "language"
	Namespace   = "ft_bot"
	Version     = "v1"
)

var ErrNotFound = errors.New("user data not found")

func GenerateKey(chatID int64, dataKey string) string {
	sb := strings.Builder{}
	sb.WriteString(Namespace)
	sb.WriteString(":")
	sb.WriteString(Version)
	sb.WriteString(":")
	sb.WriteString(strconv.FormatInt(chatID, 10))
	sb.WriteString(":")
	sb.WriteString(dataKey)
	return sb.String()
}
