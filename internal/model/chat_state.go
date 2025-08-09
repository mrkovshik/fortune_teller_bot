package model

import "container/list"

const (
	SelectStartCommand = ChatStep("select_start_command")
	SelectBook         = ChatStep("select_book")
	AskingQuestion     = ChatStep("asking_question")
	AskingQuestionMenu = ChatStep("asking_question_menu")
)

type ChatState struct {
	StepStack *StepStack
	TempData  map[string]string
}
type ChatStep string

// StepStack is a stack data structure implemented using container/list.
// It stores elements in LIFO (Last In, First Out) order.
type StepStack struct {
	steps *list.List
}

// NewStepStack creates and returns a new empty StepStack.
func NewStepStack() *StepStack {
	return &StepStack{
		steps: list.New(),
	}
}

// Push adds an element to the top of the stack.
func (s *StepStack) Push(v ChatStep) {
	s.steps.PushBack(v)
}

// Pop removes and returns the top element of the stack.
// It returns ("", false) if the stack is empty.
func (s *StepStack) Pop() (ChatStep, bool) {
	el := s.steps.Back()
	if el == nil {
		return "", false
	}
	s.steps.Remove(el)
	step := el.Value.(ChatStep)
	return step, true
}

// Peek returns the top element of the stack without removing it.
// It returns ("", false) if the stack is empty.
func (s *StepStack) Peek() (ChatStep, bool) {
	el := s.steps.Back()
	if el == nil {
		return "", false
	}
	step := el.Value.(ChatStep)
	return step, true
}

// PeekPrevious returns the second element from the top of the stack without removing it.
// It returns ("", false) if the stack is empty or has only one element.
func (s *StepStack) PeekPrevious() (ChatStep, bool) {
	if s.steps.Len() < 2 {
		return "", false
	}
	el := s.steps.Back()
	prev := el.Prev()
	if prev == nil {
		return "", false
	}
	step := prev.Value.(ChatStep)
	return step, true
}

// IsEmpty returns true if the stack has no elements.
func (s *StepStack) IsEmpty() bool {
	return s.steps.Len() == 0
}

// Size returns the number of elements currently in the stack.
func (s *StepStack) Size() int {
	return s.steps.Len()
}
