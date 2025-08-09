package textparser

type TextParcer interface {
	ParseRandomSentence(data []byte) (string, error)
}
