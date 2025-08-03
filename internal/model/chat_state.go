package model

type ChatState struct {
	CurrentStep ChatStep
	TempData    map[string]string
}

type ChatStep string

const (
	Initial        = ChatStep("Initial")
	SelectBook     = ChatStep("select_book")
	AskingQuestion = ChatStep("asking_question")
)
