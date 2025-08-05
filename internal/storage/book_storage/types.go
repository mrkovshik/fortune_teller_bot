package book_storage

type TextParser interface {
	ParseRandomSentence(data []byte) (string, error)
}
