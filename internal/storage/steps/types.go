package steps

import "errors"

type ChatStep string

var (
	ErrStepNotFound = errors.New("step not found")
	ErrChatNotFound = errors.New("chat not found")
)

const (
	SelectStartCommand    = ChatStep("select_start_command")
	SelectBook            = ChatStep("select_book")
	AskingQuestion        = ChatStep("asking_question")
	AskingQuestionMenu    = ChatStep("asking_question_menu")
	GetRandomSentenceMenu = ChatStep("get_random_sentence_menu")
)
