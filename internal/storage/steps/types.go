package steps

type ChatStep string

const (
	SelectStartCommand    = ChatStep("select_start_command")
	SelectBook            = ChatStep("select_book")
	AskingQuestion        = ChatStep("asking_question")
	AskingQuestionMenu    = ChatStep("asking_question_menu")
	GetRandomSentenceMenu = ChatStep("get_random_sentence_menu")
)
