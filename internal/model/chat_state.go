package model

type ChatState struct {
	CurrentStep ChatStep
	TempData    map[string]string
}

type ChatStep string

const (
	Initial        = ChatStep("Initial")
	BookChoosing   = ChatStep("BookChoosing")
	AskingQuestion = ChatStep("AskingQuestion")
)
