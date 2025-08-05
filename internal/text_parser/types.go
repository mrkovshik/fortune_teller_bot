package text_parser

type TextParcer interface {
	ParseRandomSentence(data []byte) (string, error)
}
