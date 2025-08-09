package bookstorage

type TextParser interface {
	ParseRandomSentence(data []byte, seed int64) (string, error)
}
